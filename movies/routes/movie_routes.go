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
		movies.GET("/", middlewares.AuthRequired(), controllers.GetMovies)
		movies.GET("/:id", middlewares.AuthRequired(), controllers.GetMovie)
		movies.GET("/sorted", middlewares.AuthRequired(), controllers.GetMoviesSorted)
		movies.GET("/recent", middlewares.AuthRequired(), controllers.GetRecentMovies)

		// Películas mejor valoradas
		movies.GET("/top-rated", middlewares.AuthRequired(), controllers.GetTopRatedMovies)
		movies.GET("/average-score", middlewares.AuthRequired(), controllers.GetMovieWithAverageScore)

		// Búsqueda avanzada
		movies.GET("/search", middlewares.AuthRequired(), controllers.SearchMovies)
		movies.GET("/genres", middlewares.AuthRequired(), controllers.GetGenres)
		movies.GET("/genres/detailed", middlewares.AuthRequired(), controllers.GetGenresDetailed)

		// Rutas restringidas a admin
		movies.POST("/", middlewares.AdminRequired(), controllers.CreateMovie)
		movies.PUT("/:id", middlewares.AdminRequired(), controllers.UpdateMovie)
		movies.DELETE("/:id", middlewares.AdminRequired(), controllers.DeleteMovie)
		movies.POST("/:id/poster", middlewares.AdminRequired(), controllers.UploadPoster)
	}

	// Registrar las rutas de recomendaciones
	RegisterRecommendationRoutes(r)
}
