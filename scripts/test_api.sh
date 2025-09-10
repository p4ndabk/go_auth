#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo -e "${YELLOW}=== API Auth Test Script ===${NC}"
echo ""

# Test 1: Register User
echo -e "${YELLOW}1. Testing user registration...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com", "password": "password123"}')

if echo "$REGISTER_RESPONSE" | grep -q "User created successfully"; then
    echo -e "${GREEN}✓ User registration successful${NC}"
    echo "$REGISTER_RESPONSE" | jq .
else
    echo -e "${RED}✗ User registration failed${NC}"
    echo "$REGISTER_RESPONSE"
fi

echo ""

# Test 2: Login
echo -e "${YELLOW}2. Testing user login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    echo -e "${GREEN}✓ User login successful${NC}"
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
    echo "Token: $TOKEN"
else
    echo -e "${RED}✗ User login failed${NC}"
    echo "$LOGIN_RESPONSE"
    exit 1
fi

echo ""

# Test 3: Get User Info
echo -e "${YELLOW}3. Testing /me endpoint...${NC}"
ME_RESPONSE=$(curl -s -X GET "$BASE_URL/me" \
  -H "Authorization: Bearer $TOKEN")

if echo "$ME_RESPONSE" | grep -q "username"; then
    echo -e "${GREEN}✓ /me endpoint successful${NC}"
    echo "$ME_RESPONSE" | jq .
else
    echo -e "${RED}✗ /me endpoint failed${NC}"
    echo "$ME_RESPONSE"
fi

echo ""

# Test 4: Unauthorized access
echo -e "${YELLOW}4. Testing unauthorized access...${NC}"
UNAUTH_RESPONSE=$(curl -s -X GET "$BASE_URL/me")

if echo "$UNAUTH_RESPONSE" | grep -q "Authorization header required"; then
    echo -e "${GREEN}✓ Unauthorized access properly blocked${NC}"
else
    echo -e "${RED}✗ Unauthorized access test failed${NC}"
fi

echo ""

# Test 5: Invalid token
echo -e "${YELLOW}5. Testing invalid token...${NC}"
INVALID_RESPONSE=$(curl -s -X GET "$BASE_URL/me" \
  -H "Authorization: Bearer invalid_token")

if echo "$INVALID_RESPONSE" | grep -q "Invalid token"; then
    echo -e "${GREEN}✓ Invalid token properly rejected${NC}"
else
    echo -e "${RED}✗ Invalid token test failed${NC}"
fi

echo ""
echo -e "${GREEN}=== All tests completed ===${NC}"
