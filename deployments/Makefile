# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go_auth
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOCMD) run main.go

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Test coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Reset database
reset-db:
	rm -f auth.db
	@echo "Database reset. Run 'make run' to recreate it."

# Seed database with sample data
seed:
	./scripts/seed_db.sh

# Test API endpoints
test-api:
	./scripts/test_api.sh

# Test Admin API endpoints
test-admin:
	./scripts/test_admin_api.sh

# Build for Linux
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v .

# Docker build
docker-build:
	docker build -t $(BINARY_NAME) .

# Show help
help:
	@echo ""
	@echo "Usage:"
	@echo "  make build        Build the application"
	@echo "  make run          Run the application"
	@echo "  make test         Run tests"
	@echo "  make test-api     Test API endpoints"
	@echo "  make test-admin   Test admin API endpoints"
	@echo "  make clean        Clean build files"
	@echo "  make deps         Install dependencies"
	@echo "  make reset-db     Reset database"
	@echo "  make seed         Seed database with sample data"
	@echo "  make help         Show this help message"
	@echo ""

.PHONY: build clean test test-coverage deps run reset-db seed test-api test-admin build-linux docker-build help
