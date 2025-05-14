#!/bin/bash

# Exit on error
set -e

# Default values
DB_NAME="userservice_test"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_HOST="localhost"
DB_PORT="5432"

# Parse command line arguments
while getopts "n:u:p:h:P:v" opt; do
    case $opt in
        n) DB_NAME=$OPTARG ;;
        u) DB_USER=$OPTARG ;;
        p) DB_PASSWORD=$OPTARG ;;
        h) DB_HOST=$OPTARG ;;
        P) DB_PORT=$OPTARG ;;
        v) VERBOSE=1 ;;
        ?) echo "Usage: $0 [-n db_name] [-u db_user] [-p db_password] [-h db_host] [-P db_port] [-v]" && exit 1 ;;
    esac
done

echo "Setting up test environment..."

# Setup test database
./scripts/setup_test_db.sh -n "$DB_NAME" -u "$DB_USER" -p "$DB_PASSWORD" -h "$DB_HOST" -P "$DB_PORT"

# Source the environment variables
source env.sh

# Run the API tests
echo "Running API tests..."
if [ -n "$VERBOSE" ]; then
    go test -v ./tests/api/...
else
    go test ./tests/api/...
fi 