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
make integration-tests
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
