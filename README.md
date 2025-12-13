# Wordle Tournament Backend

A Go HTTP API for running Wordle bot tournaments.

## Local Development Ports
- **8080** - Wordle API
- **8000** - DynamoDB Local

## Running Locally with Docker Compose

### Start all services (API + DynamoDB):
```bash
docker-compose up --build
```

### Test the API:
```bash
curl http://localhost:8080/health
```

### Initialize env vars
```bash
source scripts/local-env-setup.sh
```

### Verify DynamoDB tables:
```bash
aws dynamodb list-tables \
  --endpoint-url http://localhost:8000
```

### Stop all services:
```bash
docker-compose down
```

### Clean slate (remove persisted data):
```bash
docker-compose down -v
```

## Running Tests

All test commands must be run from the project root.

### Run all unit tests:
```bash
go test ./...
```

### Run tests for a specific package:
```bash
go test ./internal/wordle
```
