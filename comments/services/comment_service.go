package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"errors"
	"fmt"
)

// UpdateMovieRating actualiza el rating de una película basado en los comentarios
// Esta función se mantiene por compatibilidad pero ya no es necesario llamarla
// directamente ya que los ratings se actualizan automáticamente al consultar las películas
func UpdateMovieRating(movieID uint) error {
	fmt.Printf("[DEBUG-SERVICE] Actualizando rating de la película ID=%d\n", movieID)

	// Verificar si la película existe
	var movie movieModels.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return fmt.Errorf("película no encontrada: %w", err)
	}

	// Obtener todos los comentarios de la película
	var comments []models.Comment
	if err := config.DB.Where("movie_id = ?", movieID).Find(&comments).Error; err != nil {
		return fmt.Errorf("error al obtener comentarios: %w", err)
	}

	// Calcular el nuevo rating basado en el promedio de las puntuaciones de sentimientos
	var newRating float32
	if len(comments) > 0 {
		var totalScore float64
		for _, comment := range comments {
			totalScore += comment.SentimentScore
		}
		newRating = float32(totalScore / float64(len(comments)))
		fmt.Printf("[DEBUG-SERVICE] Nuevo rating calculado para película %d: %.2f (basado en %d comentarios)\n",
			movieID, newRating, len(comments))
	} else {
		// Si no hay comentarios, dejar el rating en 0 o algún valor predeterminado
		newRating = 0
		fmt.Printf("[DEBUG-SERVICE] No hay comentarios para la película %d, rating establecido a 0\n", movieID)
	}

	// Actualizar el rating de la película
	movie.Rating = newRating
	if err := config.DB.Save(&movie).Error; err != nil {
		return fmt.Errorf("error al actualizar rating: %w", err)
	}

	fmt.Printf("[DEBUG-SERVICE] Rating de película %d actualizado exitosamente a %.2f\n", movieID, newRating)
	return nil
}

// updateAllMoviesRatings actualiza los ratings de todas las películas
func updateAllMoviesRatings() error {
	// Obtener todas las películas
	var movies []movieModels.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		return fmt.Errorf("error al obtener películas: %w", err)
	}

	fmt.Printf("[DEBUG-SERVICE] Actualizando ratings para %d películas\n", len(movies))

	// Actualizar el rating de cada película
	for _, movie := range movies {
		if err := UpdateMovieRating(movie.ID); err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al actualizar rating de película %d: %v\n", movie.ID, err)
			// Continuamos con la siguiente película en caso de error
			continue
		}
	}

	return nil
}

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

	err = config.DB.Save(c).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteComment(id uint) error {
	// Obtener el comentario primero para conocer la película asociada
	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		return err
	}

	// Eliminar el comentario
	if err := config.DB.Delete(&models.Comment{}, id).Error; err != nil {
		return err
	}

	return nil
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
	err := config.DB.Exec("DELETE FROM comments").Error
	if err != nil {
		return err
	}

	return nil
}
