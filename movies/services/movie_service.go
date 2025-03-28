package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"errors"
)

// CreateMovie guarda una nueva película en la base de datos.
func CreateMovie(movie *models.Movie) error {
	return config.DB.Create(movie).Error
}

// GetMovies obtiene todas las películas.
func GetMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Find(&movies).Error
	return movies, err
}

// GetMovieByID obtiene una película por su ID.
func GetMovieByID(id uint) (models.Movie, error) {
	var movie models.Movie
	err := config.DB.First(&movie, id).Error
	return movie, err
}

// UpdateMovie actualiza una película existente.
func UpdateMovie(movie *models.Movie) error {
	return config.DB.Save(movie).Error
}

// DeleteMovie elimina una película por su ID.
func DeleteMovie(id uint) error {
	return config.DB.Delete(&models.Movie{}, id).Error
}

var allowedSortFields = map[string]bool{
	"title":  true,
	"genre":  true,
	"rating": true,
}

func GetMoviesSorted(sortBy string, order string) ([]models.Movie, error) {
	var movies []models.Movie

	// Validar columna
	if !allowedSortFields[sortBy] {
		return nil, errors.New("columna inválida para ordenar")
	}

	// Validar orden
	if order != "asc" && order != "desc" {
		return nil, errors.New("orden inválido")
	}

	sortQuery := sortBy + " " + order

	result := config.DB.Order(sortQuery).Find(&movies)

	if result.Error != nil {
		return nil, result.Error
	}

	return movies, nil
}
