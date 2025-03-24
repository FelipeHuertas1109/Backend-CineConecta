package services

import (
	"cine_conecta_backend/auth/models"
	"cine_conecta_backend/config"
)

func SaveUser(user *models.User) error {
	return config.DB.Create(user).Error
}
