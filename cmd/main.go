package main

import (
	"db-less-store/configs"
	"db-less-store/internal/product"
	"db-less-store/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	newDb := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepo := product.NewProductRepository(newDb)

	// Handlers
	product.NewProductHandler(router, product.HandlerProductDeps{
		ProductRepository: productRepo,
	})

	server := http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	fmt.Println("Server listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server listening error")
		return
	}
}
