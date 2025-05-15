package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/users/:id/recommendations
func GetUserRecommendations(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de usuario inválido")
		return
	}

	raw, ok := c.Get("claims")
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Token no encontrado")
		return
	}
	claims := raw.(*utils.Claims)

	// solo el propio usuario o rol admin
	if claims.Role != "admin" && claims.UserID != uint(id) {
		utils.ErrorResponse(c, http.StatusForbidden, "No autorizado")
		return
	}

	recs, err := services.GetRecommendationsForUser(uint(id), 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error generando recomendaciones")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":         id,
		"recommendations": recs,
	})
}
