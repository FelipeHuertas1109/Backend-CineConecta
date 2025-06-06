package models

import "time"

// Genre representa un género cinematográfico
type Genre struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`

	// Relación muchos a muchos con películas
	Movies []Movie `gorm:"many2many:movie_genres;" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
