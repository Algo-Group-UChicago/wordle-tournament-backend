# Wordle Tournament Backend

A Go HTTP API for running Wordle bot tournaments.

## Project Structure

```
wordle-tournament-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── common/
│   │   └── common.go            # Shared constants and types
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── corpus/
│   │   ├── corpus.go            # Word list management
│   │   ├── corpus.txt           # Valid Wordle guesses
│   │   └── possible_answers.txt # Possible answer words
│   ├── handlers/
│   │   ├── guesses.go           # Guess grading endpoint
│   │   └── health.go            # Health check endpoint
│   ├── server/
│   │   └── server.go            # HTTP server setup and routing
│   └── wordle/
│       ├── grader.go            # Core Wordle grading logic
│       └── validation.go        # Input validation
├── Dockerfile                    # Container definition
└── go.mod                        # Go module dependencies
```

## Running Locally

### Build the Docker image:
```bash
docker build -t wordle-api .
```

### Run the container:
```bash
docker run -d --name wordle-api -p 8080:8080 --env-file .config.local wordle-api
```

### Test the API:
```bash
curl http://localhost:8080/health
```

### Stop the container:
```bash
docker stop wordle-api && docker rm wordle-api
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
