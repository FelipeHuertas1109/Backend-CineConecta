package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
)

func CreateComment(c *models.Comment) error {
	return config.DB.Create(c).Error
}

func GetComments() ([]models.Comment, error) {
	var list []models.Comment
	err := config.DB.Preload("User").Preload("Movie").Find(&list).Error
	return list, err
}

func GetCommentByID(id uint) (models.Comment, error) {
	var c models.Comment
	err := config.DB.Preload("User").Preload("Movie").First(&c, id).Error
	return c, err
}

func UpdateComment(c *models.Comment) error {
	return config.DB.Save(c).Error
}

func DeleteComment(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}
