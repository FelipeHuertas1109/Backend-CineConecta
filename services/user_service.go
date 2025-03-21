package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/models"
)

func SaveUser(user *models.User) error {
	return config.DB.Create(user).Error
}
