package auth

import (
	"fmt"
	"log"
)

type SMSService struct {
	// Здесь могут быть настройки SMS-провайдера
}

func NewSMSService() *SMSService {
	return &SMSService{}
}

func (s *SMSService) SendVerificationCode(phone, code string) error {
	// В реальном приложении здесь будет интеграция с SMS-провайдером
	// Сейчас просто логируем
	log.Printf("SMS to %s: Your verification code is: %s", phone, code)
	fmt.Printf("=== DEMO SMS ===\n")
	fmt.Printf("To: %s\n", phone)
	fmt.Printf("Message: Your verification code: %s\n", code)
	fmt.Printf("===============\n")
	return nil
}
