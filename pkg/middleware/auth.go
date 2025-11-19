package middleware

import (
	"context"
	"db-less-store/pkg/jwt"
	"db-less-store/pkg/res"
	"log"
	"net/http"
	"strings"
)

type PhoneKey string

const PhoneContextKey PhoneKey = "userPhone"

var jwtService *jwt.JWT

func SetJWTService(service *jwt.JWT) {
	jwtService = service
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if jwtService == nil {
			res.SendJsonResponse(w, map[string]string{"error": "JWT service not configured"}, http.StatusInternalServerError)
			return
		}
		// Получаем токен из заголовка
		header := r.Header.Get("Authorization")
		if header == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Authorization header required"}, http.StatusUnauthorized)
			return
		}

		// Проверяем формат: Bearer <token>
		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			res.SendJsonResponse(w, map[string]string{"error": "Invalid authorization format. Expected: Bearer <token>"}, http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token == "" {
			res.SendJsonResponse(w, map[string]string{"error": "Token is required"}, http.StatusUnauthorized)
			return
		}

		log.Println("@@@ Verifying token:", token)

		// Верифицируем токен
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			log.Println("@@@ Token verification failed:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
			return
		}

		// Получаем телефон из токена
		phone, ok := claims["phone"].(string)
		if !ok || phone == "" {
			log.Println("@@@ Phone not found in token claims")
			res.SendJsonResponse(w, map[string]string{"error": "Invalid token claims"}, http.StatusUnauthorized)
			return
		}

		log.Println("@@@ Authenticated user phone:", phone)

		// добавил телефон в контекст для использования в хендлерах
		ctx := context.WithValue(r.Context(), PhoneContextKey, phone)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
