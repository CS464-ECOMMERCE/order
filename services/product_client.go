package services

import (
	"context"
	"fmt"
	"order/configs"
	"time"

	pb "order/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProductClient is a client for the product service
type ProductClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
	config configs.EnvConfig
}

var productClient *ProductClient

// GetProductClient returns a singleton instance of the product client
func GetProductClient() *ProductClient {
	if productClient == nil {
		config := configs.GetEnvConfig()

		// Set up gRPC connection
		addr := fmt.Sprintf("%s:%s", config.ProductServiceHost, config.ProductServicePort)
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to product service: %v", err))
		}

		productClient = &ProductClient{
			conn:   conn,
			client: pb.NewProductServiceClient(conn),
			config: config,
		}
	}

	return productClient
}

// Close closes the gRPC connection
func (pc *ProductClient) Close() error {
	if pc.conn != nil {
		return pc.conn.Close()
	}
	return nil
}

// GetProduct retrieves a product from the product service
func (pc *ProductClient) GetProduct(id uint64) (*pb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return pc.client.GetProduct(ctx, &pb.GetProductRequest{Id: id})
}

// ValidateInventory checks if there is sufficient inventory for the requested quantity
func (pc *ProductClient) ValidateInventory(productID, requestedQuantity uint64) (*pb.Product, error) {
	product, err := pc.GetProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product.Inventory < requestedQuantity {
		return nil, fmt.Errorf("insufficient inventory for product %d: requested %d, available %d",
			productID, requestedQuantity, product.Inventory)
	}

	return product, nil
}

// UpdateInventory updates the inventory of a product
func (pc *ProductClient) UpdateInventory(productID, newQuantity uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product, err := pc.GetProduct(productID)
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	product.Inventory = newQuantity

	// Convert the Product into an UpdateProductRequest for the API call
	updateRequest := &pb.UpdateProductRequest{
		Id:              product.Id,
		Name:            product.Name,
		Price:           product.Price,
		Inventory:       newQuantity,
		Description:     product.Description,
		Images:          product.Images,
		StripePriceId:   product.StripePriceId,
		StripeProductId: product.StripeProductId,
		MerchantId:      product.MerchantId,
	}

	_, err = pc.client.UpdateProduct(ctx, updateRequest)
	if err != nil {
		return fmt.Errorf("failed to update product inventory: %w", err)
	}

	return nil
}
