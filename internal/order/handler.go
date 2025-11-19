package order

import (
	"db-less-store/internal/auth"
	"db-less-store/internal/product"
	"db-less-store/pkg/middleware"
	"db-less-store/pkg/req"
	"db-less-store/pkg/res"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type OrderHandlerDeps struct {
	OrderRepository   *OrderRepository
	AuthRepository    *auth.AuthRepository
	ProductRepository *product.ProductRepository
}

type OrderHandler struct {
	OrderRepository   *OrderRepository
	AuthRepository    *auth.AuthRepository
	ProductRepository *product.ProductRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository:   deps.OrderRepository,
		AuthRepository:    deps.AuthRepository,
		ProductRepository: deps.ProductRepository,
	}

	router.Handle("POST /order", middleware.IsAuthenticated(handler.CreateOrder()))
	router.Handle("GET /order/{id}", middleware.IsAuthenticated(handler.GetOrder()))
	router.Handle("GET /my-orders", middleware.IsAuthenticated(handler.GetMyOrders()))
}

func (h *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateOrder")

		// Получаем phone из контекста (из middleware)
		phone, ok := r.Context().Value(middleware.PhoneContextKey).(string)
		if !ok || phone == "" {
			res.SendJsonResponse(w, map[string]string{"error": "User not authenticated"}, http.StatusUnauthorized)
			return
		}

		// Находим пользователя по phone
		user, err := h.AuthRepository.FindUserByPhone(phone)
		if err != nil {
			log.Println("Error finding user:", err)
			res.SendJsonResponse(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
			return
		}

		body, err := req.HandleBody[CreateOrderRequest](&w, r)
		if err != nil {
			return
		}

		// Создаем заказ
		order := &Order{
			UserID: user.ID, // используем реальный userID
			Status: "pending",
			Items:  []OrderItem{},
		}

		// Добавляем товары в заказ
		totalPrice := float32(0)
		for _, itemReq := range body.Items {
			// Получаем продукт из базы
			product, err := h.ProductRepository.GetProductById(strconv.FormatUint(uint64(itemReq.ProductID), 10))
			if err != nil {
				log.Println("Product not found:", itemReq.ProductID)
				res.SendJsonResponse(w, map[string]string{"error": "Product not found: " + strconv.FormatUint(uint64(itemReq.ProductID), 10)}, http.StatusBadRequest)
				return
			}

			// Создаем позицию заказа
			orderItem := OrderItem{
				ProductID: itemReq.ProductID,
				Quantity:  itemReq.Quantity,
				Price:     product.Price, // сохраняем цену на момент заказа
			}

			order.Items = append(order.Items, orderItem)
			totalPrice += product.Price * float32(itemReq.Quantity)
		}

		order.TotalPrice = totalPrice

		// Сохраняем заказ в базу
		err = h.OrderRepository.CreateOrder(order)
		if err != nil {
			log.Println("Error creating order:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to create order"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, CreateOrderResponse{OrderID: order.ID}, http.StatusCreated)
	}
}

func (h *OrderHandler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetMyOrders")

		// Получаем phone из контекста (из middleware)
		phone, ok := r.Context().Value(middleware.PhoneContextKey).(string)
		if !ok || phone == "" {
			res.SendJsonResponse(w, map[string]string{"error": "User not authenticated"}, http.StatusUnauthorized)
			return
		}

		// Находим пользователя по phone
		user, err := h.AuthRepository.FindUserByPhone(phone)
		if err != nil {
			log.Println("Error finding user:", err)
			res.SendJsonResponse(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
			return
		}

		orders, err := h.OrderRepository.GetOrdersByUserID(user.ID)
		if err != nil {
			log.Println("Error getting orders:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to get orders"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, orders, http.StatusOK)
	}
}

func (h *OrderHandler) GetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetOrder")

		// Получаем phone из контекста (из middleware)
		phone, ok := r.Context().Value(middleware.PhoneContextKey).(string)
		if !ok || phone == "" {
			res.SendJsonResponse(w, map[string]string{"error": "User not authenticated"}, http.StatusUnauthorized)
			return
		}

		// Находим пользователя по phone
		currentUser, err := h.AuthRepository.FindUserByPhone(phone)
		if err != nil {
			log.Println("Error finding user:", err)
			res.SendJsonResponse(w, map[string]string{"error": "User not found"}, http.StatusUnauthorized)
			return
		}

		id := r.PathValue("id")
		if id == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Order ID is required"}, http.StatusBadRequest)
			return
		}

		orderID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			res.SendJsonResponse(w, map[string]string{"error": "Invalid order ID"}, http.StatusBadRequest)
			return
		}

		order, err := h.OrderRepository.GetOrderByID(uint(orderID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				res.SendJsonResponse(w, map[string]string{"error": "Order not found"}, http.StatusNotFound)
				return
			}
			log.Println("Error getting order:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to get order"}, http.StatusInternalServerError)
			return
		}

		// Проверяем что заказ принадлежит текущему пользователю
		if order.UserID != currentUser.ID {
			res.SendJsonResponse(w, map[string]string{"error": "Access denied"}, http.StatusForbidden)
			return
		}
		order.User.Phone = currentUser.Phone

		res.SendJsonResponse(w, order, http.StatusOK)
	}
}
