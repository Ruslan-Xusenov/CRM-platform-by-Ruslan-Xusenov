# ==============================================================
# Omnichannel CRM & WebRTC PBX Platform — Makefile
# ==============================================================

.PHONY: help dev down build migrate test lint logs clean

# Default target
help: ## Show this help message
	@echo "╔══════════════════════════════════════════════════╗"
	@echo "║   Omnichannel CRM & WebRTC PBX Platform         ║"
	@echo "╚══════════════════════════════════════════════════╝"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ─────────────────────────────────────────────

dev: ## Start all services in development mode
	docker compose -f docker-compose.yml up --build -d
	@echo "\n✅ All services started!"
	@echo "   Frontend:  http://localhost:3000"
	@echo "   Backend:   http://localhost:8080"
	@echo "   RabbitMQ:  http://localhost:15672"
	@echo "   MinIO:     http://localhost:9001"
	@echo "   Grafana:   http://localhost:3001"

down: ## Stop all services
	docker compose down

restart: ## Restart all services
	docker compose down && docker compose up --build -d

build: ## Build all Docker images
	docker compose build

# ─── Database ────────────────────────────────────────────────

migrate: ## Run database migrations
	docker compose exec backend /app/server migrate up

migrate-down: ## Rollback last migration
	docker compose exec backend /app/server migrate down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=my_migration)
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=my_migration"; exit 1; fi
	touch backend/migrations/$$(date +%Y%m%d%H%M%S)_$(NAME).up.sql
	touch backend/migrations/$$(date +%Y%m%d%H%M%S)_$(NAME).down.sql
	@echo "✅ Created migration: $(NAME)"

# ─── Backend ─────────────────────────────────────────────────

backend-dev: ## Run backend locally (outside Docker)
	cd backend && go run ./cmd/server/main.go

backend-test: ## Run backend tests
	cd backend && go test -v -race -cover ./...

backend-lint: ## Lint backend code
	cd backend && golangci-lint run ./...

# ─── Frontend ────────────────────────────────────────────────

frontend-dev: ## Run frontend locally (outside Docker)
	cd frontend && npm run dev

frontend-test: ## Run frontend tests
	cd frontend && npm test

frontend-lint: ## Lint frontend code
	cd frontend && npm run lint

# ─── Testing ─────────────────────────────────────────────────

test: backend-test frontend-test ## Run all tests

lint: backend-lint frontend-lint ## Run all linters

# ─── Monitoring ──────────────────────────────────────────────

logs: ## View all service logs (follow mode)
	docker compose logs -f

logs-backend: ## View backend logs
	docker compose logs -f backend

logs-asterisk: ## View Asterisk logs
	docker compose logs -f asterisk

# ─── Utilities ───────────────────────────────────────────────

clean: ## Remove all containers, volumes, and images
	docker compose down -v --rmi local
	@echo "✅ Cleaned up all containers and volumes"

psql: ## Connect to PostgreSQL shell
	docker compose exec postgres psql -U $${POSTGRES_USER:-crm_admin} -d $${POSTGRES_DB:-crm_platform}

redis-cli: ## Connect to Redis CLI
	docker compose exec redis redis-cli -a $${REDIS_PASSWORD:-change-me}

asterisk-cli: ## Connect to Asterisk CLI
	docker compose exec asterisk asterisk -rvvv
