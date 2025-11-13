package main

import (
	"db-less-store/configs"
	"db-less-store/internal/product"
	"db-less-store/pkg/db"
	"db-less-store/pkg/middleware"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	// Настройка logrus
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	conf := configs.LoadConfig()
	newDb := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepo := product.NewProductRepository(newDb)

	// Handlers
	product.NewProductHandler(router, product.HandlerProductDeps{
		ProductRepository: productRepo,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":7777",
		Handler: stack(router),
	}

	fmt.Println("Server listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server listening error")
		return
	}
}
