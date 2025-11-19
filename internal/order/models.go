package order

import (
	"db-less-store/internal/auth"
	"db-less-store/internal/product"
	"gorm.io/gorm"
)

// Order - заказ пользователя
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" gorm:"not null"`
	User       auth.User   `json:"user" gorm:"foreignKey:UserID"`
	Status     string      `json:"status" gorm:"default:'pending'"` // pending, paid, shipped, completed, cancelled
	TotalPrice float32     `json:"total_price"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// OrderItem - позиция в заказе (связь многие-ко-многим с продуктами)
type OrderItem struct {
	gorm.Model
	OrderID   uint            `json:"order_id" gorm:"not null"`
	ProductID uint            `json:"product_id" gorm:"not null"`
	Product   product.Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int             `json:"quantity" gorm:"not null;default:1"`
	Price     float32         `json:"price" gorm:"not null"` // цена на момент заказа (фиксированная)
}
