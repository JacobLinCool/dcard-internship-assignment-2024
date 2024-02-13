package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initDatabaseConnection()

	port := os.Getenv("AD_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	ttl := os.Getenv("AD_SERVICE_CACHE_TTL")
	if ttl == "" {
		ttl = "1"
	}
	ttlInt, err := strconv.Atoi(ttl)
	if err != nil {
		log.Fatalf("Failed to parse cache TTL: %v", err)
	}

	go cacheUpdater(time.Second * time.Duration(ttlInt))

	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/ad", createAdHandler)
			v1.GET("/ad", listAdsHandler)
		}
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
