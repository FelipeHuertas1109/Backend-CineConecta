package models

import (
	"strings"
	"time"
)

type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
	PosterURL   string    `json:"poster_url"`

	// Campo para el género como texto (para facilidad de uso)
	Genre string `json:"genre" gorm:"column:genre"`

	// Relación muchos a muchos con géneros
	Genres []Genre `gorm:"many2many:movie_genres;" json:"genres"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ParseGenresString convierte una cadena de géneros en una lista de nombres de géneros
func ParseGenresString(genreStr string) []string {
	if genreStr == "" {
		return []string{}
	}

	// Dividir por comas
	parts := strings.Split(genreStr, ",")

	// Limpiar espacios y filtrar elementos vacíos
	var genres []string
	for _, part := range parts {
		genre := strings.TrimSpace(part)
		if genre != "" {
			genres = append(genres, genre)
		}
	}

	return genres
}
