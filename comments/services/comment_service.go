package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	"errors"
	"os"
)

func CreateComment(c *models.Comment) error {
	// Verificar si el usuario ya ha comentado esta película
	var existingComment models.Comment
	result := config.DB.Where("user_id = ? AND movie_id = ?", c.UserID, c.MovieID).First(&existingComment)

	if result.RowsAffected > 0 {
		return errors.New("el usuario ya ha comentado esta película")
	}

	// Realizar análisis de sentimientos (con ML si está disponible)
	// Verificar si se debe usar análisis con ML (basado en variable de entorno)
	var sentiment models.SentimentType
	var score float64
	var err error

	if os.Getenv("USE_ML_SENTIMENT") == "true" {
		// Usar ML para el análisis
		sentiment, score, err = AnalyzeSentimentWithML(c.Content)
		if err != nil {
			// Fallback al método tradicional
			sentiment, score = AnalyzeSentiment(c.Content)
		}
	} else {
		// Usar método tradicional
		sentiment, score = AnalyzeSentiment(c.Content)
	}

	c.Sentiment = sentiment
	c.SentimentScore = score

	return config.DB.Create(c).Error
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
	// Realizar análisis de sentimientos (con ML si está disponible)
	// Verificar si se debe usar análisis con ML (basado en variable de entorno)
	var sentiment models.SentimentType
	var score float64
	var err error

	if os.Getenv("USE_ML_SENTIMENT") == "true" {
		// Usar ML para el análisis
		sentiment, score, err = AnalyzeSentimentWithML(c.Content)
		if err != nil {
			// Fallback al método tradicional
			sentiment, score = AnalyzeSentiment(c.Content)
		}
	} else {
		// Usar método tradicional
		sentiment, score = AnalyzeSentiment(c.Content)
	}

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

	useML := os.Getenv("USE_ML_SENTIMENT") == "true"

	// Actualizar cada comentario con el nuevo análisis de sentimientos
	for _, comment := range comments {
		var sentiment models.SentimentType
		var score float64
		var err error

		if useML {
			// Usar ML para el análisis
			sentiment, score, err = AnalyzeSentimentWithML(comment.Content)
			if err != nil {
				// Fallback al método tradicional
				sentiment, score = AnalyzeSentiment(comment.Content)
			}
		} else {
			// Usar método tradicional
			sentiment, score = AnalyzeSentiment(comment.Content)
		}

		comment.Sentiment = sentiment
		comment.SentimentScore = score

		if err := config.DB.Save(&comment).Error; err != nil {
			return err
		}
	}

	return nil
}
