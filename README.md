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

### Stop all services:
```bash
docker-compose down
```

### Clean slate (remove persisted data):
```bash
docker-compose down -v
```

### Test the API:
```bash
curl http://localhost:8080/health
```

### Initialize env vars
```bash
source scripts/local-env-setup.sh
```

### Sample Call to /start
```bash
curl -X POST http://localhost:8080/start \
  -H "Content-Type: application/json" \
  -d '{"team_id": "TEST"}'
```

### View DynamoDB Entires
```bash
aws dynamodb scan --table-name ActiveRuns --endpoint-url http://localhost:8000 --output json
```

### List DynamoDB tables:
```bash
aws dynamodb list-tables \
  --endpoint-url http://localhost:8000
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

### Run integration tests:

Integration tests require DynamoDB Local to be running. The easiest way is using Make:

```bash
make test-integration
```

This will:
1. Start DynamoDB Local and create necessary tables
2. Set required environment variables
3. Run the integration tests
4. Clean up Docker containers

**Alternative: Run script directly:**
```bash
./scripts/run-integration-tests.sh
```

**Manual execution:**

If you prefer to run tests manually:

1. Start DynamoDB Local:
```bash
make docker-up
# or
docker-compose up -d dynamodb-local setup-local-db
```

2. Set environment variables:
```bash
export DYNAMODB_ENDPOINT=http://localhost:8000
export AWS_REGION=us-east-1
export RANDOM_SEED=1
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
```

3. Run integration tests:
```bash
go test -tags=integration -v ./internal/handlers/...
```

4. Clean up:
```bash
make docker-down
# or
docker-compose down
```

## Makefile Commands

Common commands available via `make`:

```bash
make help              # Show all available commands
make test              # Run unit tests
make test-integration  # Run integration tests
make build             # Build the application
make run               # Run the application locally
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make clean             # Clean build artifacts
```
