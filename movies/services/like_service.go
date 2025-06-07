package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"errors"
	"time"
)

// CreateLike crea un nuevo "me gusta" para una película
func CreateLike(userID, movieID uint) error {
	// Verificar si la película existe
	var movie models.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return errors.New("película no encontrada")
	}

	// Verificar si ya existe un like activo de este usuario para esta película
	var existingLike models.Like
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&existingLike)
	if result.RowsAffected > 0 {
		return errors.New("ya has dado me gusta a esta película")
	}

	// Comprobar si existe un registro eliminado con soft delete
	var deletedLike models.Like
	resultDeleted := config.DB.Unscoped().Where("user_id = ? AND movie_id = ? AND deleted_at IS NOT NULL", userID, movieID).First(&deletedLike)

	// Si existe un registro eliminado, lo restauramos actualizando sus timestamps
	if resultDeleted.RowsAffected > 0 {
		config.DB.Unscoped().Model(&deletedLike).Updates(map[string]interface{}{
			"deleted_at": nil,
			"updated_at": time.Now(),
		})
		return nil
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
	// Primero verificamos si existe un "me gusta" activo
	var like models.Like
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&like)
	if result.RowsAffected == 0 {
		return errors.New("no has dado me gusta a esta película")
	}

	// Realizar soft delete
	if err := config.DB.Delete(&like).Error; err != nil {
		return errors.New("error al quitar me gusta de la película")
	}

	return nil
}

// GetLikeByUserAndMovie verifica si un usuario ha dado me gusta a una película
func GetLikeByUserAndMovie(userID, movieID uint) (bool, error) {
	var like models.Like
	// No necesitamos Unscoped() aquí porque GORM automáticamente ignora los registros con soft delete
	result := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).First(&like)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// GetLikesByUser obtiene todas las películas a las que un usuario ha dado me gusta
func GetLikesByUser(userID uint) ([]models.Movie, error) {
	var movies []models.Movie
	// Modificamos la consulta para asegurarnos de que solo incluya likes activos (no eliminados)
	err := config.DB.Joins("JOIN movie_likes ON movies.id = movie_likes.movie_id AND movie_likes.deleted_at IS NULL").
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

// DiagnoseLikes función auxiliar para diagnosticar problemas con los likes
func DiagnoseLikes(userID, movieID uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 1. Verificar likes activos
	var activeLikes []models.Like
	if err := config.DB.Where("user_id = ? AND movie_id = ?", userID, movieID).Find(&activeLikes).Error; err != nil {
		return nil, err
	}
	result["active_likes"] = activeLikes

	// 2. Verificar likes eliminados (soft deleted)
	var deletedLikes []models.Like
	if err := config.DB.Unscoped().Where("user_id = ? AND movie_id = ? AND deleted_at IS NOT NULL", userID, movieID).Find(&deletedLikes).Error; err != nil {
		return nil, err
	}
	result["deleted_likes"] = deletedLikes

	// 3. Verificar todos los likes sin filtro de borrado
	var allLikes []models.Like
	if err := config.DB.Unscoped().Where("user_id = ? AND movie_id = ?", userID, movieID).Find(&allLikes).Error; err != nil {
		return nil, err
	}
	result["all_likes"] = allLikes

	return result, nil
}
