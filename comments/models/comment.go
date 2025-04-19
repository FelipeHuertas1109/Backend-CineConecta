package models

import (
	"time"

	userModels "cine_conecta_backend/auth/models"
	movieModels "cine_conecta_backend/movies/models"
)

type Comment struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `json:"user_id"`
	MovieID   uint              `json:"movie_id"`
	Content   string            `json:"content"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	User      userModels.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie     movieModels.Movie `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
}
