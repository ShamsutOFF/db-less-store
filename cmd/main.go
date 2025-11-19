package main

import (
	"db-less-store/configs"
	"db-less-store/internal/auth"
	"db-less-store/internal/order"
	"db-less-store/internal/product"
	"db-less-store/pkg/db"
	"db-less-store/pkg/jwt"
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
	orderRepo := order.NewOrderRepository(newDb)

	// Auth dependencies
	authRepo := auth.NewAuthRepository(newDb)
	smsService := auth.NewSMSService()
	codeGenerator := auth.NewCodeGenerator()
	jwtService := jwt.NewJWT(conf.Auth.Secret)
	middleware.SetJWTService(jwtService)

	// Handlers
	product.NewProductHandler(router, product.HandlerProductDeps{
		ProductRepository: productRepo,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthRepository: authRepo,
		SMSService:     smsService,
		CodeGenerator:  codeGenerator,
		JWTService:     jwtService,
	})
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository:   orderRepo,
		AuthRepository:    authRepo,
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
