cnf ?= ./deployments/.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

server: ## Run application server
	go run ./cmd/api/main.go

migrate-up: ## Run migrations
	docker run -v $(PWD)/deployments/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database postgres://postgres:$(POSTGRES_PASSWORD)@localhost:5434/postgres?sslmode=disable up

migrate-down: ## Rollback migrations
	docker run -v $(PWD)/deployments/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database postgres://postgres:$(POSTGRES_PASSWORD)@localhost:5434/postgres?sslmode=disable down -all

gen-oas: ## Generate OpenAPI from Postman
	p2o ./api/postman_collection.json -f ./api/oas.yml
