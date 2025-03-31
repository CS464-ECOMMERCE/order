package services

import (
	"context"
	pb "order/proto"
	"time"

	"google.golang.org/grpc"
)

type CartService struct {
	client pb.CartServiceClient
}

func NewCartService(conn *grpc.ClientConn) *CartService {
	return &CartService{
		client: pb.NewCartServiceClient(conn),
	}
}

// GetCart retrieves cart items from the cart service
func (c *CartService) GetCart(sessionId string) (*pb.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetCart(ctx, &pb.GetCartRequest{SessionId: sessionId})
	return resp, err
}

func (c *CartService) DeleteCart(sessionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.EmptyCart(ctx, &pb.EmptyCartRequest{
		SessionId: sessionId,
	})
	return err
}
