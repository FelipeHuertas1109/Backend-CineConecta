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

	c.JSON(http.StatusOK, gin.H{
		"message": "Comentario creado correctamente",
		"comment": input,
	})
}

// GET /api/comments
func GetComments(c *gin.Context) {
	list, err := services.GetComments()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener comentarios")
		return
	}
	c.JSON(http.StatusOK, list)
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
	c.JSON(http.StatusOK, comment)
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
		"message": "Comentario actualizado correctamente",
		"comment": input,
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
