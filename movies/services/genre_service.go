package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"fmt"
	"sort"
	"strings"
)

// GetUniqueGenres obtiene todos los géneros únicos utilizados en las películas
func GetUniqueGenres() ([]string, error) {
	var movies []models.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	// Usar un mapa para evitar duplicados
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

	// Ordenar géneros alfabéticamente
	sort.Strings(genres)

	return genres, nil
}

// GetAllGenres obtiene todos los géneros únicos (compatibilidad con código existente)
func GetAllGenres() ([]string, error) {
	return GetUniqueGenres()
}

// GetGenreForMovie obtiene el género de una película
func GetGenreForMovie(movieID uint) (string, error) {
	var movie models.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return "", err
	}
	return movie.Genre, nil
}

// AddGenreToMovie asigna un género a una película
func AddGenreToMovie(movieID uint, genreName string) error {
	var movie models.Movie

	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	// Asignar el género a la película
	movie.Genre = genreName

	return config.DB.Save(&movie).Error
}

// RemoveGenreFromMovie elimina el género de una película
func RemoveGenreFromMovie(movieID uint, genreName string) error {
	var movie models.Movie

	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	// Verificar que el género actual sea el que queremos eliminar
	if movie.Genre == genreName {
		movie.Genre = ""
		return config.DB.Save(&movie).Error
	}

	return nil
}

// GetMoviesByGenre obtiene todas las películas de un género específico
func GetMoviesByGenre(genreName string) ([]models.Movie, error) {
	var movies []models.Movie

	fmt.Printf("[DEBUG-GENRE] Buscando películas con género: %s\n", genreName)

	// Usar LOWER para hacer la búsqueda case-insensitive
	if err := config.DB.Where("LOWER(genre) LIKE LOWER(?)", "%"+genreName+"%").Find(&movies).Error; err != nil {
		fmt.Printf("[DEBUG-GENRE] Error al buscar películas por género: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG-GENRE] Encontradas %d películas con género %s\n", len(movies), genreName)

	return movies, nil
}

// UpdateMovieGenre actualiza el género de una película
func UpdateMovieGenre(movieID uint, genreName string) error {
	var movie models.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	// Actualizar el género
	movie.Genre = genreName
	return config.DB.Save(&movie).Error
}

// GetGenreStats obtiene estadísticas de un género
func GetGenreStats(genreName string) (*GenreInfo, error) {
	fmt.Printf("[DEBUG-GENRE] Obteniendo estadísticas para el género: %s\n", genreName)

	// Obtener películas con este género usando LOWER para búsqueda case-insensitive
	var movies []models.Movie
	if err := config.DB.Where("LOWER(genre) LIKE LOWER(?)", "%"+genreName+"%").Find(&movies).Error; err != nil {
		fmt.Printf("[DEBUG-GENRE] Error al obtener películas para estadísticas: %v\n", err)
		return nil, err
	}

	info := &GenreInfo{
		Name:        genreName,
		Count:       len(movies),
		TotalRating: 0,
		AvgRating:   0,
	}

	if info.Count > 0 {
		for _, movie := range movies {
			info.TotalRating += float64(movie.Rating)
		}
		info.AvgRating = info.TotalRating / float64(info.Count)
	}

	fmt.Printf("[DEBUG-GENRE] Estadísticas para %s: %d películas, rating promedio %.2f\n",
		genreName, info.Count, info.AvgRating)

	return info, nil
}
