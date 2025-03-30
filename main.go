package main

import (
	"fmt"
	"log"
	"order/configs"
	"order/controllers"
	"order/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	configs.InitEnv()
	config := configs.GetEnvConfig()

	// Initialize database
	storage.Initialize()
	defer func() {
		if err := storage.GetInstance().Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Create router
	router := gin.Default()

	// Set up controllers
	orderController := controllers.NewOrderController()
	orderController.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%s", config.OrderPort)
	log.Printf("Order service starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
