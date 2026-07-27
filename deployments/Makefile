# Run from the repo root: make -f deployments/Makefile <target>
GOCMD=go
BINARY_NAME=go_auth

.PHONY: build run migrate rollback seed test test-coverage vet swagger clean docker-build help

build:
	$(GOCMD) build -o bin/api ./cmd/api
	$(GOCMD) build -o bin/migrate ./cmd/migrate
	$(GOCMD) build -o bin/seed ./cmd/seed

run:
	$(GOCMD) run ./cmd/api

migrate:
	$(GOCMD) run ./cmd/migrate

rollback:
	$(GOCMD) run ./cmd/migrate -rollback

seed:
	$(GOCMD) run ./cmd/seed

test:
	$(GOCMD) test -v ./...

test-coverage:
	$(GOCMD) test -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

vet:
	$(GOCMD) vet ./...

swagger:
	$(GOCMD) tool swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

clean:
	rm -rf bin coverage.out

docker-build:
	docker build -f deployments/Dockerfile -t $(BINARY_NAME) .

help:
	@echo ""
	@echo "Usage (from repo root): make -f deployments/Makefile <target>"
	@echo "  build          Build api/migrate/seed binaries into ./bin"
	@echo "  run            Run the API server"
	@echo "  migrate        Apply pending migrations"
	@echo "  rollback       Roll back the last migration"
	@echo "  seed           Seed the database with sample data"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with an HTML coverage report"
	@echo "  vet            Run go vet"
	@echo "  swagger        Regenerate the Swagger docs"
	@echo "  clean          Remove build artifacts"
	@echo "  docker-build   Build the Docker image"
	@echo ""
