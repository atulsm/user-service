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

## Performance

The service has been load tested using `wrk` with the following results:

```bash
$ wrk -t1 -c8 -d15 --latency http://localhost:8080/users
Running 15s test @ http://localhost:8080/users
  1 threads and 8 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.99ms    2.36ms  38.97ms   90.89%
    Req/Sec    23.10k     1.87k   25.70k    84.11%
  Latency Distribution
     50%  288.00us
     75%  375.00us
     90%    2.91ms
     99%   11.50ms
  347064 requests in 15.10s, 470.66MB read
Requests/sec:  22981.57
Transfer/sec:     31.17MB
```

### Key Performance Metrics
- Average Latency: 0.99ms
- Requests per Second: 22,981
- 99th Percentile Latency: 11.50ms
- Throughput: 31.17MB/s

## Getting Started

### Prerequisites
- Go 1.21 or later
- Node.js 18 or later
- PostgreSQL 14 or later

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/user-service.git
cd user-service
```

2. Set up the backend:
```bash
cd internal
go mod download
go run cmd/server/main.go
```

3. Set up the frontend:
```bash
cd web
npm install
npm start
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## API Documentation

### Authentication Endpoints

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/logout` - Logout user

### User Endpoints

- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users/profile` - Get current user profile
- `PUT /api/v1/users/profile` - Update current user profile

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

