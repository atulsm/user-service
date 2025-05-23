package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/atulsm/user-service/internal/models"
	"github.com/atulsm/user-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	CreateUser(user *models.RegisterRequest) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uuid.UUID, updates *models.UpdateProfileRequest) (*models.User, error)
	ListUsers(limit, offset int) ([]*models.User, error)
	DeleteUser(id uuid.UUID) error
	Close() error
	UpdatePassword(id uuid.UUID, passwordHash string) error
	GetUsers(ctx context.Context, page, pageSize int) ([]*models.User, int, error)
}

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(dbURL string) UserRepository {
	// Parse the database URL to extract database name
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		log.Printf("Warning: Could not parse database URL: %v", err)
	} else {
		dbName := parsedURL.Path[1:] // Remove leading '/'
		log.Printf("Connecting to database: %s on host: %s", dbName, parsedURL.Host)
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Printf("Successfully connected to PostgreSQL database")
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(req *models.RegisterRequest) (*models.User, error) {
	// Check if user with this email already exists
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", req.Email)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		ID:          uuid.New(),
		Email:       req.Email,
		Password:    passwordHash,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert user into database
	_, err = r.db.NamedExec(`
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone_number, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :first_name, :last_name, :phone_number, :created_at, :updated_at)
	`, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepository) UpdateUser(id uuid.UUID, updates *models.UpdateProfileRequest) (*models.User, error) {
	// Get current user
	user, err := r.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.FirstName != "" {
		user.FirstName = updates.FirstName
	}
	if updates.LastName != "" {
		user.LastName = updates.LastName
	}
	if updates.PhoneNumber != "" {
		user.PhoneNumber = sql.NullString{String: updates.PhoneNumber, Valid: true}
	}
	if updates.Email != "" && updates.Email != user.Email {
		// Check if email is already taken
		var count int
		err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2", updates.Email, id)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("email already in use")
		}
		user.Email = updates.Email
	}

	user.UpdatedAt = time.Now()

	// Save updates
	_, err = r.db.NamedExec(`
		UPDATE users 
		SET first_name = :first_name, 
			last_name = :last_name, 
			email = :email, 
			phone_number = :phone_number,
			updated_at = :updated_at
		WHERE id = :id
	`, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) ListUsers(limit, offset int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	log.Printf("Executing ListUsers query with limit: %d, offset: %d", limit, offset)

	users := []*models.User{}
	err := r.db.Select(&users, "SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Printf("Database error in ListUsers: %v", err)
		return nil, err
	}

	log.Printf("Successfully retrieved %d users from database", len(users))
	return users, nil
}

func (r *PostgresUserRepository) DeleteUser(id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *PostgresUserRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresUserRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET password_hash = $1,
			updated_at = NOW()
		WHERE id = $2
	`, passwordHash, id)
	return err
}

func (r *PostgresUserRepository) GetUsers(ctx context.Context, page, pageSize int) ([]*models.User, int, error) {
	offset := (page - 1) * pageSize
	users, err := r.ListUsers(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.Get(&total, "SELECT COUNT(*) FROM users")
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
