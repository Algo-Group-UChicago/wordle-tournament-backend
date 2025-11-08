# Wordle Tournament Backend

A Go HTTP API for running Wordle bot tournaments.

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
