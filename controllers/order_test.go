package controllers

import (
	"order/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockOrderService is a mock of OrderService
type MockOrderService struct {
	mock.Mock
}

// Implement all service methods for the mock
func (m *MockOrderService) PlaceOrder(sessionId string, userId uint64) (*models.Order, error) {
	args := m.Called(sessionId, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrder(id uint64) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrdersByUser(userId uint64) ([]*models.Order, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrdersByMerchant(merchantId uint64) ([]*models.Order, error) {
	args := m.Called(merchantId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(id uint64, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderService) CancelOrder(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrderService) DeleteOrder(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

// Setup test router and controller
func setupTestRouter(mockService *MockOrderService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	controller := &OrderController{
		orderService: mockService,
	}

	controller.SetupRoutes(router)
	return router
}
