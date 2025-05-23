package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/comments/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /api/comments    (AuthRequired)
func CreateComment(c *gin.Context) {
	var input models.CommentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID

	// Buscar la película por nombre
	movie, err := services.FindMovieByName(input.MovieName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	// Crear el comentario con el ID de la película encontrada
	comment := &models.Comment{
		UserID:  userID,
		MovieID: movie.ID,
		Content: input.Content,
	}

	if err := services.CreateComment(comment); err != nil {
		// Verificar si es el error específico de usuario que ya ha comentado
		if strings.Contains(err.Error(), "ya ha comentado") {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo crear el comentario")
		return
	}

	// Generar texto descriptivo para la puntuación
	ratingText := getRatingDescription(comment.SentimentScore)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comentario creado correctamente",
		"comment": comment,
		"movie": gin.H{
			"id":    movie.ID,
			"title": movie.Title,
		},
		"sentiment_info": gin.H{
			"rating":         comment.SentimentScore,
			"description":    ratingText,
			"sentiment":      comment.Sentiment,
			"sentiment_text": getSentimentText(string(comment.Sentiment)),
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

// DELETE /api/comments/all (AdminRequired)
func DeleteAllComments(c *gin.Context) {
	// Verificación adicional de seguridad
	claims, _ := c.Get("claims")
	if claims.(*utils.Claims).Role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Solo administradores pueden ejecutar esta acción")
		return
	}

	// Ya no se requiere confirmación específica, solo ser administrador es suficiente

	if err := services.DeleteAllComments(); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar los comentarios: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todos los comentarios han sido eliminados correctamente",
		"time":    time.Now().Format(time.RFC3339),
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
