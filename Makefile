.PHONY: dev setup gen-swagger gen-mocks

# Use variables to define important paths and versions
AIR_VERSION = latest
MOCKERY_VERSION = v2.32.0
SWAGGER_VERSION = v1.16.1
SWAGGER_IMAGE = ghcr.io/swaggo/swag:$(SWAGGER_VERSION)
MOCKERY_IMAGE = vektra/mockery:$(MOCKERY_VERSION)


# Define targets and their dependencies
dev: setup
	air -c .air.toml

setup:
	@go install github.com/cosmtrek/air@$(AIR_VERSION)
	@go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	@go install github.com/volatiletech/sqlboiler/v4@latest
	@go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

gen-swagger:
	@docker run --rm -v $(shell pwd):/app -w /app $(SWAGGER_IMAGE) /root/swag init -g ./cmd/server/main.go

gen-mocks:
	mockery --all --keeptree --output ./mocks
	# @docker run --rm -v $(shell pwd):/app -w /app $(MOCKERY_IMAGE) --all --keeptree --output ./mocks

pg-start-dev:
	docker-compose up -d db

pg-start-test:
	docker-compose up -d db-test

pg-stop-dev:
	docker-compose down -t 0 db

pg-stop-test:
	docker-compose down -t 0 db-test

pg-migrate-up:
	sql-migrate up -env=development

pg-migrate-down:
	sql-migrate down -env=development

test: pg-start-test
	sql-migrate down -env="test" -limit=0
	sql-migrate up -env="test"
	ENV=test go test -v ./... --cover

gen-models: pg-migrate-up
	sqlboiler psql

.PHONY: help

help:
	@echo "Available targets:"
	@echo "  setup              Install required dependencies"
	@echo "  dev                Run the development server"
	@echo "  gen-swagger        Generate Swagger documentation"
	@echo "  gen-mocks          Generate mock interfaces"
	@echo "  pg-start-dev       Start the development database container"
	@echo "  pg-start-test      Start the testing database container"
	@echo "  pg-stop-dev        Stop the development database container"
	@echo "  pg-stop-test       Stop the testing database container"
	@echo "  pg-migrate-up      Apply pending migrations"
	@echo "  pg-migrate-down    Rollback the last migration"
	@echo "  gen-models    Generate models using sqlboiler"
	@echo "  test               Start the testing database container, run tests, and stop the container"
