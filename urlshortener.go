// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"urlshortener/internal/cache"
	"urlshortener/internal/config"
	"urlshortener/internal/service"
	"urlshortener/internal/storage"

	"github.com/gin-gonic/gin"
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

	// Create an instance of MongoDBService
	var dbClient storage.MongoDBClient = &storage.MongoDBService{}

	// Connect to MongoDB using the provided URI from the config
	connectObservable := dbClient.Connect(cfg.MongoURI)
	connectObservable.ForEach(func(item interface{}) {
		log.Println("Connected to MongoDB")

		// Initialize the URL collection in MongoDB
		service.InitDatabase(dbClient.GetClient(), cfg.MongoDBName, cfg.MongoCollection)
	}, func(err error) {
		log.Fatalf("MongoDB connection error: %v", err)
	}, func() {
		log.Println("MongoDB connection sequence completed")
	})

	// Ensure MongoDB connection is closed when the server stops
	defer func() {
		disconnectObservable := dbClient.Disconnect()
		disconnectObservable.ForEach(func(item interface{}) {
			log.Println(item)
		}, func(err error) {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}, func() {
			log.Println("MongoDB disconnection sequence completed")
		})
	}()

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
