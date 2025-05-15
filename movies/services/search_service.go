package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"strings"
)

// SearchParams contiene los parámetros de búsqueda
type SearchParams struct {
	Title  string  `json:"title"`  // Búsqueda por título
	Genre  string  `json:"genre"`  // Filtro por género
	Rating float64 `json:"rating"` // Puntuación mínima
}

// SearchMovies busca películas según criterios de búsqueda
func SearchMovies(params SearchParams) ([]models.Movie, error) {
	var movies []models.Movie
	query := config.DB

	// Filtro por título
	if params.Title != "" {
		query = query.Where("title ILIKE ?", "%"+params.Title+"%")
	}

	// Filtro por género
	if params.Genre != "" {
		query = query.Where("genre ILIKE ?", "%"+params.Genre+"%")
	}

	// Filtro por puntuación
	if params.Rating > 0 {
		query = query.Where("rating >= ?", params.Rating)
	}

	// Ejecutar la consulta
	if err := query.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

// GetAllGenres obtiene todos los géneros disponibles
func GetAllGenres() ([]string, error) {
	var movies []models.Movie
	if err := config.DB.Select("genre").Find(&movies).Error; err != nil {
		return nil, err
	}

	// Usar un mapa para evitar géneros duplicados
	genresMap := make(map[string]bool)
	for _, movie := range movies {
		if movie.Genre != "" {
			// Algunos géneros pueden ser compuestos (ej: "Acción, Aventura")
			genreParts := strings.Split(movie.Genre, ",")
			for _, part := range genreParts {
				genreTrimmed := strings.TrimSpace(part)
				if genreTrimmed != "" {
					genresMap[genreTrimmed] = true
				}
			}
		}
	}

	// Convertir el mapa a slice
	var genres []string
	for genre := range genresMap {
		genres = append(genres, genre)
	}

	return genres, nil
}
