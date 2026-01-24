#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

echo "Building API..."
docker compose build api

echo "Starting DynamoDB Local, setup-local-db, and API..."
docker compose up -d dynamodb-local setup-local-db api

echo "Waiting for DynamoDB Local (port 8000)..."
timeout=30
counter=0
while ! curl -s http://localhost:8000 > /dev/null 2>&1; do
    if [ "$counter" -ge "$timeout" ]; then
        echo "DynamoDB Local failed to become ready"
        docker compose down
        exit 1
    fi
    sleep 1
    counter=$((counter + 1))
done

echo "Waiting for API (port 8080)..."
counter=0
while ! curl -sf http://localhost:8080/health > /dev/null 2>&1; do
    if [ "$counter" -ge "$timeout" ]; then
        echo "API failed to become ready"
        docker compose down
        exit 1
    fi
    sleep 1
    counter=$((counter + 1))
done

echo "DynamoDB and backend are up."
echo "  DynamoDB: http://localhost:8000"
echo "  API:      http://localhost:8080"
