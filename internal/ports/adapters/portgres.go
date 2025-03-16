package adapters

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"yandexLyceumTheme3gRPC/internal/models"
	"yandexLyceumTheme3gRPC/internal/ports"
)

type OrdersRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewOrdersRepositoryPostgres(pool *pgxpool.Pool) *OrdersRepositoryPostgres {
	return &OrdersRepositoryPostgres{
		pool: pool,
	}
}

func (repo *OrdersRepositoryPostgres) CreateOrder(order models.Order) (models.Order, error) {
	// validate
	err := ports.ValidateOrder(order)
	if err != nil {
		return models.Order{}, err
	}

	// build insert query
	sql, args, err := squirrel.Insert("schema_name.orders").Columns("item", "quantity").
		Values(order.Item, order.Quantity).
		Suffix("RETURNING ID").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.Order{}, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	// perform insert query
	_ = repo.pool.QueryRow(context.Background(), sql, args...).Scan(&order.ID)

	return order, nil
}

func (repo *OrdersRepositoryPostgres) GetOrder(id string) (models.Order, error) {
	// build select query
	sql, args, err := squirrel.Select("id", "item", "quantity").
		From("schema_name.orders").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.Order{}, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	var order models.Order

	// perform select query
	_ = repo.pool.QueryRow(context.Background(), sql, args...).Scan(&order.ID, &order.Item, &order.Quantity)

	if order.ID == "" {
		return models.Order{}, fmt.Errorf("order not found by id %s", id)
	}

	return order, nil
}

func (repo *OrdersRepositoryPostgres) UpdateOrder(newOrder models.Order) (models.Order, error) {
	// validate
	err := ports.ValidateOrder(newOrder)
	if err != nil {
		return models.Order{}, err
	}

	// build select query
	sql, args, err := squirrel.Select("id", "item", "quantity").
		From("schema_name.orders").
		Where(squirrel.Eq{"id": newOrder.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return models.Order{}, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	var order models.Order

	// perform select query
	_ = repo.pool.QueryRow(context.Background(), sql, args...).Scan(&order.ID, &order.Item, &order.Quantity)

	if order.ID == "" {
		return models.Order{}, fmt.Errorf("order not found by id %s", newOrder.ID)
	}

	// build update query
	sql, args, err = squirrel.Update("schema_name.orders").
		Set("item", newOrder.Item).
		Set("quantity", newOrder.Quantity).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	_ = repo.pool.QueryRow(context.Background(), sql, args...).Scan(&order.ID, &order.Item, &order.Quantity)

	return newOrder, nil
}

func (repo *OrdersRepositoryPostgres) DeleteOrder(id string) (bool, error) {
	// build select query
	sql, args, err := squirrel.Select("id").
		From("schema_name.orders").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	var foundId string

	// perform select query
	_ = repo.pool.QueryRow(context.Background(), sql, args...).Scan(&foundId)

	if foundId == "" {
		return false, fmt.Errorf("order not found by id %s", id)
	}

	// build delete query
	sql, args, err = squirrel.Delete("schema_name.orders").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	// perform delete query
	_ = repo.pool.QueryRow(context.Background(), sql, args...)

	return true, nil
}

func (repo *OrdersRepositoryPostgres) ListOrders() ([]models.Order, error) {
	// build select query
	sql, args, err := squirrel.Select("id", "item", "quantity").
		From("schema_name.orders").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Couldn't build and SQL query: %v", err)
	}

	// perform select query
	rows, err := repo.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Couldn't perform select query: %v", err)
	}
	defer rows.Close()

	var result []models.Order

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Item, &order.Quantity)
		if err != nil {
			return nil, fmt.Errorf("Couldn't scan row: %v", err)
		}
		result = append(result, order)
	}

	return result, nil
}
