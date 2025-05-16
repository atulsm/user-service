#!/bin/bash

# Configuration
API_URL="http://localhost:8080"
EMAIL="test@example.com"
PASSWORD="Password123!"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "🔍 Testing User Service API"
echo "=========================="

# Test health endpoint
echo -e "\n1. Testing Health Check"
curl -s "$API_URL/health" | jq .

# Register new user
echo -e "\n2. Testing User Registration"
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\",
    \"first_name\": \"Test\",
    \"last_name\": \"User\",
    \"phone_number\": \"+1234567890\"
  }")
echo "$REGISTER_RESPONSE" | jq .

# Extract user ID from registration response
USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.user.id')

# Login
echo -e "\n3. Testing Login"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/users/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")
echo "$LOGIN_RESPONSE" | jq .

# Extract token from login response
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ Login successful, got token${NC}"

    # Get profile
    echo -e "\n4. Testing Get Profile"
    curl -s "$API_URL/api/users/profile" \
      -H "Authorization: Bearer $TOKEN" | jq .

    # Update profile
    echo -e "\n5. Testing Update Profile"
    curl -s -X PUT "$API_URL/api/users/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "first_name": "Updated",
        "last_name": "Name",
        "phone_number": "+1987654321"
      }' | jq .

    # List users
    echo -e "\n6. Testing List Users"
    curl -s "$API_URL/api/users?limit=10&offset=0" \
      -H "Authorization: Bearer $TOKEN" | jq .

    # Get specific user
    echo -e "\n7. Testing Get Specific User"
    curl -s "$API_URL/api/users/$USER_ID" \
      -H "Authorization: Bearer $TOKEN" | jq .

    # Delete user
    echo -e "\n8. Testing Delete User"
    curl -s -X DELETE "$API_URL/api/users/$USER_ID" \
      -H "Authorization: Bearer $TOKEN" | jq .
else
    echo -e "${RED}✗ Login failed, cannot continue with authenticated requests${NC}"
fi 