package ports

import (
	"fmt"
	"yandexLyceumTheme3gRPC/internal/models"
)

func ValidateOrder(order models.Order) error {
	if order.Item == "" {
		return fmt.Errorf("item cannot be blank")
	}
	if order.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	return nil
}
