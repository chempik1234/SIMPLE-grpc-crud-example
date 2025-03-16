package integrationTests

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"yandexLyceumTheme3gRPC/pkg/api/test"
)

type APITestCase struct {
	Name                 string
	Rpc                  string
	Data                 interface{}
	ExpectedResult       interface{}
	ResultFieldToSave    string
	ResultFieldToSaveKey string
}

func TestGRPC(t *testing.T) {
	baseURL := "http://127.0.0.1:50051"

	grpcConn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create grpc connection: %v", err)
	}
	defer grpcConn.Close()

	ordersManager := test.NewOrderServiceClient(grpcConn)

	ctx := context.Background()

	var ok bool

	// test list empty orders
	ok = t.Run("list empty orders", func(t *testing.T) {
		response, err := ordersManager.ListOrders(ctx, &test.ListOrdersRequest{})
		if err != nil {
			t.Fatalf("Failed to list empty orders: %v", err)
		}
		if len(response.Orders) != 0 {
			t.Fatalf("Expected empty orders, got %v", response.Orders)
		}
	})
	if !ok {
		return
	}

	orderToCreate := &test.Order{
		Id:       "",
		Item:     "new_item",
		Quantity: 1234,
	}

	ok = t.Run("create an order", func(t *testing.T) {
		response, err := ordersManager.CreateOrder(ctx, &test.CreateOrderRequest{
			Item:     "new_item",
			Quantity: 1234,
		})
		if err != nil {
			t.Fatalf("Failed to create an order: %v", err)
		}
		orderToCreate.Id = response.Id
	})
	if !ok {
		return
	}

	ok = t.Run("get an order", func(t *testing.T) {
		response, err := ordersManager.GetOrder(ctx, &test.GetOrderRequest{
			Id: orderToCreate.Id,
		})
		if err != nil {
			t.Fatalf("Failed to get an order: %v", err)
		}
		if response.Order != orderToCreate {
			t.Fatalf("Failed to get created order\nexpected: %v\ngot %v", orderToCreate, response.Order)
		}
	})
	if !ok {
		return
	}
}
