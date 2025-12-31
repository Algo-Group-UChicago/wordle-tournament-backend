package config

import (
	"os"
	"sync"
)

type Config struct {
	Port             string
	Region           string
	DynamoDBEndpoint string
	RandomSeed       string
}

var (
	cfg  Config
	once sync.Once
)

// Get returns the application configuration, loading it from environment
// variables on the first call. Subsequent calls return the same cached config.
func Get() Config {
	once.Do(initializeConfig)
	return cfg
}

func initializeConfig() {
	cfg = Config{
		Port:             getEnv("PORT", "8080"),
		Region:           getEnv("AWS_REGION", "us-east-1"),
		DynamoDBEndpoint: getEnv("DYNAMODB_ENDPOINT", ""),
		RandomSeed:       getEnv("RANDOM_SEED", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
