package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"errors"
)

// CreateLike crea un nuevo "me gusta" para una película
func CreateLike(userID, movieID uint) error {
	// Verificar si la película existe
	var movie models.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return errors.New("película no encontrada")
	}

	// Verificar si ya existe un like de este usuario para esta película
	var existingLike models.Like
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&existingLike)
	if result.RowsAffected > 0 {
		return errors.New("ya has dado me gusta a esta película")
	}

	// Crear el nuevo like
	like := models.Like{
		UserID:  userID,
		MovieID: movieID,
	}

	if err := config.DB.Create(&like).Error; err != nil {
		return errors.New("error al dar me gusta a la película")
	}

	return nil
}

// DeleteLike elimina un "me gusta" de una película
func DeleteLike(userID, movieID uint) error {
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&models.Like{})
	if result.RowsAffected == 0 {
		return errors.New("no has dado me gusta a esta película")
	}
	return nil
}

// GetLikeByUserAndMovie verifica si un usuario ha dado me gusta a una película
func GetLikeByUserAndMovie(userID, movieID uint) (bool, error) {
	var like models.Like
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&like)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// GetLikesByUser obtiene todas las películas a las que un usuario ha dado me gusta
func GetLikesByUser(userID uint) ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Joins("JOIN movie_likes ON movies.id = movie_likes.movie_id").
		Where("movie_likes.user_id = ?", userID).
		Find(&movies).Error

	if err != nil {
		return nil, errors.New("error al obtener las películas con me gusta")
	}

	return movies, nil
}

// GetLikesByMovie obtiene todos los usuarios que han dado me gusta a una película
func GetLikesByMovie(movieID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Like{}).Where("movie_id = ?", movieID).Count(&count).Error
	if err != nil {
		return 0, errors.New("error al obtener el conteo de me gusta")
	}
	return count, nil
}

// GetAllLikes obtiene todos los "me gusta" con información de películas y usuarios
func GetAllLikes() ([]models.Like, error) {
	var likes []models.Like
	if err := config.DB.Find(&likes).Error; err != nil {
		return nil, errors.New("error al obtener los me gusta")
	}
	return likes, nil
}
