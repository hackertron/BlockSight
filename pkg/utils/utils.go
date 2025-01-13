package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// LoadEnv loads environment variables from a .env file
func LoadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %v", err)
	}
	defer file.Close()

	vars := make(map[string]string)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&vars); err != nil {
		return fmt.Errorf("failed to decode .env file: %v", err)
	}

	for key, value := range vars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable: %v", err)
		}
	}

	return nil
}

// LogError logs an error message
func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// LogInfo logs an informational message
func LogInfo(message string) {
	log.Printf("Info: %s", message)
}

// LogWarning logs a warning message
func LogWarning(message string) {
	log.Printf("Warning: %s", message)
}
