.PHONY: help build run test lint migrate-up migrate-down migrate-create docker-build docker-up docker-down seed clean

# Variables
APP_NAME=tiny-school-hub-api
MAIN_PATH=./cmd/api
BUILD_DIR=./bin
MIGRATIONS_DIR=./migrations
DATABASE_URL?=postgres://tinyschool:tinyschool@localhost:5432/tinyschoolhub?sslmode=disable

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run --timeout 5m

lint-fix: ## Run linter with auto-fix
	@golangci-lint run --fix --timeout 5m

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

migrate-force: ## Force migration version
	@read -p "Enter version: " version; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $$version

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then echo "NAME is required. Usage: make migrate-create NAME=create_users_table"; exit 1; fi
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	@docker-compose up -d

docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-logs: ## View Docker Compose logs
	@docker-compose logs -f

swagger: ## Open Swagger UI in browser
	@echo "Opening Swagger UI..."
	@open http://localhost:8081 || xdg-open http://localhost:8081 || echo "Open http://localhost:8081 in your browser"

redoc: ## Open ReDoc documentation in browser
	@echo "Opening ReDoc..."
	@open http://localhost:8082 || xdg-open http://localhost:8082 || echo "Open http://localhost:8082 in your browser"

docs: ## Open all API documentation
	@echo "Opening API documentation..."
	@make swagger
	@sleep 1
	@make redoc

seed: ## Seed the database with test data
	@echo "Seeding database..."
	@go run ./scripts/seed/main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt coverage.html
	@go clean

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

install-hooks: ## Install git pre-commit hooks
	@echo "Installing git hooks..."
	@chmod +x scripts/pre-commit.sh
	@ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit
	@echo "✓ Pre-commit hook installed successfully"
	@echo "To bypass hooks use: git commit --no-verify"

uninstall-hooks: ## Uninstall git pre-commit hooks
	@echo "Uninstalling git hooks..."
	@rm -f .git/hooks/pre-commit
	@echo "✓ Pre-commit hook uninstalled"

pre-commit: ## Run pre-commit checks manually
	@echo "Running pre-commit checks..."
	@./scripts/pre-commit.sh

install-pre-commit-framework: ## Install pre-commit framework (Python-based)
	@echo "Installing pre-commit framework..."
	@pip install pre-commit
	@pre-commit install
	@echo "✓ Pre-commit framework installed"
	@echo "Run 'pre-commit run --all-files' to check all files"

.DEFAULT_GOAL := help
