package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllGenres obtiene todos los géneros
// GET /api/movies/genres
func GetAllGenres(c *gin.Context) {
	genres, err := services.GetAllGenres()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener géneros")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
		"count":  len(genres),
	})
}

// GetGenreByID obtiene un género por su ID
// GET /api/movies/genres/:id
func GetGenreByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de género inválido")
		return
	}

	genre, err := services.GetGenreByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Género no encontrado")
		return
	}

	c.JSON(http.StatusOK, genre)
}

// CreateGenre crea un nuevo género
// POST /api/movies/genres
func CreateGenre(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	genre, err := services.CreateGenre(input.Name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al crear género")
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// GetMovieGenres obtiene los géneros de una película
// GET /api/movies/:movieId/genres
func GetMovieGenres(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	genres, err := services.GetGenresForMovie(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Película no encontrada")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
		"count":  len(genres),
	})
}

// AddGenreToMovie añade un género a una película
// POST /api/movies/:movieId/genres
func AddGenreToMovie(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	var input struct {
		GenreID uint `json:"genre_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if err := services.AddGenreToMovie(uint(movieID), input.GenreID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al añadir género a la película")
		return
	}

	c.Status(http.StatusOK)
}

// RemoveGenreFromMovie elimina un género de una película
// DELETE /api/movies/:movieId/genres/:genreId
func RemoveGenreFromMovie(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	genreID, err := strconv.ParseUint(c.Param("genreId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de género inválido")
		return
	}

	if err := services.RemoveGenreFromMovie(uint(movieID), uint(genreID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar género de la película")
		return
	}

	c.Status(http.StatusOK)
}

// UpdateMovieGenres actualiza todos los géneros de una película
// PUT /api/movies/:movieId/genres
func UpdateMovieGenres(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	var input struct {
		GenreIDs []uint `json:"genre_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if err := services.UpdateMovieGenres(uint(movieID), input.GenreIDs); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar géneros de la película")
		return
	}

	c.Status(http.StatusOK)
}

// GetGenreStats obtiene estadísticas de un género
// GET /api/movies/genres/:id/stats
func GetGenreStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de género inválido")
		return
	}

	stats, err := services.GetGenreStats(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Género no encontrado")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetGenreInfoList obtiene lista de géneros con estadísticas
// GET /api/movies/genres/stats
func GetGenreInfoList(c *gin.Context) {
	stats, err := services.GetGenreInfoList()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener estadísticas de géneros")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": stats,
		"count":  len(stats),
	})
}
