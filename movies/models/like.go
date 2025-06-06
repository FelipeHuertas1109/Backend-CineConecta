package models

import (
	"time"

	"gorm.io/gorm"
)

// Like representa un "me gusta" de un usuario a una película
type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	MovieID   uint           `json:"movie_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName especifica el nombre de la tabla en la base de datos
func (Like) TableName() string {
	return "movie_likes"
}
