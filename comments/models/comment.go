package models

import (
	"time"

	userModels "cine_conecta_backend/auth/models"
	movieModels "cine_conecta_backend/movies/models"
)

// SentimentType representa el tipo de sentimiento del comentario
type SentimentType string

const (
	SentimentPositive SentimentType = "positive"
	SentimentNeutral  SentimentType = "neutral"
	SentimentNegative SentimentType = "negative"
)

type Comment struct {
	ID             uint              `gorm:"primaryKey" json:"id"`
	UserID         uint              `gorm:"uniqueIndex:idx_user_movie" json:"user_id"`
	MovieID        uint              `gorm:"uniqueIndex:idx_user_movie" json:"movie_id"`
	Content        string            `json:"content"`
	Sentiment      SentimentType     `json:"sentiment"`
	SentimentScore float64           `json:"sentiment_score"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	User           userModels.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie          movieModels.Movie `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}
