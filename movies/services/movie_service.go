package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
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
