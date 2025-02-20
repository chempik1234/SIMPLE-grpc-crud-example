package ports

import "yandexLyceumTheme3gRPC/internal/models"

type OrdersRepository interface {
	CreateOrder(newOrder models.Order) (models.Order, error)
	GetOrder(id string) (models.Order, error)
	UpdateOrder(newOrder models.Order) (models.Order, error)
	DeleteOrder(id string) (bool, error)
	ListOrders() ([]models.Order, error)
}
