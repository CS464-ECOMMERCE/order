package grpc

import (
	"log"
	"order/configs"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiServer struct {
	ProductServiceConn *grpc.ClientConn
	CartServiceConn    *grpc.ClientConn
}

func mustConnGRPC(conn **grpc.ClientConn, addr string) {
	var err error
	*conn, err = grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy": "round_robin",
			"methodConfig": [{
				"name": [{"service": ""}],
				"retryPolicy": {
					"maxAttempts": 5,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2.0,
					"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`),
	)
	log.Printf("grpc: connecting to %s", addr)

	if err != nil {
		log.Printf("grpc: failed to connect %s", addr)
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}

var ApiServerInstance *ApiServer

func ClientInit() {
	ApiServerInstance = &ApiServer{}
	configs := configs.GetEnvConfig()
	mustConnGRPC(&ApiServerInstance.ProductServiceConn, configs.ProductServiceAddr)
	mustConnGRPC(&ApiServerInstance.CartServiceConn, configs.CartServiceAddr)
}