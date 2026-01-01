.PHONY: help test test-unit test-integration build run clean docker-up docker-down

.DEFAULT_GOAL := help

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk '/^[a-zA-Z_-]+:/ {printf "  %s\n", $$1}' $(MAKEFILE_LIST) | grep -v '^  help$$' | sort

unit-tests:
	go test ./...

integration-tests: clean
	@./scripts/run-integration-tests.sh

build:
	go build -o bin/api ./cmd/api

run:
	@echo "Starting application..."
	@go run ./cmd/api

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -rf bin/
	docker-compose down -v
