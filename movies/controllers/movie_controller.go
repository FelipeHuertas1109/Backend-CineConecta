package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Método: POST /api/movies (restringido a admin)
func CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	if err := services.CreateMovie(&movie); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo crear la película")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Película creada correctamente",
		"movie":   movie,
	})
}

// GetMovies devuelve todas las películas.
// Método: GET /api/movies
func GetMovies(c *gin.Context) {
	movies, err := services.GetMovies()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudieron obtener las películas")
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetMovie devuelve una película por su ID.
// Método: GET /api/movies/:id
func GetMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	movie, err := services.GetMovieByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Película no encontrada")
		return
	}
	c.JSON(http.StatusOK, movie)
}

// UpdateMovie actualiza una película existente.
// Método: PUT /api/movies/:id (restringido a admin)
func UpdateMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}
	movie.ID = uint(id)

	if err := services.UpdateMovie(&movie); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo actualizar la película")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Película actualizada correctamente",
		"movie":   movie,
	})
}

// DeleteMovie elimina una película por su ID.
// Método: DELETE /api/movies/:id (restringido a admin)
func DeleteMovie(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := services.DeleteMovie(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo eliminar la película")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Película eliminada correctamente"})
}

func GetMoviesSorted(c *gin.Context) {
	sortBy := c.DefaultQuery("sortBy", "title")
	order := c.DefaultQuery("order", "asc")

	movies, err := services.GetMoviesSorted(sortBy, order)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, movies)
}

func UploadPoster(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	// Verificar si la película existe
	movie, err := services.GetMovieByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Película no encontrada")
		return
	}

	fileHeader, err := c.FormFile("poster")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Archivo no encontrado (campo 'poster')")
		return
	}

	// Validar tamaño y tipo de archivo
	if fileHeader.Size > 50<<20 { // 50 MB
		utils.ErrorResponse(c, http.StatusBadRequest, "El archivo supera los 50 MB")
		return
	}

	mime := fileHeader.Header.Get("Content-Type")
	if mime != "image/jpeg" && mime != "image/png" && mime != "image/webp" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Formato no permitido. Sólo se aceptan JPEG, PNG o WEBP")
		return
	}

	url, err := services.UploadPoster(uint(id), fileHeader)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Devolver respuesta con la película actualizada
	c.JSON(http.StatusOK, gin.H{
		"message":    "Póster subido correctamente",
		"poster_url": url,
		"movie":      movie,
	})
}
