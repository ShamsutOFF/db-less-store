package order

import (
	"db-less-store/pkg/db"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{Database: database}
}

func (r *OrderRepository) CreateOrder(order *Order) error {
	return r.Database.DB.Create(order).Error
}

func (r *OrderRepository) GetOrderByID(orderID uint) (*Order, error) {
	var order Order
	result := r.Database.DB.
		Preload("Items").
		Preload("Items.Product").
		First(&order, orderID)

	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(userID uint) ([]Order, error) {
	var orders []Order
	result := r.Database.DB.
		Preload("Items").
		Preload("Items.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (r *OrderRepository) OrderExists(orderID uint) (bool, error) {
	var count int64
	result := r.Database.DB.Model(&Order{}).Where("id = ?", orderID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
