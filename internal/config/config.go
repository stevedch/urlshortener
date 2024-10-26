package config

import (
	"log"
	"os"
	"strconv"
	"urlshortener/internal/models"
)

// LoadConfig loads configuration from environment variables or default values
func LoadConfig() *models.Config {
	config := &models.Config{
		Port:            getEnv("PORT", "8080"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:     getEnv("MONGO_DB_NAME", "urlshortener"),
		MongoCollection: getEnv("MONGO_COLLECTION", "urls"),
		RedisAddress:    getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""), // No password by default
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
	}

	log.Println("Configuration loaded successfully")
	return config
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s; using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
