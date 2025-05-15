package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/comments/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /api/comments    (AuthRequired)
func CreateComment(c *gin.Context) {
	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Obtener userID del token, si quieres forzar que sólo usuarios logueados comenten
	claims, _ := c.Get("claims")
	input.UserID = claims.(*utils.Claims).UserID

	if err := services.CreateComment(&input); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo crear el comentario")
		return
	}

	// Generar texto descriptivo para la puntuación
	ratingText := getRatingDescription(input.SentimentScore)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comentario creado correctamente",
		"comment": input,
		"sentiment_info": gin.H{
			"rating":         input.SentimentScore,
			"description":    ratingText,
			"sentiment":      input.Sentiment,
			"sentiment_text": getSentimentText(string(input.Sentiment)),
		},
	})
}

// GET /api/comments
func GetComments(c *gin.Context) {
	list, err := services.GetComments()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}

	// Añadir información de puntuación para cada comentario
	var enhancedComments []gin.H
	for _, comment := range list {
		enhancedComments = append(enhancedComments, gin.H{
			"id":             comment.ID,
			"user_id":        comment.UserID,
			"movie_id":       comment.MovieID,
			"content":        comment.Content,
			"created_at":     comment.CreatedAt,
			"updated_at":     comment.UpdatedAt,
			"user":           comment.User,
			"movie":          comment.Movie,
			"rating":         comment.SentimentScore,
			"rating_text":    getRatingDescription(comment.SentimentScore),
			"sentiment":      comment.Sentiment,
			"sentiment_text": getSentimentText(string(comment.Sentiment)),
		})
	}

	c.JSON(http.StatusOK, enhancedComments)
}

// GET /api/comments/:id
func GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	comment, err := services.GetCommentByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Comentario no encontrado")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comment":        comment,
		"rating":         comment.SentimentScore,
		"rating_text":    getRatingDescription(comment.SentimentScore),
		"sentiment":      comment.Sentiment,
		"sentiment_text": getSentimentText(string(comment.Sentiment)),
	})
}

// PUT /api/comments/:id   (AuthRequired)
func UpdateComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}
	input.ID = uint(id)

	if err := services.UpdateComment(&input); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo actualizar el comentario")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "Comentario actualizado correctamente",
		"comment":     input,
		"rating":      input.SentimentScore,
		"rating_text": getRatingDescription(input.SentimentScore),
	})
}

// DELETE /api/comments/:id  (AuthRequired o AdminRequired)
func DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := services.DeleteComment(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo eliminar el comentario")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comentario eliminado correctamente"})
}

// POST /api/comments/update-sentiments (AdminRequired)
func UpdateAllSentiments(c *gin.Context) {
	err := services.UpdateAllCommentSentiments()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar sentimientos de comentarios")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sentimientos de comentarios actualizados correctamente",
	})
}

// Función auxiliar para obtener descripción textual de la puntuación
func getRatingDescription(score float64) string {
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

// Función auxiliar para convertir tipo de sentimiento a texto descriptivo
func getSentimentText(sentiment string) string {
	switch sentiment {
	case "positive":
		return "positivo"
	case "negative":
		return "negativo"
	default:
		return "neutro"
	}
}
