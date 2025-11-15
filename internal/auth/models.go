package auth

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Phone     string    `json:"phone" gorm:"not null"`
	Code      string    `json:"code" gorm:"not null"`
	SessionID string    `json:"session_id" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
}

type User struct {
	gorm.Model
	Phone string `json:"phone" gorm:"uniqueIndex;not null"`
	// Можно добавить другие поля позже (имя, email и т.д.)
}
