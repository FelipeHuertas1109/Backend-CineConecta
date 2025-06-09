package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Estructura para recibir los datos de creación/actualización de películas
type MovieInput struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Director    string   `json:"director"`
	ReleaseDate string   `json:"release_date"`
	Rating      float32  `json:"rating"`
	PosterURL   string   `json:"poster_url"`
	Genres      []string `json:"genres"` // Lista de nombres de géneros
	Genre       string   `json:"genre"`  // Campo legacy para compatibilidad
}

// Método: POST /api/movies (restringido a admin)
func CreateMovie(c *gin.Context) {
	var input struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		Director    string    `json:"director"`
		ReleaseDate time.Time `json:"release_date"`
		Rating      float32   `json:"rating"`
		PosterURL   string    `json:"poster_url"`
		Genre       string    `json:"genre"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Crear la película
	movie := models.Movie{
		Title:       input.Title,
		Description: input.Description,
		Director:    input.Director,
		ReleaseDate: input.ReleaseDate,
		Rating:      input.Rating,
		PosterURL:   input.PosterURL,
	}

	// Procesar géneros desde el string del género
	var genreNames []string
	if input.Genre != "" {
		genreNames = models.ParseGenresString(input.Genre)
	}

	// Crear la película con géneros
	if err := services.CreateMovieWithGenres(&movie, genreNames); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la película: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// GetMovies devuelve todas las películas.
// Método: GET /api/movies
func GetMovies(c *gin.Context) {
	var movies []models.Movie

	// Realizar la consulta con preload de géneros
	query := config.DB.Model(&models.Movie{}).Preload("Genres")

	// Filtrar por género si se proporciona
	if genre := c.Query("genre"); genre != "" {
		// Buscar primero por el campo de texto del género para mayor rapidez
		query = query.Where("genre ILIKE ?", "%"+genre+"%")
		// También se podría hacer con una subconsulta más compleja para buscar en la tabla de géneros
	}

	// Aplicar ordenamiento si se especifica
	if sort := c.Query("sort"); sort != "" {
		direction := "ASC"
		if strings.HasPrefix(sort, "-") {
			direction = "DESC"
			sort = sort[1:]
		}
		query = query.Order(sort + " " + direction)
	} else {
		// Ordenar por ID de forma descendente por defecto
		query = query.Order("id DESC")
	}

	// Ejecutar la consulta
	if err := query.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetRecentMovies devuelve las películas más recientes según su fecha de lanzamiento.
// Método: GET /api/movies/recent
func GetRecentMovies(c *gin.Context) {
	// Obtener el parámetro 'limit' de la consulta (por defecto 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Si hay error o el límite es inválido, usar 10 por defecto
	}

	movies, err := services.GetRecentMovies(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudieron obtener las películas recientes")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":   len(movies),
		"results": movies,
	})
}

// GetMovie devuelve una película por su ID.
// Método: GET /api/movies/:movieId
func GetMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := config.DB.Preload("Genres").First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Película no encontrada"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// UpdateMovie actualiza una película existente.
// Método: PUT /api/movies/:movieId (restringido a admin)
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie models.Movie
	if err := config.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Película no encontrada"})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Director    string    `json:"director"`
		ReleaseDate time.Time `json:"release_date"`
		Rating      float32   `json:"rating"`
		PosterURL   string    `json:"poster_url"`
		Genre       string    `json:"genre"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos de la película
	if input.Title != "" {
		movie.Title = input.Title
	}
	if input.Description != "" {
		movie.Description = input.Description
	}
	if input.Director != "" {
		movie.Director = input.Director
	}
	if !input.ReleaseDate.IsZero() {
		movie.ReleaseDate = input.ReleaseDate
	}
	if input.Rating != 0 {
		movie.Rating = input.Rating
	}
	if input.PosterURL != "" {
		movie.PosterURL = input.PosterURL
	}

	// Procesar géneros desde el string
	var genreNames []string
	if input.Genre != "" {
		genreNames = models.ParseGenresString(input.Genre)
	}

	// Actualizar la película con géneros
	if err := services.UpdateMovieWithGenres(&movie, genreNames); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la película: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// DeleteMovie elimina una película por su ID.
// Método: DELETE /api/movies/:movieId (restringido a admin)
func DeleteMovie(c *gin.Context) {
	idParam := c.Param("movieId")
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
	idParam := c.Param("movieId")
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
