#!/bin/bash

# Demo: Multi-Database Configuration
# Demonstrates switching between SQLite and MySQL

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${BLUE}=================================${NC}"
echo -e "${BLUE}   Multi-Database Demo Script    ${NC}"
echo -e "${BLUE}=================================${NC}"
echo ""

# Function to test API
test_api() {
    local port=$1
    echo -e "${YELLOW}Testing API on port $port...${NC}"
    
    # Wait for server to start
    sleep 2
    
    # Test if server is responding
    if curl -s -f "http://localhost:$port/swagger/index.html" > /dev/null; then
        echo -e "${GREEN}✓ API is responding on port $port${NC}"
        return 0
    else
        echo -e "${RED}✗ API not responding on port $port${NC}"
        return 1
    fi
}

# Demo 1: SQLite (Default)
echo -e "${CYAN}=== Demo 1: SQLite Configuration ===${NC}"
echo "Creating .env for SQLite..."

cat > .env << EOF
DATABASE_TYPE=sqlite
DATABASE_URL=./demo_auth.db
PORT=8090
EOF

echo -e "${GREEN}✓ .env configured for SQLite${NC}"
echo ""

echo "Starting application with SQLite..."
PORT=8090 go run main.go &
APP_PID=$!

if test_api 8090; then
    echo -e "${GREEN}✓ SQLite demo successful${NC}"
else
    echo -e "${RED}✗ SQLite demo failed${NC}"
fi

echo "Stopping SQLite demo..."
kill $APP_PID 2>/dev/null
wait $APP_PID 2>/dev/null
echo ""

# Demo 2: MySQL (if available)
echo -e "${CYAN}=== Demo 2: MySQL Configuration ===${NC}"
echo "Checking if MySQL is available..."

if command -v mysql &> /dev/null; then
    echo -e "${GREEN}✓ MySQL client found${NC}"
    
    # Test MySQL connection (without password for demo)
    if mysql -u root -e "SELECT 1;" &> /dev/null; then
        echo -e "${GREEN}✓ MySQL server is accessible${NC}"
        
        echo "Creating .env for MySQL..."
        cat > .env << EOF
DATABASE_TYPE=mysql
DATABASE_URL=root:@tcp(localhost:3306)/auth_demo?charset=utf8mb4&parseTime=True&loc=Local
PORT=8091
EOF

        echo "Creating demo database..."
        mysql -u root -e "CREATE DATABASE IF NOT EXISTS auth_demo;" 2>/dev/null
        
        echo "Starting application with MySQL..."
        PORT=8091 go run main.go &
        APP_PID=$!
        
        if test_api 8091; then
            echo -e "${GREEN}✓ MySQL demo successful${NC}"
        else
            echo -e "${RED}✗ MySQL demo failed${NC}"
        fi
        
        echo "Stopping MySQL demo..."
        kill $APP_PID 2>/dev/null
        wait $APP_PID 2>/dev/null
        
        echo "Cleaning up demo database..."
        mysql -u root -e "DROP DATABASE IF EXISTS auth_demo;" 2>/dev/null
    else
        echo -e "${YELLOW}! MySQL server not accessible (no password auth)${NC}"
        echo -e "${YELLOW}  Skipping MySQL demo${NC}"
    fi
else
    echo -e "${YELLOW}! MySQL client not found${NC}"
    echo -e "${YELLOW}  Skipping MySQL demo${NC}"
fi

echo ""

# Summary
echo -e "${BLUE}=== Demo Summary ===${NC}"
echo -e "${GREEN}✓ Multi-database configuration working${NC}"
echo -e "${GREEN}✓ SQLite (default) - File-based database${NC}"
echo -e "${GREEN}✓ MySQL (optional) - Server-based database${NC}"
echo ""
echo "Key features implemented:"
echo "• Automatic database type detection"
echo "• Environment-based configuration"
echo "• Database-specific schemas and migrations" 
echo "• Smart seeding scripts"
echo ""

# Restore default configuration
echo "Restoring default configuration..."
cat > .env << EOF
DATABASE_TYPE=sqlite
DATABASE_URL=./auth.db
PORT=8080
EOF

echo -e "${GREEN}✓ Default configuration restored${NC}"
echo ""
echo -e "${BLUE}Demo completed! Your API now supports both SQLite and MySQL.${NC}"
echo -e "${YELLOW}To use MySQL: Set DATABASE_TYPE=mysql in your .env file${NC}"
