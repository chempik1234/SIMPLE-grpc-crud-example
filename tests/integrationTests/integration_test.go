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
	baseURL := "nginx:80"

	grpcConn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to create grpc connection: %v", err)
	}
	defer grpcConn.Close()

	ordersManager := test.NewOrderServiceClient(grpcConn)

	ctx := context.Background()

	var ok bool

	var expectedOrdersCount int
	var ordersList *test.ListOrdersResponse

	// test list orders before creation
	ok = t.Run("list orders before creation", func(t *testing.T) {
		response, err := ordersManager.ListOrders(ctx, &test.ListOrdersRequest{})
		if err != nil {
			t.Fatalf("Failed to list orders: %v", err)
		}
		expectedOrdersCount = len(response.Orders)
		ordersList = response
	})
	if !ok {
		return
	}

	orderToCreate := &test.Order{
		Id:       "",
		Item:     "new_item",
		Quantity: 1234,
	}

	expectedOrdersCount += 1

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
		if ordersNotEqual(response.Order, orderToCreate) {
			t.Fatalf("Failed to get created order\n"+
				"expected: %v\n"+
				"got %v", orderToCreate, response.Order)
		}
	})
	if !ok {
		return
	}

	// test update an order
	ok = t.Run("update an order", func(t *testing.T) {
		// new order data
		orderToCreate.Item = "new_item2"
		orderToCreate.Quantity = 1235

		response, err := ordersManager.UpdateOrder(ctx, &test.UpdateOrderRequest{
			Id:       orderToCreate.Id,
			Item:     orderToCreate.Item,
			Quantity: orderToCreate.Quantity,
		})
		if err != nil {
			t.Fatalf("Failed to update an order: %v", err)
		}
		if ordersNotEqual(response.Order, orderToCreate) {
			t.Fatalf("Failed to update an order\n"+
				"expected: %v\n"+
				"got %v", orderToCreate, response.Order)
		}
	})

	// test list orders after creation
	ok = t.Run("list orders after creation", func(t *testing.T) {
		response, err := ordersManager.ListOrders(ctx, &test.ListOrdersRequest{})
		if err != nil {
			t.Fatalf("Failed to list orders: %v", err)
		}
		if len(response.Orders) != expectedOrdersCount {
			t.Fatalf("Failed to list order\n"+
				"expected length: %v\n"+
				"got length: %v", expectedOrdersCount, len(response.Orders))
		}
	})
	if !ok {
		return
	}

	expectedOrdersCount -= 1

	// test delete an order
	ok = t.Run("delete an order", func(t *testing.T) {
		_, err := ordersManager.DeleteOrder(ctx, &test.DeleteOrderRequest{
			Id: orderToCreate.Id,
		})
		if err != nil {
			t.Fatalf("Failed to delete an order: %v", err)
		}

		responseList, err := ordersManager.ListOrders(ctx, &test.ListOrdersRequest{})
		if err != nil {
			t.Fatalf("Failed to list orders: %v", err)
		}
		if ordersListsNotEqual(responseList, ordersList) {
			t.Fatalf("Failed to list orders after deletion of an order\n"+
				"length is right, but not the content\n"+
				"expected: %v\n"+
				"got: %v", responseList, ordersList)
		}
	})
}

func ordersListsNotEqual(list1 *test.ListOrdersResponse, list2 *test.ListOrdersResponse) bool {
	if len(list1.Orders) != len(list2.Orders) {
		return true
	}
	for index, item := range list1.Orders {
		if ordersNotEqual(item, list2.Orders[index]) {
			return true
		}
	}
	return false
}

func ordersNotEqual(order1 *test.Order, order2 *test.Order) bool {
	return order1.Item != order2.Item ||
		order1.Quantity != order2.Quantity ||
		order1.Id != order2.Id
}
