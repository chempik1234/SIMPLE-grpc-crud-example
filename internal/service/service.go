package service

import (
	"context"
	"go.uber.org/zap"
	"yandexLyceumTheme3gRPC/internal/models"
	"yandexLyceumTheme3gRPC/internal/ports"
	"yandexLyceumTheme3gRPC/pkg/api/test"
	"yandexLyceumTheme3gRPC/pkg/logger"
)

const (
	orderId     = "order_id"
	ordersCount = "orders_count"
)

type Service struct {
	test.OrderServiceServer
	ordersRepoStorage ports.OrdersRepository
}

func New(ordersRepoStorage ports.OrdersRepository) *Service {
	return &Service{
		ordersRepoStorage: ordersRepoStorage,
	}
}

func (s *Service) CreateOrder(ctx context.Context, req *test.CreateOrderRequest) (*test.CreateOrderResponse, error) {
	result, err := s.ordersRepoStorage.CreateOrder(models.Order{
		Item:     req.Item,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "created an order", zap.String(orderId, result.ID))
	return &test.CreateOrderResponse{Id: result.ID}, nil
}

func (s *Service) GetOrder(ctx context.Context, req *test.GetOrderRequest) (*test.GetOrderResponse, error) {
	result, err := s.ordersRepoStorage.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "got an order", zap.String(orderId, result.ID))
	return &test.GetOrderResponse{Order: &test.Order{
		Id:       result.ID,
		Item:     result.Item,
		Quantity: result.Quantity,
	}}, nil
}

func (s *Service) UpdateOrder(ctx context.Context, req *test.UpdateOrderRequest) (*test.UpdateOrderResponse, error) {
	result, err := s.ordersRepoStorage.UpdateOrder(models.Order{
		ID:       req.Id,
		Item:     req.Item,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "updated an order", zap.String(orderId, result.ID))
	return &test.UpdateOrderResponse{Order: &test.Order{
		Id:       result.ID,
		Item:     result.Item,
		Quantity: result.Quantity,
	}}, nil
}

func (s *Service) DeleteOrder(ctx context.Context, req *test.DeleteOrderRequest) (*test.DeleteOrderResponse, error) {
	result, err := s.ordersRepoStorage.DeleteOrder(req.Id)
	if err != nil {
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "deleted an order", zap.String(orderId, req.Id))
	return &test.DeleteOrderResponse{Success: result}, nil
}

func (s *Service) ListOrders(ctx context.Context, req *test.ListOrdersRequest) (*test.ListOrdersResponse, error) {
	result, err := s.ordersRepoStorage.ListOrders()
	if err != nil {
		return nil, err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "someone listed orders", zap.Int(ordersCount, len(result)))

	ordersArray := make([]*test.Order, len(result))
	for i, order := range result {
		ordersArray[i] = &test.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		}
	}
	return &test.ListOrdersResponse{Orders: ordersArray}, nil
}
