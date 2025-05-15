package models

import "time"

type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
	PosterURL   string    `json:"poster_url"` // 🆕

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
