package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LikeMovie permite a un usuario dar "me gusta" a una película
// POST /api/movies/:movieId/like
func LikeMovie(c *gin.Context) {
	// Obtener ID de la película
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	// Obtener ID del usuario del token
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// Crear el "me gusta"
	if err := services.CreateLike(userID, uint(movieID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Has dado me gusta a la película",
	})
}

// UnlikeMovie permite a un usuario quitar su "me gusta" de una película
// DELETE /api/movies/:movieId/like
func UnlikeMovie(c *gin.Context) {
	// Obtener ID de la película
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	// Obtener ID del usuario del token
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// Eliminar el "me gusta"
	if err := services.DeleteLike(userID, uint(movieID)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Has quitado tu me gusta de la película",
	})
}

// CheckLikeStatus verifica si un usuario ha dado "me gusta" a una película
// GET /api/movies/:movieId/like
func CheckLikeStatus(c *gin.Context) {
	// Obtener ID de la película
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	// Obtener ID del usuario del token
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// Verificar el estado del "me gusta"
	liked, err := services.GetLikeByUserAndMovie(userID, uint(movieID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al verificar estado del me gusta")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
	})
}

// GetLikedMovies obtiene todas las películas que un usuario ha marcado con "me gusta"
// GET /api/movies/liked
func GetLikedMovies(c *gin.Context) {
	// Obtener ID del usuario del token
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// Obtener películas con "me gusta"
	movies, err := services.GetLikesByUser(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener películas con me gusta")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"count":  len(movies),
	})
}

// GetMovieLikes obtiene la cantidad de "me gusta" de una película
// GET /api/movies/:movieId/likes/count
func GetMovieLikes(c *gin.Context) {
	// Obtener ID de la película
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	// Obtener cantidad de "me gusta"
	count, err := services.GetLikesByMovie(uint(movieID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener cantidad de me gusta")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

// DiagnoseLikesHandler obtiene información de diagnóstico sobre los likes
// GET /api/movies/:movieId/like/diagnose
func DiagnoseLikesHandler(c *gin.Context) {
	// Obtener ID de la película
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	// Obtener ID del usuario del token
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}
	userClaims := claims.(*utils.Claims)
	userID := userClaims.UserID

	// Obtener diagnóstico
	diagnosis, err := services.DiagnoseLikes(userID, uint(movieID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al diagnosticar likes")
		return
	}

	c.JSON(http.StatusOK, diagnosis)
}
