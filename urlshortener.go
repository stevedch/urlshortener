package main

import (
	"fmt"
	"log"
	"os"
	"urlshortener/internal/cache"

	"github.com/gin-gonic/gin"
	"urlshortener/internal/config"
	"urlshortener/internal/service"
	"urlshortener/internal/storage" // Import MongoDB connection package
)

func main() {
	// Initialize the Gin server
	router := gin.Default()

	// Set the port from an environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Configuration loaded: %s", cfg)

	// Connect to MongoDB using the provided URI from the config
	client := storage.ConnectMongoDB(cfg.MongoURI)

	// Initialize the URL collection in MongoDB
	service.InitDatabase(client, cfg.MongoDBName, cfg.MongoCollection)

	// Ensure MongoDB connection is closed when the server stops
	defer storage.DisconnectMongoDB()

	// Connect to Redis using configuration details
	cache.InitRedis(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)

	// Define routes
	router.POST("/shorten", func(c *gin.Context) {
		service.ShortenURLHandler(c)
	})

	router.GET("/:id", func(c *gin.Context) {
		service.RedirectURLHandler(c)
	})

	router.PATCH("/:id", func(c *gin.Context) {
		service.ToggleURLStateHandler(c)
	})

	// Start the server
	log.Printf("Listening on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
