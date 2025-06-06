package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"errors"
	"fmt"
)

func CreateComment(c *models.Comment) error {
	fmt.Printf("[DEBUG-SERVICE] Creando comentario: UserID=%d, MovieID=%d, Content=%s\n", c.UserID, c.MovieID, c.Content)

	// Verificar si el usuario ya ha comentado esta película
	var existingComment models.Comment
	result := config.DB.Where("user_id = ? AND movie_id = ?", c.UserID, c.MovieID).First(&existingComment)

	if result.RowsAffected > 0 {
		fmt.Printf("[DEBUG-SERVICE] El usuario %d ya ha comentado la película %d\n", c.UserID, c.MovieID)
		return errors.New("el usuario ya ha comentado esta película")
	}

	// Usar la API de Cine Conecta ML para el análisis de sentimientos
	fmt.Printf("[DEBUG-SERVICE] Analizando sentimiento con API para: %s\n", c.Content)
	sentiment, score, err := AnalyzeSentimentWithML(c.Content)
	if err != nil {
		fmt.Printf("[DEBUG-SERVICE] Error al analizar sentimiento con API: %v\n", err)
		return fmt.Errorf("error al analizar sentimiento: %w", err)
	}

	fmt.Printf("[DEBUG-SERVICE] Sentimiento obtenido de API: %s, Score: %.2f\n", sentiment, score)

	c.Sentiment = sentiment
	c.SentimentScore = score

	// Crear el comentario en la base de datos
	fmt.Println("[DEBUG-SERVICE] Guardando comentario en la base de datos")
	err = config.DB.Create(c).Error
	if err != nil {
		fmt.Printf("[DEBUG-SERVICE] Error al guardar en la base de datos: %v\n", err)
		return fmt.Errorf("error al guardar en la base de datos: %w", err)
	}

	fmt.Printf("[DEBUG-SERVICE] Comentario guardado exitosamente con ID: %d\n", c.ID)
	return nil
}

func GetComments() ([]models.Comment, error) {
	var list []models.Comment
	err := config.DB.Preload("User").Preload("Movie").Find(&list).Error
	return list, err
}

func GetCommentByID(id uint) (models.Comment, error) {
	var c models.Comment
	err := config.DB.Preload("User").Preload("Movie").First(&c, id).Error
	return c, err
}

func UpdateComment(c *models.Comment) error {
	// Usar la API de Cine Conecta ML para el análisis de sentimientos
	fmt.Printf("[DEBUG-SERVICE] Analizando sentimiento con API para actualización: %s\n", c.Content)
	sentiment, score, err := AnalyzeSentimentWithML(c.Content)
	if err != nil {
		fmt.Printf("[DEBUG-SERVICE] Error al analizar sentimiento con API: %v\n", err)
		return fmt.Errorf("error al analizar sentimiento: %w", err)
	}

	fmt.Printf("[DEBUG-SERVICE] Sentimiento obtenido de API: %s, Score: %.2f\n", sentiment, score)

	c.Sentiment = sentiment
	c.SentimentScore = score

	return config.DB.Save(c).Error
}

func DeleteComment(id uint) error {
	return config.DB.Delete(&models.Comment{}, id).Error
}

// GetCommentsByMovie obtiene todos los comentarios de una película
func GetCommentsByMovie(movieID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Where("movie_id = ?", movieID).
		Preload("User").
		Find(&comments).Error
	return comments, err
}

// GetCommentsByMovieName obtiene todos los comentarios de una película por su nombre
func GetCommentsByMovieName(movieName string) ([]models.Comment, error) {
	var comments []models.Comment

	// Primero buscamos la película por nombre
	var movieIDs []uint
	err := config.DB.Model(&movieModels.Movie{}).
		Where("LOWER(title) LIKE LOWER(?)", "%"+movieName+"%").
		Pluck("id", &movieIDs).Error

	if err != nil {
		return nil, err
	}

	if len(movieIDs) == 0 {
		return []models.Comment{}, nil
	}

	// Luego buscamos los comentarios para esas películas
	err = config.DB.Where("movie_id IN ?", movieIDs).
		Preload("User").
		Preload("Movie").
		Find(&comments).Error

	return comments, err
}

// GetCommentsByUser obtiene todos los comentarios de un usuario
func GetCommentsByUser(userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Where("user_id = ?", userID).
		Preload("Movie").
		Find(&comments).Error
	return comments, err
}

// GetSentimentStats obtiene estadísticas de sentimientos para todas las películas
func GetSentimentStats() (map[string]int, error) {
	stats := map[string]int{
		"positive": 0,
		"neutral":  0,
		"negative": 0,
	}

	var comments []models.Comment
	if err := config.DB.Find(&comments).Error; err != nil {
		return stats, err
	}

	for _, comment := range comments {
		switch comment.Sentiment {
		case models.SentimentPositive:
			stats["positive"]++
		case models.SentimentNeutral:
			stats["neutral"]++
		case models.SentimentNegative:
			stats["negative"]++
		}
	}

	return stats, nil
}

// UpdateAllCommentSentiments actualiza el análisis de sentimientos para todos los comentarios existentes
func UpdateAllCommentSentiments() error {
	var comments []models.Comment

	// Obtener todos los comentarios
	if err := config.DB.Find(&comments).Error; err != nil {
		return err
	}

	fmt.Printf("[DEBUG-SERVICE] Actualizando sentimientos para %d comentarios\n", len(comments))

	// Actualizar cada comentario con el nuevo análisis de sentimientos
	for _, comment := range comments {
		// Usar la API de Cine Conecta ML para el análisis de sentimientos
		fmt.Printf("[DEBUG-SERVICE] Analizando sentimiento para comentario ID=%d: %s\n", comment.ID, comment.Content)
		sentiment, score, err := AnalyzeSentimentWithML(comment.Content)
		if err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al analizar sentimiento con API para comentario %d: %v\n", comment.ID, err)
			// Continuamos con el siguiente comentario en caso de error
			continue
		}

		fmt.Printf("[DEBUG-SERVICE] Sentimiento obtenido de API para comentario %d: %s, Score: %.2f\n", comment.ID, sentiment, score)
		comment.Sentiment = sentiment
		comment.SentimentScore = score

		if err := config.DB.Save(&comment).Error; err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al guardar comentario %d: %v\n", comment.ID, err)
			continue
		}

		fmt.Printf("[DEBUG-SERVICE] Comentario %d actualizado correctamente\n", comment.ID)
	}

	return nil
}

// DeleteAllComments elimina todos los comentarios de la base de datos
func DeleteAllComments() error {
	// Usar eliminación en masa para mayor eficiencia
	return config.DB.Exec("DELETE FROM comments").Error
}
