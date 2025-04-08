package grpc

import (
	"fmt"
	"log"
	"net"
	"order/configs"
	"order/controllers"
	pb "order/proto"
	"order/services"

	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Init initializes and starts the gRPC server
func ServerInit() {
	config := configs.GetEnvConfig()
	address := fmt.Sprintf(":%s", config.GrpcPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("OrderService", grpc_health_v1.HealthCheckResponse_SERVING)

	orderService := services.NewOrderService()

	pb.RegisterOrderServiceServer(s, controllers.NewOrderController(orderService))

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
