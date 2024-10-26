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
		urlService := &service.URLServiceImpl{}
		urlService.InitDatabase(dbClient.GetClient(), cfg.MongoDBName, cfg.MongoCollection)

		// Set URLServiceInstance to the initialized URLService
		service.URLServiceInstance = urlService
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

	// Instantiate services
	statService := service.NewURLStatService()
	urlService := service.NewURLShortenerService(statService)

	// Define routes
	router.POST("/shorten", urlService.ShortenURLHandler)
	router.GET("/:id", urlService.RedirectURLHandler)
	router.PATCH("/:id", urlService.ToggleURLStateHandler)
	//router.GET("/stats/:id", statService.GetURLStats)

	/*router.GET("/stats/:id", func(c *gin.Context) {
		shortID := c.Param("id")
		statsObservable := statService.GetURLStats(shortID)

		result := <-statsObservable.Observe()
		if result.E != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get URL stats"})
			return
		}

		stats := result.V.(map[string]interface{})
		c.JSON(http.StatusOK, stats)
	})*/

	// Start the server
	log.Printf("Listening on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
