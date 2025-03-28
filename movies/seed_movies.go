package main

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"fmt"
	"time"
)

func main() {
	config.ConnectDB()

	movies := []models.Movie{
		{
			Title:       "Inception",
			Description: "Un ladrón que roba secretos a través de los sueños debe realizar una misión inversa: implantar una idea.",
			Genre:       "Ciencia Ficción",
			Director:    "Christopher Nolan",
			ReleaseDate: time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC),
			Rating:      4.8,
		},
		{
			Title:       "The Godfather",
			Description: "La historia de una familia mafiosa italoamericana que lucha por mantener su imperio.",
			Genre:       "Crimen",
			Director:    "Francis Ford Coppola",
			ReleaseDate: time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
			Rating:      4.9,
		},
		{
			Title:       "Parasite",
			Description: "Una familia pobre se infiltra en una casa rica con consecuencias inesperadas.",
			Genre:       "Drama",
			Director:    "Bong Joon-ho",
			ReleaseDate: time.Date(2019, 5, 30, 0, 0, 0, 0, time.UTC),
			Rating:      4.6,
		},
		{
			Title:       "The Dark Knight",
			Description: "Batman se enfrenta al Joker, un enemigo que desata el caos en Gotham.",
			Genre:       "Acción",
			Director:    "Christopher Nolan",
			ReleaseDate: time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC),
			Rating:      4.9,
		},
		{
			Title:       "Pulp Fiction",
			Description: "Historias entrelazadas de crimen, redención y violencia en Los Ángeles.",
			Genre:       "Crimen",
			Director:    "Quentin Tarantino",
			ReleaseDate: time.Date(1994, 10, 14, 0, 0, 0, 0, time.UTC),
			Rating:      4.7,
		},
		{
			Title:       "Spirited Away",
			Description: "Una niña entra en un mundo espiritual y debe rescatar a sus padres.",
			Genre:       "Animación",
			Director:    "Hayao Miyazaki",
			ReleaseDate: time.Date(2001, 7, 20, 0, 0, 0, 0, time.UTC),
			Rating:      4.8,
		},
	}

	for _, movie := range movies {
		var existing models.Movie
		// Verificar si ya existe una película con el mismo título
		result := config.DB.Where("title = ?", movie.Title).First(&existing)
		if result.RowsAffected > 0 {
			fmt.Printf("⚠️  La película '%s' ya existe. Saltando...\n", movie.Title)
			continue
		}

		// Crear si no existe
		if err := config.DB.Create(&movie).Error; err != nil {
			fmt.Printf("❌ Error al insertar '%s': %v\n", movie.Title, err)
		} else {
			fmt.Printf("✅ Película '%s' insertada correctamente.\n", movie.Title)
		}
	}
}
