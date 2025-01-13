package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	AlchemyAPIKey string
	KafkaBrokers  []string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config := &Config{
		AlchemyAPIKey: envFile["ALCHEMY_API_KEY"],
		KafkaBrokers:  []string{envFile["KAFKA_BROKER_1"], envFile["KAFKA_BROKER_2"]},
	}

	return config, nil
}
