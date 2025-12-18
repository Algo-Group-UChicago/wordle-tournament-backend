package storage

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

// initializeDynamo initializes the DynamoDB client from environment variables.
// It reads DYNAMODB_ENDPOINT (required) and AWS_REGION (defaults to us-east-1),
// then creates a client configured to use the specified endpoint.
func initializeDynamo() {
	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		panic("DYNAMODB_ENDPOINT environment variable must be set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	ctx := context.Background()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	dynamoClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}
