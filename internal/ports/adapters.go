package ports

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"yandexLyceumTheme3gRPC/internal/models"
)

type OrdersRepositoryInMemory struct {
	orders map[string]models.Order
	mu     sync.RWMutex
}

func NewOrdersRepositoryInMemory() *OrdersRepositoryInMemory {
	return &OrdersRepositoryInMemory{
		orders: make(map[string]models.Order),
	}
}

func (repo *OrdersRepositoryInMemory) CreateOrder(order models.Order) (models.Order, error) {
	err := ValidateOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	newId := uuid.New().String()
	order.ID = newId

	repo.mu.Lock()
	repo.orders[newId] = order
	repo.mu.Unlock()

	return order, nil
}

func (repo *OrdersRepositoryInMemory) GetOrder(id string) (models.Order, error) {
	repo.mu.RLock()
	order, ok := repo.orders[id]
	repo.mu.RUnlock()
	if !ok {
		return models.Order{}, fmt.Errorf("order not found by id %s", id)
	}
	return order, nil
}

func (repo *OrdersRepositoryInMemory) UpdateOrder(newOrder models.Order) (models.Order, error) {
	err := ValidateOrder(newOrder)
	if err != nil {
		return models.Order{}, err
	}

	repo.mu.Lock()
	_, ok := repo.orders[newOrder.ID]
	repo.mu.Unlock()

	if !ok {
		return models.Order{}, fmt.Errorf("order not found by id %s", newOrder.ID)
	}

	repo.mu.Lock()
	repo.orders[newOrder.ID] = newOrder
	repo.mu.Unlock()

	return newOrder, nil
}

func (repo *OrdersRepositoryInMemory) DeleteOrder(id string) (bool, error) {
	repo.mu.Lock()
	_, ok := repo.orders[id]
	repo.mu.Unlock()

	if !ok {
		return false, fmt.Errorf("order not found by id %s", id)
	}

	repo.mu.Lock()
	delete(repo.orders, id)
	repo.mu.Unlock()

	return true, nil
}

func (repo *OrdersRepositoryInMemory) ListOrders() ([]models.Order, error) {
	repo.mu.RLock()
	result := make([]models.Order, 0, len(repo.orders))
	for _, order := range repo.orders {
		result = append(result, order)
	}
	repo.mu.RUnlock()
	return result, nil
}
