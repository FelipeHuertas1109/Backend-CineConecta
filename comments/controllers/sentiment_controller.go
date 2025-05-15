package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/services"
	"net/http"
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
