.PHONY: dev dev-bot dev-web build build-bot build-web test lint clean

ENV_FILE := .env

# Load .env if it exists
ifneq (,$(wildcard $(ENV_FILE)))
  include $(ENV_FILE)
  export
endif

# --- Development ---

dev: ## Start bot and web dev server in parallel
	$(MAKE) dev-bot &
	$(MAKE) dev-web &
	wait

dev-bot: build-bot ## Build and run the Go bot with .env
	cd bot && ./server

dev-web: ## Start the Vite dev server
	cd web && pnpm dev

# --- Build ---

build: build-bot build-web ## Build both bot and web

build-bot: ## Build the Go binary
	$(MAKE) -C bot build

build-web: ## Build the Vite frontend
	cd web && pnpm build

# --- Test ---

test: ## Run all tests (Go + web)
	$(MAKE) -C bot test
	cd web && pnpm typecheck

# --- Lint ---

lint: ## Lint both bot and web
	$(MAKE) -C bot lint
	cd web && pnpm lint

# --- Clean ---

clean: ## Clean build artifacts
	$(MAKE) -C bot clean
	rm -rf web/dist

# --- Docker ---

docker: ## Build docker image
	docker compose build

docker-up: ## Start with docker compose
	docker compose up

docker-down: ## Stop docker compose
	docker compose down

# --- Help ---

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
