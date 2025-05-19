package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	"encoding/json"
	"errors"
	"time"
)

// SaveRecommendationsToDataset guarda un conjunto de recomendaciones para un usuario
func SaveRecommendationsToDataset(userID uint, recommendations []models.RecommendationItem) (uint, error) {
	// Verificar que haya al menos una recomendación
	if len(recommendations) == 0 {
		return 0, errors.New("no hay recomendaciones para guardar")
	}

	// Convertir las recomendaciones a JSON
	recommendationsJSON, err := json.Marshal(recommendations)
	if err != nil {
		return 0, err
	}

	// Crear el dataset
	dataset := models.RecommendationDataset{
		UserID:              userID,
		Recommendations:     recommendations,
		RecommendationsJSON: string(recommendationsJSON),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Guardar en la base de datos
	if err := config.DB.Create(&dataset).Error; err != nil {
		return 0, err
	}

	return dataset.ID, nil
}

// GetRecommendationDatasets obtiene todos los datasets de recomendaciones para un usuario
func GetRecommendationDatasets(userID uint) ([]models.RecommendationDataset, error) {
	var datasets []models.RecommendationDataset
	if err := config.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&datasets).Error; err != nil {
		return nil, err
	}

	// Cargar las recomendaciones desde JSON para cada dataset
	for i := range datasets {
		var recommendations []models.RecommendationItem
		if err := json.Unmarshal([]byte(datasets[i].RecommendationsJSON), &recommendations); err != nil {
			return nil, err
		}
		datasets[i].Recommendations = recommendations
	}

	return datasets, nil
}

// GetRecommendationDatasetByID obtiene un dataset específico
func GetRecommendationDatasetByID(id uint) (models.RecommendationDataset, error) {
	var dataset models.RecommendationDataset
	if err := config.DB.First(&dataset, id).Error; err != nil {
		return dataset, err
	}

	// Cargar las recomendaciones desde JSON
	var recommendations []models.RecommendationItem
	if err := json.Unmarshal([]byte(dataset.RecommendationsJSON), &recommendations); err != nil {
		return dataset, err
	}
	dataset.Recommendations = recommendations

	return dataset, nil
}

// DeleteRecommendationDataset elimina un dataset de recomendaciones
func DeleteRecommendationDataset(id uint, userID uint) error {
	// Verificar que el dataset pertenezca al usuario
	var dataset models.RecommendationDataset
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&dataset).Error; err != nil {
		return errors.New("dataset no encontrado o no pertenece al usuario")
	}

	return config.DB.Delete(&dataset).Error
}
