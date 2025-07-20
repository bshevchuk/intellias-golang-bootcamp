.DEFAULT_GOAL := help
help: # adjust number in substring "-30s\" to adjust tabulation from left border
	@grep -E '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#.PHONY: cli-build
#cli-build: ## Build CLI tool
#	go build ./cmd/cli
#
#.PHOHY: cli-run
#cli-run: ## Run CLI tool
#	go run ./cmd/cli

.PHONY: cli-build
server-build: ## Build API server
	go build ./cmd/api-server

.PHOHY: cli-run
server-run: ## Run API server
	go run ./cmd/api-server


.PHONY: infra-up
infra-up: ## Start local infrastructure
	docker compose up -d

.PHONY: infra-down
infra-down: ## Stop local infrastructure
	docker compose down

.PHONY: infra-amnesia
infra-amnesia: ## !!! Stop local infrastructure with cleaning volumes (pgdata)
	docker compose down -v

.PHONY: db-schema-migrations
db-schema-migrations: ## Run database schema migrations
	docker compose run db_schema_migrations
