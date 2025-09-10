#!/bin/bash

# Smart Database Seed Script
# Automatically detects database type and runs appropriate seed script

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Load environment variables if .env file exists
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Get database type from environment or default to sqlite
DATABASE_TYPE=${DATABASE_TYPE:-"sqlite"}

echo -e "${BLUE}=== Smart Database Seeder ===${NC}"
echo -e "Detected database type: ${YELLOW}$DATABASE_TYPE${NC}"

case "$DATABASE_TYPE" in
    "mysql")
        echo -e "${YELLOW}Using MySQL seeder...${NC}"
        if [ -f "./scripts/seed_db_mysql.sh" ]; then
            ./scripts/seed_db_mysql.sh
        else
            echo -e "${RED}✗ MySQL seed script not found at ./scripts/seed_db_mysql.sh${NC}"
            exit 1
        fi
        ;;
    "sqlite")
        echo -e "${YELLOW}Using SQLite seeder...${NC}"
        if [ -f "./scripts/seed_db.sh" ]; then
            ./scripts/seed_db.sh
        else
            echo -e "${RED}✗ SQLite seed script not found at ./scripts/seed_db.sh${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}✗ Unsupported database type: $DATABASE_TYPE${NC}"
        echo "Supported types: sqlite, mysql"
        exit 1
        ;;
esac

echo -e "${GREEN}✓ Database seeding completed for $DATABASE_TYPE${NC}"
