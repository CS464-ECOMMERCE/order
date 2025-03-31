package services

import (
	"context"
	"fmt"
	pb "order/proto"
	"time"

	"google.golang.org/grpc"
)

type ProductService struct {
	client pb.ProductServiceClient
}

func NewProductService(conn *grpc.ClientConn) *ProductService {
	return &ProductService{
		client: pb.NewProductServiceClient(conn),
	}
}

// GetProduct retrieves a product from the product service
func (pc *ProductService) GetProduct(productID uint64) (*pb.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := pc.client.GetProduct(ctx, &pb.GetProductRequest{Id: productID})
	return resp, err
}

// ValidateInventory checks if there is sufficient inventory for the requested quantity
func (pc *ProductService) ValidateInventory(productID, requestedQuantity uint64) (*pb.Product, error) {
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
func (pc *ProductService) UpdateInventory(productID, newQuantity uint64) error {
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
