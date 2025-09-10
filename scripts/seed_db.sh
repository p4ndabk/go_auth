#!/bin/bash

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

DB_FILE="auth.db"

echo -e "${YELLOW}=== Seeding Database with Sample Data ===${NC}"

# Check if database exists
if [ ! -f "$DB_FILE" ]; then
    echo "Database file not found. Please run the application first to create the database."
    exit 1
fi

# Apply seed data
sqlite3 "$DB_FILE" < migrations/002_seed_data.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Sample data inserted successfully${NC}"
    
    echo ""
    echo "Sample data includes:"
    echo "- 2 Applications (main-app, admin-app)"
    echo "- 3 Roles (admin, user, moderator)"
    echo "- 7 Permissions (read_users, write_users, etc.)"
    echo "- Role-Permission associations"
    echo ""
    echo "To assign roles to users, you can run SQL commands like:"
    echo "INSERT INTO user_roles (user_id, role_id) VALUES (1, 1); -- Assign admin role to user 1"
    echo ""
    echo "You can also check the data with:"
    echo "sqlite3 $DB_FILE 'SELECT * FROM roles;'"
    echo "sqlite3 $DB_FILE 'SELECT * FROM permissions;'"
else
    echo "Error inserting sample data"
    exit 1
fi
