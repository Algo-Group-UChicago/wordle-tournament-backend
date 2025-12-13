package config

import "os"

type Config struct {
	Port             string
	Region           string
	DynamoDBEndpoint string
}

func Load() Config {
	return Config{
		Port:             getEnv("PORT", "8080"),
		Region:           getEnv("AWS_REGION", "us-east-1"),
		DynamoDBEndpoint: getEnv("DYNAMODB_ENDPOINT", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
