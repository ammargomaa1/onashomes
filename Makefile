.PHONY: help run build clean test migrate-up migrate-down migration deps install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@echo "🚀 Starting application..."
	go run cmd/api/main.go

build: ## Build executables for all platforms
	@echo "🔨 Building executables..."
	./scripts/build.sh

build-linux: ## Build for Linux only
	@echo "🔨 Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/ecommerce-api-linux-amd64 ./cmd/api

build-windows: ## Build for Windows only
	@echo "🔨 Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/ecommerce-api-windows-amd64.exe ./cmd/api

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	go clean

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	go mod download

tidy: ## Tidy dependencies
	@echo "🧹 Tidying dependencies..."
	go mod tidy

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

migration: ## Create a new migration (usage: make migration name=create_products)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migration name=create_products"; \
		exit 1; \
	fi
	./scripts/create_migration.sh $(name)

migrate-down: ## Rollback last migration (usage: make migrate-down steps=1)
	@STEPS=$${steps:-1}; \
	echo "⚠️  Rolling back $$STEPS migration(s)..."; \
	./scripts/migrate_down.sh $$STEPS

dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without auto-reload..."; \
		make run; \
	fi

docker-up: ## Start PostgreSQL with Docker
	@echo "🐳 Starting PostgreSQL..."
	docker-compose up -d

docker-down: ## Stop PostgreSQL
	@echo "🐳 Stopping PostgreSQL..."
	docker-compose down

docker-logs: ## View PostgreSQL logs
	docker-compose logs -f

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin"; \
	fi

setup: ## Initial setup (install deps, create .env)
	@echo "🔧 Setting up project..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✓ Created .env file"; \
	fi
	@make install
	@echo "✓ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env with your database credentials"
	@echo "2. Create database: createdb ecommerce_db"
	@echo "3. Run: make run"

.DEFAULT_GOAL := help
