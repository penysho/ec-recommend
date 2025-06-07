# EC Recommend System Makefile

# Docker and Database
.PHONY: db-up db-down db-restart db-setup db-seed db-reset db-connect

# Variables
DB_HOST=localhost
DB_PORT=5436
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DOCKER_COMPOSE=docker compose

# Database operations
db-up: ## Start the database
	$(DOCKER_COMPOSE) up -d backend-db

db-down: ## Stop the database
	$(DOCKER_COMPOSE) down

db-restart: ## Restart the database
	$(DOCKER_COMPOSE) restart backend-db

db-setup: ## Create database schema
	@echo "Setting up database schema..."
	@docker exec -i ec_recommend-db psql -U $(DB_USER) -d $(DB_NAME) < db/schema.sql

db-seed: ## Seed the database with sample data
	@echo "Seeding database with sample data..."
	@docker exec -i ec_recommend-db psql -U $(DB_USER) -d $(DB_NAME) < db/sample_data.sql

db-reset: db-down ## Reset database (drop, recreate, setup, seed)
	@echo "Resetting database..."
	$(DOCKER_COMPOSE) up -d backend-db
	@sleep 3
	@docker exec -i ec_recommend-db psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@make db-setup
	@make db-seed
	@echo "Database reset complete!"

db-connect: ## Connect to database via psql
	@docker exec -it ec_recommend-db psql -U $(DB_USER) -d $(DB_NAME)

# Application operations
.PHONY: build run test clean dev

build: ## Build the application
	go build -o ec-recommend cmd/main.go

run: ## Run the application
	go run cmd/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -f ec-recommend

dev: db-up ## Start development environment
	@echo "Starting development environment..."
	@sleep 3
	@make run

# Recommendation system operations
.PHONY: recommend-setup recommend-test

recommend-setup: db-up db-setup db-seed ## Setup complete recommendation system
	@echo "Recommendation system setup complete!"
	@echo "Database is running on port $(DB_PORT)"
	@echo "Use 'make db-connect' to connect to the database"

recommend-test: ## Test recommendation queries
	@echo "Testing recommendation queries..."
	@docker exec -i ec_recommend-db psql -U $(DB_USER) -d $(DB_NAME) -c "\
		SELECT 'Customer Purchase Summary:' as test; \
		SELECT customer_id, email, total_orders, total_spent, unique_products_bought FROM customer_purchase_summary LIMIT 3; \
		SELECT 'Product Popularity:' as test; \
		SELECT name, category_id, price, order_count, avg_rating FROM product_popularity WHERE order_count > 0 ORDER BY order_count DESC LIMIT 5;"

# Help
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
