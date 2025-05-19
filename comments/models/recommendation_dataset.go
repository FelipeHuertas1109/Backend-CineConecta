package models

import (
	"time"
)

// RecommendationItem representa un elemento individual de recomendación
type RecommendationItem struct {
	MovieID         uint    `json:"movie_id"`
	Title           string  `json:"title"`
	PredictedRating float64 `json:"predicted_rating"`
	RatingText      string  `json:"rating_text"`
	Reason          string  `json:"reason"`
}

// RecommendationDataset representa un conjunto de recomendaciones guardadas
type RecommendationDataset struct {
	ID                  uint                 `gorm:"primaryKey" json:"id"`
	UserID              uint                 `json:"user_id"`
	Recommendations     []RecommendationItem `gorm:"-" json:"recommendations"` // No se almacena directamente en la base de datos
	RecommendationsJSON string               `gorm:"type:jsonb" json:"-"`      // Se almacena como JSON en la base de datos
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}
