package config

import "os"

type Config struct {
	Port             string
	Environment      string
	AWSRegion        string
	DynamoDBEndpoint string
}

func Load() *Config {
	return &Config{
		Port:             getEnvOrDefault("PORT", "8080"),
		Environment:      getEnvOrDefault("ENV", "dev"),
		AWSRegion:        getEnvOrDefault("AWS_REGION", "us-east-1"),
		DynamoDBEndpoint: getEnvOrDefault("DYNAMODB_ENDPOINT", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
