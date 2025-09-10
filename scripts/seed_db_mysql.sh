#!/bin/bash

# MySQL Seed Script
# Populates MySQL database with sample data

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database configuration (can be overridden by environment variables)
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"root"}
DB_PASSWORD=${DB_PASSWORD:-""}
DB_NAME=${DB_NAME:-"auth_db"}

echo -e "${YELLOW}=== Seeding MySQL Database with Sample Data ===${NC}"

# Check if mysql client is available
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}✗ MySQL client not found. Please install MySQL client.${NC}"
    exit 1
fi

# Test connection
echo "Testing MySQL connection..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" &> /dev/null; then
    echo -e "${RED}✗ Failed to connect to MySQL database${NC}"
    echo "Please check your MySQL connection parameters:"
    echo "  Host: $DB_HOST"
    echo "  Port: $DB_PORT" 
    echo "  User: $DB_USER"
    echo "  Database: $DB_NAME"
    exit 1
fi

echo -e "${GREEN}✓ MySQL connection successful${NC}"

# Apply schema
echo "Applying MySQL schema..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < migrations/001_initial_schema_mysql.sql; then
    echo -e "${GREEN}✓ Schema applied successfully${NC}"
else
    echo -e "${RED}✗ Failed to apply schema${NC}"
    exit 1
fi

# Apply seed data
echo "Inserting sample data..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < migrations/002_seed_data_mysql.sql; then
    echo -e "${GREEN}✓ Sample data inserted successfully${NC}"
else
    echo -e "${RED}✗ Failed to insert sample data${NC}"
    exit 1
fi

echo ""
echo "Sample data includes:"
echo "- 2 Applications (main-app, admin-app)"
echo "- 4 Roles (admin, user, moderator, admin-panel-admin)"
echo "- 8 Permissions (read_users, write_users, etc.)"
echo "- Role-Permission associations"
echo ""
echo "To assign roles to users, you can run SQL commands like:"
echo "INSERT INTO user_application_role (user_id, application_id, role_id) VALUES (1, 1, 1);"
echo ""
echo "You can also check the data with:"
echo "mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e 'SELECT * FROM roles;'"
echo "mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME -e 'SELECT * FROM permissions;'"
