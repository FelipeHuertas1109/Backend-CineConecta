package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

// DeleteAllComments elimina todos los comentarios de la base de datos
func DeleteAllComments() error {
	// Usar eliminación en masa para mayor eficiencia
	return config.DB.Exec("DELETE FROM comments").Error
}

// FindMovieByName busca una película por nombre, primero buscando coincidencia exacta
// y luego por coincidencia parcial si no se encuentra
func FindMovieByName(movieName string) (movieModels.Movie, error) {
	var movie movieModels.Movie

	// Limpiar el nombre de la película
	cleanName := strings.TrimSpace(movieName)
	if cleanName == "" {
		return movie, errors.New("el nombre de la película no puede estar vacío")
	}

	log.Printf("🔍 Buscando película: '%s'", cleanName)

	// 1. Intentar coincidencia exacta (ignorando mayúsculas/minúsculas)
	exactResult := config.DB.Where("LOWER(title) = LOWER(?)", cleanName).First(&movie)
	if exactResult.Error == nil {
		log.Printf("✅ Película encontrada por coincidencia exacta: ID=%d, Título='%s'", movie.ID, movie.Title)
		return movie, nil
	}

	// 2. Intentar coincidencia parcial
	partialResult := config.DB.Where("LOWER(title) LIKE LOWER(?)", "%"+cleanName+"%").First(&movie)
	if partialResult.Error == nil {
		log.Printf("✅ Película encontrada por coincidencia parcial: ID=%d, Título='%s'", movie.ID, movie.Title)
		return movie, nil
	}

	// 3. Mostrar todas las películas disponibles para depuración
	var allMovies []movieModels.Movie
	config.DB.Select("id, title").Find(&allMovies)
	var titles []string
	for _, m := range allMovies {
		titles = append(titles, fmt.Sprintf("%d: %s", m.ID, m.Title))
	}
	log.Printf("📋 Películas disponibles: %s", strings.Join(titles, ", "))

	// No se encontró ninguna película
	log.Printf("❌ No se encontró ninguna película con el nombre: '%s'", cleanName)
	return movie, fmt.Errorf("no se encontró ninguna película con el nombre: '%s'", cleanName)
}
