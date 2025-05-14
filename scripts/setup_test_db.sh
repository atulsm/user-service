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
while getopts "n:u:p:h:P:" opt; do
    case $opt in
        n) DB_NAME=$OPTARG ;;
        u) DB_USER=$OPTARG ;;
        p) DB_PASSWORD=$OPTARG ;;
        h) DB_HOST=$OPTARG ;;
        P) DB_PORT=$OPTARG ;;
        ?) echo "Usage: $0 [-n db_name] [-u db_user] [-p db_password] [-h db_host] [-P db_port]" && exit 1 ;;
    esac
done

echo "Setting up test database..."

# Drop database if it exists
PGPASSWORD=$DB_PASSWORD dropdb -U $DB_USER -h $DB_HOST -p $DB_PORT $DB_NAME --if-exists

# Create fresh database
PGPASSWORD=$DB_PASSWORD createdb -U $DB_USER -h $DB_HOST -p $DB_PORT $DB_NAME

# Set DATABASE_URL for Liquibase
export DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Run Liquibase migrations with test context
echo "Running migrations with test data..."
liquibase --contexts=test update

echo "Test database setup complete!"
echo "Test database connection string: $DATABASE_URL"

# Print test user credentials
echo -e "\nTest Users Available:"
echo "1. Admin User"
echo "   Email: admin@example.com"
echo "   Password: Admin123!"
echo "   Phone: +1234567890"
echo ""
echo "2. Regular User"
echo "   Email: user@example.com"
echo "   Password: User123!"
echo "   Phone: +1987654321"
echo ""
echo "3. No Phone User"
echo "   Email: nophone@example.com"
echo "   Password: Test123!"
echo ""
echo "4. Inactive User"
echo "   Email: inactive@example.com"
echo "   Password: Inactive123!"
echo "   Phone: +1555555555" 