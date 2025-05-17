package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchMovies busca películas por título, género o puntuación
// GET /api/movies/search
func SearchMovies(c *gin.Context) {
	// Obtener parámetros de búsqueda
	title := c.Query("title")
	genre := c.Query("genre")
	ratingStr := c.Query("rating")

	// Convertir rating a float64
	var rating float64
	if ratingStr != "" {
		var err error
		rating, err = strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Puntuación inválida")
			return
		}
	}

	// Crear parámetros de búsqueda
	params := services.SearchParams{
		Title:  title,
		Genre:  genre,
		Rating: rating,
	}

	// Realizar la búsqueda
	movies, err := services.SearchMovies(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error en la búsqueda")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": movies,
		"count":   len(movies),
		"filters": params,
	})
}

// GetGenres devuelve todos los géneros disponibles (solo nombres)
// GET /api/movies/genres
func GetGenres(c *gin.Context) {
	genres, err := services.GetSimpleGenres()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener géneros")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
	})
}

// GetGenresDetailed devuelve todos los géneros con información detallada
// GET /api/movies/genres/detailed
func GetGenresDetailed(c *gin.Context) {
	genres, err := services.GetAllGenres()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener información de géneros")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"genres": genres,
		"count":  len(genres),
	})
}
