package auth

import (
	jwt "db-less-store/pkg/jwt"
	"db-less-store/pkg/req"
	"db-less-store/pkg/res"

	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandlerDeps struct {
	AuthRepository *AuthRepository
	SMSService     *SMSService
	CodeGenerator  *CodeGenerator
	JWTService     *jwt.JWT
}

type AuthHandler struct {
	AuthRepository *AuthRepository
	SMSService     *SMSService
	CodeGenerator  *CodeGenerator
	JWTService     *jwt.JWT
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		AuthRepository: deps.AuthRepository,
		SMSService:     deps.SMSService,
		CodeGenerator:  deps.CodeGenerator,
		JWTService:     deps.JWTService,
	}

	router.HandleFunc("POST /auth/init", handler.InitAuth())
	router.HandleFunc("POST /auth/verify", handler.VerifyCode())
}

func (h *AuthHandler) InitAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("InitAuth")

		body, err := req.HandleBody[InitAuthRequest](&w, r)
		if err != nil {
			return
		}

		// Генерируем код и сессию
		code := h.CodeGenerator.GenerateCode()
		sessionID := uuid.New().String()

		// Создаем сессию
		session := &Session{
			Phone:     body.Phone,
			Code:      code,
			SessionID: sessionID,
			ExpiresAt: time.Now().Add(10 * time.Minute), // 10 минут на ввод кода
			Used:      false,
		}

		err = h.AuthRepository.CreateSession(session)
		if err != nil {
			log.Println("Error creating session:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to create session"}, http.StatusInternalServerError)
			return
		}

		// Отправляем SMS (в демо-режиме)
		err = h.SMSService.SendVerificationCode(body.Phone, code)
		if err != nil {
			log.Println("Error sending SMS:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to send SMS"}, http.StatusInternalServerError)
			return
		}

		// Возвращаем sessionId клиенту
		res.SendJsonResponse(w, InitAuthResponse{SessionID: sessionID}, http.StatusOK)
	}
}

func (h *AuthHandler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("VerifyCode")

		body, err := req.HandleBody[VerifyCodeRequest](&w, r)
		if err != nil {
			return
		}

		// Ищем сессию
		session, err := h.AuthRepository.FindSession(body.SessionID)
		if err != nil {
			log.Println("Session not found or expired:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Invalid or expired session"}, http.StatusBadRequest)
			return
		}

		// Проверяем код
		if session.Code != body.Code {
			res.SendJsonResponse(w, map[string]string{"error": "Invalid code"}, http.StatusBadRequest)
			return
		}

		// Помечаем сессию использованной
		err = h.AuthRepository.MarkSessionUsed(body.SessionID)
		if err != nil {
			log.Println("Error marking session used:", err)
		}

		// Создаем/обновляем пользователя
		user, err := h.AuthRepository.CreateOrUpdateUser(session.Phone)
		if err != nil {
			log.Println("Error creating user:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to create user"}, http.StatusInternalServerError)
			return
		}

		// Генерируем JWT токен
		token, err := h.JWTService.GenerateToken(user.Phone)
		if err != nil {
			log.Println("Error generating token:", err)
			res.SendJsonResponse(w, map[string]string{"error": "Failed to generate token"}, http.StatusInternalServerError)
			return
		}

		res.SendJsonResponse(w, VerifyCodeResponse{Token: token}, http.StatusOK)
	}
}
