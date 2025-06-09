package controllers

import (
	movieServices "cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTopRatedMovies devuelve las 5 películas mejor valoradas según los comentarios
func GetTopRatedMovies(c *gin.Context) {
	// Obtener películas mejor valoradas
	movies, err := movieServices.GetTopRatedMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Responder con las películas
	c.JSON(http.StatusOK, gin.H{
		"message": "Películas mejor valoradas por la comunidad",
		"count":   len(movies),
		"movies":  movies,
	})
}

// GetMovieWithAverageScore devuelve una película con su puntuación media de comentarios
func GetMovieWithAverageScore(c *gin.Context) {
	// Obtener ID de la película
	movieIDStr := c.Query("id")
	if movieIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere el ID de la película",
		})
		return
	}

	// Convertir ID a uint
	movieID, err := strconv.ParseUint(movieIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID de película inválido",
		})
		return
	}

	// Obtener película con puntuación media
	movie, err := movieServices.GetMovieWithAverageScore(uint(movieID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Responder con la película
	c.JSON(http.StatusOK, gin.H{
		"message": "Puntuación media de la película",
		"movie":   movie,
	})
}
