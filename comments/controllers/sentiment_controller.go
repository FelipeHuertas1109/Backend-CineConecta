package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/services"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/movies/:id/sentiment
func GetMovieSentiment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	sentiment, score, err := services.GetMovieSentiment(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener sentimiento")
		return
	}

	// Mapear el tipo de sentimiento a un texto descriptivo
	sentimentText := "neutro"
	if sentiment == "positive" {
		sentimentText = "positivo"
	} else if sentiment == "negative" {
		sentimentText = "negativo"
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id":       id,
		"sentiment":      sentiment,
		"sentiment_text": sentimentText,
		"rating":         score,
		"rating_text":    getPuntuacionTexto(score),
	})
}

// GET /api/movies/:id/comments
func GetMovieComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	comments, err := services.GetCommentsByMovie(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GET /api/users/:id/comments
func GetUserComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	comments, err := services.GetCommentsByUser(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GET /api/sentiment/stats   (AdminRequired)
func GetSentimentStats(c *gin.Context) {
	stats, err := services.GetSentimentStats()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sentiment_stats": stats,
	})
}

// POST /api/comments/settings (AdminRequired)
func UpdateSentimentSettings(c *gin.Context) {
	// Solo administradores pueden modificar esta configuración
	claims, _ := c.Get("claims")
	if claims.(*utils.Claims).Role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Solo administradores pueden cambiar esta configuración")
		return
	}

	// Estructura para recibir la configuración
	var settings struct {
		UseML     bool   `json:"use_ml"`
		OpenAIKey string `json:"openai_key,omitempty"`
	}

	if err := c.ShouldBindJSON(&settings); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Actualizar configuraciones
	if settings.UseML {
		os.Setenv("USE_ML_SENTIMENT", "true")
	} else {
		os.Setenv("USE_ML_SENTIMENT", "false")
	}

	// Actualizar API key si se proporciona
	if settings.OpenAIKey != "" {
		os.Setenv("OPENAI_API_KEY", settings.OpenAIKey)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuración de análisis de sentimientos actualizada",
		"use_ml":  settings.UseML,
	})
}

// GET /api/comments/settings (AdminRequired)
func GetSentimentSettings(c *gin.Context) {
	// Solo administradores pueden ver esta configuración
	claims, _ := c.Get("claims")
	if claims.(*utils.Claims).Role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Solo administradores pueden ver esta configuración")
		return
	}

	// Obtener configuración actual
	useML := os.Getenv("USE_ML_SENTIMENT") == "true"
	hasOpenAIKey := os.Getenv("OPENAI_API_KEY") != ""

	c.JSON(http.StatusOK, gin.H{
		"use_ml":         useML,
		"has_openai_key": hasOpenAIKey,
	})
}

// Función auxiliar para obtener descripción textual de la puntuación
func getPuntuacionTexto(score float64) string {
	switch {
	case score >= 9.5:
		return "Obra maestra"
	case score >= 9.0:
		return "Excepcional"
	case score >= 8.0:
		return "Excelente"
	case score >= 7.0:
		return "Muy buena"
	case score >= 6.0:
		return "Buena"
	case score >= 5.0:
		return "Aceptable"
	case score >= 4.0:
		return "Regular"
	case score >= 3.0:
		return "Mala"
	case score >= 2.0:
		return "Muy mala"
	default:
		return "Pésima"
	}
}

// GetPublicMovieComments obtiene todos los comentarios de una película sin requerir autenticación
func GetPublicMovieComments(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	comments, err := services.GetCommentsByMovie(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}

	// Añadir información enriquecida de sentimiento para cada comentario
	var enhancedComments []gin.H
	for _, comment := range comments {
		enhancedComments = append(enhancedComments, gin.H{
			"id":             comment.ID,
			"user_id":        comment.UserID,
			"movie_id":       comment.MovieID,
			"content":        comment.Content,
			"created_at":     comment.CreatedAt,
			"sentiment":      comment.Sentiment,
			"sentiment_text": getSentimentText(string(comment.Sentiment)),
			"rating":         comment.SentimentScore,
			"rating_text":    getPuntuacionTexto(comment.SentimentScore),
			"user":           comment.User,
		})
	}

	c.JSON(http.StatusOK, enhancedComments)
}

// GetPublicMovieSentiment obtiene el sentimiento de una película sin requerir autenticación
func GetPublicMovieSentiment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	sentiment, score, err := services.GetMovieSentiment(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener sentimiento")
		return
	}

	// Mapear el tipo de sentimiento a un texto descriptivo
	sentimentText := "neutro"
	if sentiment == "positive" {
		sentimentText = "positivo"
	} else if sentiment == "negative" {
		sentimentText = "negativo"
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id":       id,
		"sentiment":      sentiment,
		"sentiment_text": sentimentText,
		"rating":         score,
		"rating_text":    getPuntuacionTexto(score),
	})
}

// GetPublicMovieCommentsByName obtiene todos los comentarios de una película por su nombre sin requerir autenticación
func GetPublicMovieCommentsByName(c *gin.Context) {
	movieName := c.Param("name")
	if movieName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nombre de película inválido")
		return
	}

	comments, err := services.GetCommentsByMovieName(movieName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}

	// Añadir información enriquecida de sentimiento para cada comentario
	// y limitar información del usuario a solo ID y nombre
	var enhancedComments []gin.H
	for _, comment := range comments {
		enhancedComments = append(enhancedComments, gin.H{
			"id":             comment.ID,
			"user_id":        comment.UserID,
			"movie_id":       comment.MovieID,
			"content":        comment.Content,
			"created_at":     comment.CreatedAt,
			"sentiment":      comment.Sentiment,
			"sentiment_text": getSentimentText(string(comment.Sentiment)),
			"rating":         comment.SentimentScore,
			"rating_text":    getPuntuacionTexto(comment.SentimentScore),
			"user": gin.H{
				"id":   comment.User.ID,
				"name": comment.User.Name,
			},
		})
	}

	c.JSON(http.StatusOK, enhancedComments)
}

// GetPublicMovieSentimentByName obtiene el sentimiento de una película por su nombre sin requerir autenticación
func GetPublicMovieSentimentByName(c *gin.Context) {
	movieName := c.Param("name")
	if movieName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nombre de película inválido")
		return
	}

	sentiment, score, err := services.GetMovieSentimentByName(movieName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener sentimiento")
		return
	}

	// Mapear el tipo de sentimiento a un texto descriptivo
	sentimentText := "neutro"
	if sentiment == "positive" {
		sentimentText = "positivo"
	} else if sentiment == "negative" {
		sentimentText = "negativo"
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_name":     movieName,
		"sentiment":      sentiment,
		"sentiment_text": sentimentText,
		"rating":         score,
		"rating_text":    getPuntuacionTexto(score),
	})
}
