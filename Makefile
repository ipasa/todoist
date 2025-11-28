.PHONY: help setup build run-backend run-frontend test clean migrate docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Setup development environment
	@echo "Setting up development environment..."
	cd frontend && npm install
	@echo "Setup complete!"

build: ## Build all services
	@echo "Building all services..."
	cd backend/auth-service && go build -o ../../bin/auth-service ./cmd
	cd backend/task-service && go build -o ../../bin/task-service ./cmd
	cd backend/project-service && go build -o ../../bin/project-service ./cmd
	cd backend/notification-service && go build -o ../../bin/notification-service ./cmd
	cd backend/websocket-gateway && go build -o ../../bin/websocket-gateway ./cmd
	cd backend/api-gateway && go build -o ../../bin/api-gateway ./cmd
	@echo "Backend build complete!"
	cd frontend && npm run build
	@echo "Frontend build complete!"
	@echo "All services built successfully!"

run-backend: ## Run all backend services
	@echo "Starting backend services..."
	docker-compose -f infrastructure/docker/docker-compose.yml up postgres rabbitmq redis auth-service task-service project-service notification-service websocket-gateway api-gateway mailhog

run-frontend: ## Run frontend development server
	cd frontend && npm run dev

run-frontend-docker: ## Run frontend in Docker
	@echo "Starting frontend in Docker..."
	docker-compose -f infrastructure/docker/docker-compose.yml up frontend

test: ## Run tests for all services
	@echo "Running tests..."
	cd backend/auth-service && go test ./...
	cd backend/task-service && go test ./...
	cd backend/project-service && go test ./...
	cd backend/notification-service && go test ./...
	cd backend/websocket-gateway && go test ./...
	cd backend/api-gateway && go test ./...
	cd frontend && npm test

migrate: ## Run database migrations
	@echo "Running migrations..."
	cd infrastructure/scripts && ./migrate.sh

docker-up: ## Start Docker containers
	docker-compose -f infrastructure/docker/docker-compose.yml up -d

docker-up-build: ## Start Docker containers with build
	docker-compose -f infrastructure/docker/docker-compose.yml up -d --build

docker-down: ## Stop Docker containers
	docker-compose -f infrastructure/docker/docker-compose.yml down

clean: ## Clean build artifacts
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

lint: ## Run linters
	@echo "Running linters..."
	cd backend/auth-service && golangci-lint run
	cd backend/task-service && golangci-lint run
	cd backend/project-service && golangci-lint run
	cd backend/notification-service && golangci-lint run
	cd backend/websocket-gateway && golangci-lint run
	cd backend/api-gateway && golangci-lint run
	cd frontend && npm run lint
