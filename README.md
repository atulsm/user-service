# User Service 🚀

A high-performance Go API server providing comprehensive user management functionality including registration, authentication, profile management, and CRUD operations.

## 📋 Features

### User Management
- 👤 User registration and authentication
- 🔐 JWT-based authentication
- 👥 Profile management
- 📃 User listing with pagination
- 🗑️ User account deletion

### Technical Stack
- 🛠️ RESTful API with Gin framework
- 🗄️ PostgreSQL database integration
- 🔒 Secure password hashing with bcrypt
- 🔄 Graceful server shutdown
- 📦 Well-organized project structure
- 🔄 Database migrations with Liquibase

## 🏗️ Project Structure

```
user-service/
├── cmd/
│   └── server/          # Application entry point
├── db/
│   └── changelog/       # Database migration files
│       ├── changes/     # Individual change files
│       └── db.changelog-master.xml
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   └── repository/      # Database operations
├── pkg/
│   └── utils/          # Shared utilities
├── scripts/
│   └── test.sh         # Test runner script
├── liquibase.properties # Liquibase configuration
└── README.md
```

## 🚀 Getting Started

### Prerequisites
- Go 1.16 or higher
- PostgreSQL database
- Git
- Liquibase (for database migrations)

### Environment Variables
```bash
# Required in production
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/userservice?sslmode=disable"
export JWT_SECRET="your-secure-random-string"

# Optional
export PORT="8080"  # defaults to 8080
```

### Database Setup

#### Using Liquibase (Recommended)
```bash
# Install Liquibase
brew install liquibase

# Create database
createdb userservice

# Run migrations
liquibase update
```

#### Manual Setup (Alternative)
```sql
-- Connect to PostgreSQL
psql -U postgres

-- Create database
CREATE DATABASE userservice;
\c userservice

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone_number ON users(phone_number);
```

### Database Migrations

The project uses Liquibase for database version control. Migration files are located in the `db/changelog` directory.

```bash
# View pending changes
liquibase status

# Apply pending changes
liquibase update

# Rollback last change
liquibase rollbackCount 1

# Generate change log from existing database
liquibase generateChangeLog

# Validate change log
liquibase validate
```

### Building and Running
```bash
# Install dependencies
go mod tidy

# Build the service
go build -o userservice ./cmd/server

# Run the service
./userservice
```

## 🔍 API Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| POST | `/api/users/register` | Create new user account | Public |
| POST | `/api/users/login` | Authenticate and get JWT token | Public |
| GET | `/api/users/profile` | Get current user profile | Required |
| PUT | `/api/users/profile` | Update user profile | Required |
| GET | `/api/users` | List all users | Required |
| GET | `/api/users/:id` | Get specific user | Required |
| DELETE | `/api/users/:id` | Delete user | Required |

## 🧪 Running Tests

### Setting Up Test Database

The project includes a script to set up a test database with sample data:

```bash
# Set up test database with default settings
./scripts/setup_test_db.sh

# Custom database settings
./scripts/setup_test_db.sh -n custom_db_name -u db_user -p db_password -h db_host -P db_port
```

### Test Users

The following test users are available in the test database:

1. Admin User
   - Email: admin@example.com
   - Password: Admin123!
   - Phone: +1234567890

2. Regular User
   - Email: user@example.com
   - Password: User123!
   - Phone: +1987654321

3. No Phone User
   - Email: nophone@example.com
   - Password: Test123!
   - Phone: (not set)

4. Inactive User
   - Email: inactive@example.com
   - Password: Inactive123!
   - Phone: +1555555555

### Running Test Suite

```bash
# Run all tests
./scripts/test.sh

# Run tests with coverage report
./scripts/test.sh -c

# Run tests in verbose mode
./scripts/test.sh -v

# Run tests with race detector
./scripts/test.sh -r

# Run tests for a specific package
./scripts/test.sh -p ./internal/handlers

# Run tests with multiple options
./scripts/test.sh -cvr
```

### Test Options
| Option | Description |
|--------|-------------|
| `-c` | Generate coverage report (creates coverage/coverage.html) |
| `-v` | Run tests in verbose mode |
| `-r` | Run tests with race detector |
| `-s` | Run tests in short mode |
| `-p` | Run tests for specific package |
| `-h` | Show help message |

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

