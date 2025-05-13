# user-service
High performant Go API server for a user service. This will include standard functionality like user registration, authentication, profile management, and CRUD operations.

This service includes:

# User Management

Registration
Authentication with JWT
Profile management
User listing with pagination
User deletion


# Project Structure

Well-organized Go project following standard practices
Clear separation of concerns (handlers, repositories, models)
Configuration management


# Technical Features

RESTful API with Gin framework
PostgreSQL database integration
JWT-based authentication
Password hashing with bcrypt
Graceful server shutdown



# The service contains all the necessary endpoints:

POST /api/users/register: Create new user account
POST /api/users/login: Authenticate and get JWT token
GET /api/users/profile: Get current user profile (protected)
PUT /api/users/profile: Update user profile (protected)
GET /api/users: List all users (protected)
GET /api/users/:id: Get specific user (protected)
DELETE /api/users/:id: Delete user (protected)

# build
go mod tidy
go build -o userservice ./cmd/server

# database setup
psql -U postgres
CREATE DATABASE userservice;
\c userservice
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_users_email ON users(email);
\q

# To run this service, you'll need:

Go installed (1.16+)
PostgreSQL database
Set environment variables:

DATABASE_URL
JWT_SECRET
PORT (optional, defaults to 8080)

export DATABASE_URL="postgres://postgres:postgres@localhost:5432/userservice?sslmode=disable"
export JWT_SECRET="your-secure-random-string"
./userservice

