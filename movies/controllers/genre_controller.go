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
	genres, err := services.GetUniqueGenres()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener géneros")
		return
	}

	// Construir la respuesta con formato
	var formattedGenres []gin.H
	for _, genreName := range genres {
		formattedGenres = append(formattedGenres, gin.H{
			"name": genreName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": formattedGenres,
		"count":  len(genres),
	})
}

// GetGenreByName obtiene un género por su nombre
// GET /api/movies/genres/:name
func GetGenreByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nombre de género no proporcionado")
		return
	}

	// Obtener estadísticas del género
	stats, err := services.GetGenreStats(name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Género no encontrado")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  name,
		"stats": stats,
	})
}

// GetMovieGenres obtiene el género de una película
// GET /api/movies/:movieId/genres
func GetMovieGenres(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	genre, err := services.GetGenreForMovie(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Película no encontrada")
		return
	}

	// Si la película no tiene género, devolver un array vacío
	if genre == "" {
		c.JSON(http.StatusOK, gin.H{
			"genres": []gin.H{},
			"count":  0,
		})
		return
	}

	// Devolver el género en un array para mantener compatibilidad con el frontend
	c.JSON(http.StatusOK, gin.H{
		"genres": []gin.H{{"name": genre}},
		"count":  1,
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
		Genre string `json:"genre" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if err := services.AddGenreToMovie(uint(movieID), input.Genre); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al añadir género a la película")
		return
	}

	c.Status(http.StatusOK)
}

// RemoveGenreFromMovie elimina un género de una película
// DELETE /api/movies/:movieId/genres/:genre
func RemoveGenreFromMovie(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	genre := c.Param("genre")
	if genre == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Género no proporcionado")
		return
	}

	if err := services.RemoveGenreFromMovie(uint(movieID), genre); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar género de la película")
		return
	}

	c.Status(http.StatusOK)
}

// UpdateMovieGenre actualiza el género de una película
// PUT /api/movies/:movieId/genres
func UpdateMovieGenre(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("movieId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de película inválido")
		return
	}

	var input struct {
		Genre string `json:"genre" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if err := services.UpdateMovieGenre(uint(movieID), input.Genre); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al actualizar género de la película")
		return
	}

	c.Status(http.StatusOK)
}

// GetGenreStats obtiene estadísticas de un género
// GET /api/movies/genres/:name/stats
func GetGenreStats(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nombre de género no proporcionado")
		return
	}

	stats, err := services.GetGenreStats(name)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Género no encontrado")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetGenreInfoList obtiene lista de géneros con estadísticas
// GET /api/movies/genres/stats
func GetGenreInfoList(c *gin.Context) {
	genres, err := services.GetUniqueGenres()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener géneros")
		return
	}

	var stats []services.GenreInfo
	for _, genreName := range genres {
		genreStats, err := services.GetGenreStats(genreName)
		if err != nil {
			continue
		}
		stats = append(stats, *genreStats)
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": stats,
		"count":  len(stats),
	})
}
