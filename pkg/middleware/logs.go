package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		// Выполняем запрос
		next.ServeHTTP(wrapper, r)

		// Создаем JSON лог с помощью logrus
		logrus.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     wrapper.StatusCode,
			"duration":   time.Since(start).String(),
			"user_agent": r.UserAgent(),
			"ip":         r.RemoteAddr,
		}).Info("HTTP request")
	})
}
