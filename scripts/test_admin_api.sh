#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo -e "${YELLOW}=== Admin API Test Script ===${NC}"
echo ""

# First, login to get JWT token
echo -e "${YELLOW}1. Getting JWT token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "joao@email.com", "password": "123456"}')

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓ Token obtained${NC}"
else
    echo -e "${RED}✗ Failed to get token. Make sure user joao@email.com exists${NC}"
    exit 1
fi

AUTH_HEADER="Authorization: Bearer $TOKEN"

echo ""

# Test Applications
echo -e "${BLUE}=== TESTING APPLICATIONS ===${NC}"

echo -e "${YELLOW}2. Creating application...${NC}"
CREATE_APP_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/applications" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d '{"slug": "test-app", "name": "Test Application", "description": "A test application"}')

if echo "$CREATE_APP_RESPONSE" | grep -q "Application created successfully"; then
    echo -e "${GREEN}✓ Application created${NC}"
    APP_ID=$(echo "$CREATE_APP_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "Application ID: $APP_ID"
else
    echo -e "${RED}✗ Failed to create application${NC}"
    echo "$CREATE_APP_RESPONSE"
fi

echo ""

echo -e "${YELLOW}3. Getting applications...${NC}"
GET_APPS_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/applications" \
  -H "$AUTH_HEADER")

if echo "$GET_APPS_RESPONSE" | grep -q "applications"; then
    echo -e "${GREEN}✓ Applications retrieved${NC}"
else
    echo -e "${RED}✗ Failed to get applications${NC}"
fi

echo ""

# Test Roles
echo -e "${BLUE}=== TESTING ROLES ===${NC}"

echo -e "${YELLOW}4. Creating role...${NC}"
CREATE_ROLE_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/roles" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d "{\"application_id\": $APP_ID, \"slug\": \"test-role\", \"name\": \"Test Role\", \"description\": \"A test role\"}")

if echo "$CREATE_ROLE_RESPONSE" | grep -q "Role created successfully"; then
    echo -e "${GREEN}✓ Role created${NC}"
    ROLE_ID=$(echo "$CREATE_ROLE_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "Role ID: $ROLE_ID"
else
    echo -e "${RED}✗ Failed to create role${NC}"
    echo "$CREATE_ROLE_RESPONSE"
fi

echo ""

echo -e "${YELLOW}5. Getting roles...${NC}"
GET_ROLES_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/roles?application_id=$APP_ID" \
  -H "$AUTH_HEADER")

if echo "$GET_ROLES_RESPONSE" | grep -q "roles"; then
    echo -e "${GREEN}✓ Roles retrieved${NC}"
else
    echo -e "${RED}✗ Failed to get roles${NC}"
fi

echo ""

# Test Permissions
echo -e "${BLUE}=== TESTING PERMISSIONS ===${NC}"

echo -e "${YELLOW}6. Creating permission...${NC}"
CREATE_PERM_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/permissions" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d "{\"application_id\": $APP_ID, \"name\": \"Test Permission\", \"slug\": \"test-permission\", \"description\": \"A test permission\"}")

if echo "$CREATE_PERM_RESPONSE" | grep -q "Permission created successfully"; then
    echo -e "${GREEN}✓ Permission created${NC}"
    PERM_ID=$(echo "$CREATE_PERM_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    echo "Permission ID: $PERM_ID"
else
    echo -e "${RED}✗ Failed to create permission${NC}"
    echo "$CREATE_PERM_RESPONSE"
fi

echo ""

echo -e "${YELLOW}7. Getting permissions...${NC}"
GET_PERMS_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/permissions?application_id=$APP_ID" \
  -H "$AUTH_HEADER")

if echo "$GET_PERMS_RESPONSE" | grep -q "permissions"; then
    echo -e "${GREEN}✓ Permissions retrieved${NC}"
else
    echo -e "${RED}✗ Failed to get permissions${NC}"
fi

echo ""

# Test Role-Permission relationship
echo -e "${BLUE}=== TESTING ROLE-PERMISSION RELATIONSHIP ===${NC}"

echo -e "${YELLOW}8. Assigning permission to role...${NC}"
ASSIGN_PERM_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/roles/$ROLE_ID/permissions" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d "{\"permission_id\": $PERM_ID}")

if echo "$ASSIGN_PERM_RESPONSE" | grep -q "Permission assigned to role successfully"; then
    echo -e "${GREEN}✓ Permission assigned to role${NC}"
else
    echo -e "${RED}✗ Failed to assign permission to role${NC}"
    echo "$ASSIGN_PERM_RESPONSE"
fi

echo ""

echo -e "${YELLOW}9. Getting role permissions...${NC}"
GET_ROLE_PERMS_RESPONSE=$(curl -s -X GET "$BASE_URL/admin/roles/$ROLE_ID/permissions" \
  -H "$AUTH_HEADER")

if echo "$GET_ROLE_PERMS_RESPONSE" | grep -q "permissions"; then
    echo -e "${GREEN}✓ Role permissions retrieved${NC}"
else
    echo -e "${RED}✗ Failed to get role permissions${NC}"
fi

echo ""

# Test User-Role relationship
echo -e "${BLUE}=== TESTING USER-ROLE RELATIONSHIP ===${NC}"

echo -e "${YELLOW}10. Assigning role to user...${NC}"
ASSIGN_ROLE_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/users/1/roles" \
  -H "Content-Type: application/json" \
  -H "$AUTH_HEADER" \
  -d "{\"role_id\": $ROLE_ID}")

if echo "$ASSIGN_ROLE_RESPONSE" | grep -q "Role assigned to user successfully"; then
    echo -e "${GREEN}✓ Role assigned to user${NC}"
else
    echo -e "${RED}✗ Failed to assign role to user${NC}"
    echo "$ASSIGN_ROLE_RESPONSE"
fi

echo ""

echo -e "${YELLOW}11. Testing login with new role...${NC}"
NEW_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "joao@email.com", "password": "123456"}')

if echo "$NEW_LOGIN_RESPONSE" | grep -q "test-role"; then
    echo -e "${GREEN}✓ User now has the new role in JWT${NC}"
else
    echo -e "${YELLOW}! New role might not be showing in JWT (check permissions)${NC}"
fi

echo ""
echo -e "${GREEN}=== Admin API tests completed ===${NC}"

echo ""
echo -e "${BLUE}Available Admin Endpoints:${NC}"
echo "Applications:"
echo "  POST   /admin/applications"
echo "  GET    /admin/applications"
echo "  GET    /admin/applications/:id"
echo "  PUT    /admin/applications/:id"
echo "  DELETE /admin/applications/:id"
echo ""
echo "Roles:"
echo "  POST   /admin/roles"
echo "  GET    /admin/roles[?application_id=X]"
echo "  GET    /admin/roles/:id"
echo "  PUT    /admin/roles/:id"
echo "  DELETE /admin/roles/:id"
echo ""
echo "Permissions:"
echo "  POST   /admin/permissions"
echo "  GET    /admin/permissions[?application_id=X]"
echo "  GET    /admin/permissions/:id"
echo "  PUT    /admin/permissions/:id"
echo "  DELETE /admin/permissions/:id"
echo ""
echo "Role-Permission Relationships:"
echo "  POST   /admin/roles/:id/permissions"
echo "  GET    /admin/roles/:id/permissions"
echo "  DELETE /admin/roles/:id/permissions/:permissionId"
echo ""
echo "User-Role Relationships:"
echo "  POST   /admin/users/:id/roles"
echo "  DELETE /admin/users/:id/roles/:roleId"
