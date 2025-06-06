package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

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
	var input MovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Crear objeto película
	var movie models.Movie
	movie.Title = input.Title
	movie.Description = input.Description
	movie.Director = input.Director
	movie.Rating = input.Rating
	movie.PosterURL = input.PosterURL

	// Procesar fecha de lanzamiento
	if input.ReleaseDate != "" {
		releaseDate, err := utils.ParseDate(input.ReleaseDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Formato de fecha inválido")
			return
		}
		movie.ReleaseDate = releaseDate
	}

	// Determinar los géneros a usar
	var genreNames []string
	if len(input.Genres) > 0 {
		// Usar los géneros proporcionados en el array
		genreNames = input.Genres
	} else if input.Genre != "" {
		// Compatibilidad con el campo legacy
		genreNames = models.ParseGenresString(input.Genre)
		movie.Genre = input.Genre // Mantener el campo legacy
	}

	// Crear película con géneros
	if err := services.CreateMovieWithGenres(&movie, genreNames); err != nil {
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
	idParam := c.Param("movieId")
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
// Método: PUT /api/movies/:movieId (restringido a admin)
func UpdateMovie(c *gin.Context) {
	idParam := c.Param("movieId")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	// Obtener la película existente
	existingMovie, err := services.GetMovieByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Película no encontrada")
		return
	}

	// Recibir datos de actualización
	var input MovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Actualizar campos
	existingMovie.Title = input.Title
	existingMovie.Description = input.Description
	existingMovie.Director = input.Director
	existingMovie.Rating = input.Rating
	if input.PosterURL != "" {
		existingMovie.PosterURL = input.PosterURL
	}

	// Procesar fecha de lanzamiento
	if input.ReleaseDate != "" {
		releaseDate, err := utils.ParseDate(input.ReleaseDate)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Formato de fecha inválido")
			return
		}
		existingMovie.ReleaseDate = releaseDate
	}

	// Determinar los géneros a usar
	var genreNames []string
	if len(input.Genres) > 0 {
		// Usar los géneros proporcionados en el array
		genreNames = input.Genres
	} else if input.Genre != "" {
		// Compatibilidad con el campo legacy
		genreNames = models.ParseGenresString(input.Genre)
		existingMovie.Genre = input.Genre // Mantener el campo legacy
	}

	// Actualizar película con géneros
	if err := services.UpdateMovieWithGenres(&existingMovie, genreNames); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo actualizar la película")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Película actualizada correctamente",
		"movie":   existingMovie,
	})
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
