.PHONY: help run build test migrate-up migrate-down migrate-create migrate-force migrate-status clean docker-build docker-up docker-up-d docker-down docker-logs docker-clean docker-restart

BINARY_NAME=media-content-api
MIGRATION_DIR=database/migrations

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run:
	go run main.go

build:
	go build -o $(BINARY_NAME) main.go

test:
	go test -v ./...

migrate-up:
	docker run --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations --network host migrate/migrate -path=/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	docker run --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations --network host migrate/migrate -path=/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

migrate-create:
	docker run --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations migrate/migrate create -ext sql -dir /migrations -seq $(NAME)

migrate-status:
	docker run --rm -v $(shell pwd)/$(MIGRATION_DIR):/migrations --network host migrate/migrate -path=/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" version

clean:
	rm -f $(BINARY_NAME)
	go clean

docker-build:
	docker-compose build

docker-up:
	docker-compose up

docker-up-d:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-restart:
	docker-compose restart