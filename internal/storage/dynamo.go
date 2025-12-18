package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"wordle-tournament-backend/internal/config"
)

var (
	dynamoOnce   sync.Once
	dynamoClient *dynamodb.Client
)

// getDynamoClient returns the shared DynamoDB client, initializing it
// lazily on the first call using initializeDynamo. Subsequent calls return
// the same client instance.
func getDynamoClient() *dynamodb.Client {
	dynamoOnce.Do(initializeDynamo)
	return dynamoClient
}

// initializeDynamo initializes the DynamoDB client from configuration.
// It uses the DynamoDBEndpoint and Region from the config package,
// then creates a client configured to use the specified endpoint.
func initializeDynamo() {
	cfg := config.Get()
	if cfg.DynamoDBEndpoint == "" {
		panic("DYNAMODB_ENDPOINT environment variable must be set")
	}

	region := cfg.Region
	ctx := context.Background()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	dynamoClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(cfg.DynamoDBEndpoint)
	})
}
