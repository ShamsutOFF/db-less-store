package auth

import (
	"db-less-store/pkg/db"
	"time"
)

type AuthRepository struct {
	Database *db.Db
}

func NewAuthRepository(database *db.Db) *AuthRepository {
	return &AuthRepository{Database: database}
}

func (r *AuthRepository) CreateSession(session *Session) error {
	return r.Database.DB.Create(session).Error
}

func (r *AuthRepository) FindSession(sessionID string) (*Session, error) {
	var session Session
	result := r.Database.DB.Where("session_id = ? AND expires_at > ? AND used = ?",
		sessionID, time.Now(), false).First(&session)

	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func (r *AuthRepository) MarkSessionUsed(sessionID string) error {
	return r.Database.DB.Model(&Session{}).
		Where("session_id = ?", sessionID).
		Update("used", true).Error
}

func (r *AuthRepository) CreateOrUpdateUser(phone string) (*User, error) {
	var user User

	// Ищем существующего пользователя
	result := r.Database.DB.Where("phone = ?", phone).First(&user)
	if result.Error == nil {
		return &user, nil // пользователь уже существует
	}

	// Создаем нового пользователя
	user.Phone = phone
	err := r.Database.DB.Create(&user).Error
	return &user, err
}
