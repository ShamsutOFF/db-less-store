package order

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type CreateOrderResponse struct {
	OrderID uint `json:"order_id"`
}
