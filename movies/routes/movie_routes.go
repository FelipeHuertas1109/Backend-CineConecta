package routes

import (
	"cine_conecta_backend/auth/middlewares"
	"cine_conecta_backend/movies/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterMovieRoutes(r *gin.Engine) {
	movies := r.Group("/api/movies")
	{
		// Lectura pública o con autenticación básica
		movies.GET("/", controllers.GetMovies)
		movies.GET("/:id", controllers.GetMovie)
		movies.GET("/sorted", controllers.GetMoviesSorted)
		movies.GET("/recent", controllers.GetRecentMovies)

		// Búsqueda avanzada
		movies.GET("/search", controllers.SearchMovies)
		movies.GET("/genres", controllers.GetGenres)
		movies.GET("/genres/detailed", controllers.GetGenresDetailed)

		// Rutas restringidas a admin
		movies.POST("/", middlewares.AdminRequired(), controllers.CreateMovie)
		movies.PUT("/:id", middlewares.AdminRequired(), controllers.UpdateMovie)
		movies.DELETE("/:id", middlewares.AdminRequired(), controllers.DeleteMovie)
		movies.POST("/:id/poster", middlewares.AdminRequired(), controllers.UploadPoster)
	}

	// Registrar las rutas de recomendaciones
	RegisterRecommendationRoutes(r)
}
