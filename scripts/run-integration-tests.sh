#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_ROOT"

docker-compose build test

docker-compose up -d dynamodb-local setup-local-db

# Wait for DynamoDB Local to be ready before running tests.
# docker-compose up -d returns immediately but doesn't wait for services to be ready.
# This loop ensures DynamoDB Local is actually accepting connections before proceeding.
timeout=30
counter=0
while ! curl -s http://localhost:8000 > /dev/null 2>&1; do
    if [ $counter -ge $timeout ]; then
        docker-compose down
        exit 1
    fi
    sleep 1
    counter=$((counter + 1))
done

sleep 2

docker-compose run --rm test

TEST_EXIT_CODE=$?

docker-compose down

exit $TEST_EXIT_CODE
