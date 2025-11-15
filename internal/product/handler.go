package product

import (
	"db-less-store/pkg/middleware"
	"db-less-store/pkg/req"
	"db-less-store/pkg/res"
	"log"
	"net/http"
)

type HandlerProductDeps struct {
	ProductRepository *ProductRepository
}

type HandlerProduct struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps HandlerProductDeps) {
	handler := &HandlerProduct{
		ProductRepository: deps.ProductRepository,
	}
	router.Handle("POST /products", middleware.IsAuthenticated(handler.CreateProduct()))
	router.Handle("GET /products", middleware.IsAuthenticated(handler.GetProducts()))
	router.Handle("GET /products/{id}", middleware.IsAuthenticated(handler.GetProduct()))
	router.Handle("PATCH /products/{id}", middleware.IsAuthenticated(handler.UpdateProduct()))
	router.Handle("DELETE /products/{id}", middleware.IsAuthenticated(handler.DeleteProduct()))
}

func (h *HandlerProduct) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateProduct")

		body, err := req.HandleBody[CreateProductRequest](&w, r)
		if err != nil {
			return
		}

		product := &Product{
			Name:        body.Name,
			Description: body.Description,
			Price:       body.Price,
			Images:      body.Images,
		}

		err = h.ProductRepository.CreateProduct(product)
		if err != nil {
			log.Println("Error creating product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to create product"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, product, http.StatusCreated)
	}
}

func (h *HandlerProduct) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetProducts")

		products, err := h.ProductRepository.GetAllProducts()
		if err != nil {
			log.Println("Error getting products:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to get products"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, products, http.StatusOK)
	}
}

func (h *HandlerProduct) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetProduct")

		id := r.PathValue("id")
		if id == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Product ID is required"}, http.StatusBadRequest)
			return
		}

		product, err := h.ProductRepository.GetProductById(id)
		if err != nil {
			log.Println("Error getting product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Product not found"}, http.StatusNotFound)
			return
		}

		res.SendJsonResponse(w, product, http.StatusOK)
	}
}

func (h *HandlerProduct) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateProduct")

		id := r.PathValue("id")
		if id == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Product ID is required"}, http.StatusBadRequest)
			return
		}

		// Получаем существующий продукт
		existingProduct, err := h.ProductRepository.GetProductById(id)
		if err != nil {
			log.Println("Error finding product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Product not found"}, http.StatusNotFound)
			return
		}

		body, err := req.HandleBody[UpdateProductRequest](&w, r)
		if err != nil {
			return
		}

		// Обновляем только переданные поля
		if body.Name != "" {
			existingProduct.Name = body.Name
		}
		if body.Description != "" {
			existingProduct.Description = body.Description
		}
		if body.Price != 0 {
			existingProduct.Price = body.Price
		}
		if body.Images != nil {
			existingProduct.Images = body.Images
		}

		err = h.ProductRepository.UpdateProduct(existingProduct)
		if err != nil {
			log.Println("Error updating product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to update product"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, existingProduct, http.StatusOK)
	}
}

func (h *HandlerProduct) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DeleteProduct")

		id := r.PathValue("id")
		if id == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Product ID is required"}, http.StatusBadRequest)
			return
		}

		// Проверяем существование продукта
		_, err := h.ProductRepository.GetProductById(id)
		if err != nil {
			log.Println("Error finding product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Product not found"}, http.StatusNotFound)
			return
		}

		err = h.ProductRepository.DeleteProduct(id)
		if err != nil {
			log.Println("Error deleting product:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to delete product"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, map[string]string{"message": "Product deleted successfully"}, http.StatusOK)
	}
}
