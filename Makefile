# ==============================================================================
# Development Commands
# ==============================================================================

.PHONY: dev-up
dev-up:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.dev.yml up --build

.PHONY: dev-down
dev-down:
	@echo "Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

# ==============================================================================
# Production Commands
# ==============================================================================

.PHONY: prod-up
prod-up:
	@echo "Starting production environment..."
	docker-compose -f docker-compose.prod.yml up -d --build

.PHONY: prod-down
prod-down:
	@echo "Stopping production environment..."
	docker-compose -f docker-compose.prod.yml down