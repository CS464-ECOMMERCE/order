package main

import (
	"fmt"
	"order/configs"
	"order/grpc"
	"order/storage"
)

func main() {
	// Load configuration
	fmt.Println("Initializing Order Service...")
	configs.InitEnv()

	// Initialize database
	db := storage.GetInstance() // init db
	defer db.Close()

	// Start gRPC server
	fmt.Println("Starting gRPC server...")
	grpc.ServerInit()
}
