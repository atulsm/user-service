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

### HTTP API Load Test Results

Load test results using `wrk`:

```bash
wrk -t12 -c400 -d30s http://localhost:8080/users
```

```
Running 30s test @ http://localhost:8080/users
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     0.99ms    1.22ms  31.25ms   97.65%
    Req/Sec    22.98k     4.23k   35.00k    68.50%
  689431 requests in 30.00s, 935.05MB read
Requests/sec:  22,981.03
Transfer/sec:     31.17MB
```

Key metrics:
- Average latency: 0.99ms
- Requests per second: 22,981
- Throughput: 31.17MB/s

### gRPC API Load Test Results

Load test results using `ghz`:

```bash
ghz --proto=proto/user.proto --call=user.UserService/GetUsers --insecure --duration=15s -c 5 localhost:50051
```

```
Summary:
  Count:        187335
  Total:        15.00 s
  Slowest:      15.72 ms
  Fastest:      0.17 ms
  Average:      0.35 ms
  Requests/sec: 12490.50

Response time histogram:
  0.169  [1]      |
  1.724  [187267] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  3.279  [48]     |
  4.834  [5]      |
  6.389  [0]      |
  7.944  [0]      |
  9.498  [0]      |
  11.053 [5]      |
  12.608 [0]      |
  14.163 [0]      |
  15.718 [5]      |

Latency distribution:
  10 % in 0.26 ms 
  25 % in 0.29 ms 
  50 % in 0.33 ms 
  75 % in 0.38 ms 
  90 % in 0.45 ms 
  95 % in 0.50 ms 
  99 % in 0.72 ms 

Status code distribution:
  [OK]            187331 responses   
  [Unavailable]   4 responses
```

Key metrics:
- Average latency: 0.35ms
- Requests per second: 12,490
- 99th percentile latency: 0.72ms
- Success rate: 99.998% (187,331/187,335)

The gRPC endpoint shows significantly lower latency compared to the HTTP endpoint, making it ideal for high-performance microservices communication.

## Testing gRPC Service

The service exposes a gRPC endpoint for the `/users` endpoint. You can test it using `grpcurl`, a command-line tool for interacting with gRPC servers.

### Installation

```bash
# macOS
brew install grpcurl

# Linux
sudo apt-get install grpcurl  # Ubuntu/Debian
```

### Basic Usage

1. **List available services:**
```bash
grpcurl -plaintext localhost:50051 list
```

2. **Get service information:**
```bash
grpcurl -plaintext -proto proto/user.proto localhost:50051 describe user.UserService
```

3. **Get users with pagination:**
```bash
# Get first page with 5 users
grpcurl -plaintext -proto proto/user.proto -d '{"page": 1, "page_size": 5}' localhost:50051 user.UserService/GetUsers

# Get second page with 3 users
grpcurl -plaintext -proto proto/user.proto -d '{"page": 2, "page_size": 3}' localhost:50051 user.UserService/GetUsers
```

4. **Pretty print output (requires jq):**
```bash
grpcurl -plaintext -proto proto/user.proto -d '{"page": 1, "page_size": 5}' localhost:50051 user.UserService/GetUsers | jq
```

### Command Options

- `-plaintext`: Use plaintext (no TLS)
- `-proto proto/user.proto`: Specify the proto file
- `-d '{"page": 1, "page_size": 5}'`: Send request data
- `localhost:50051`: Server address
- `user.UserService/GetUsers`: Service and method name

### Example Response

```json
{
  "users": [
    {
      "id": "5236353d-1d1a-4369-ba1a-f5b313f418a1",
      "email": "anjana@example.com",
      "firstName": "Anjana",
      "lastName": "Paulose",
      "createdAt": "2025-05-19T22:13:43Z",
      "updatedAt": "2025-05-19T22:13:43Z"
    }
  ],
  "total": 6,
  "page": 1,
  "pageSize": 5
}
```

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

