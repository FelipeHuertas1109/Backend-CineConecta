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
		movies.GET("/:movieId", middlewares.AuthRequired(), controllers.GetMovie)
		movies.GET("/sorted", middlewares.AuthRequired(), controllers.GetMoviesSorted)
		movies.GET("/recent", middlewares.AuthRequired(), controllers.GetRecentMovies)

		// Películas mejor valoradas
		movies.GET("/top-rated", middlewares.AuthRequired(), controllers.GetTopRatedMovies)
		movies.GET("/average-score", middlewares.AuthRequired(), controllers.GetMovieWithAverageScore)

		// Búsqueda avanzada
		movies.GET("/search", middlewares.AuthRequired(), controllers.SearchMovies)

		// Rutas para "me gusta"
		movies.GET("/liked", middlewares.AuthRequired(), controllers.GetLikedMovies)
		movies.GET("/:movieId/like", middlewares.AuthRequired(), controllers.CheckLikeStatus)
		movies.POST("/:movieId/like", middlewares.AuthRequired(), controllers.LikeMovie)
		movies.DELETE("/:movieId/like", middlewares.AuthRequired(), controllers.UnlikeMovie)
		movies.GET("/:movieId/likes/count", middlewares.AuthRequired(), controllers.GetMovieLikes)
		movies.GET("/:movieId/like/diagnose", middlewares.AuthRequired(), controllers.DiagnoseLikesHandler)

		// Rutas para géneros
		movies.GET("/genres", middlewares.AuthRequired(), controllers.GetAllGenres)
		movies.GET("/genres/detailed", middlewares.AuthRequired(), controllers.GetGenreInfoList)
		movies.GET("/genres/:id", middlewares.AuthRequired(), controllers.GetGenreByID)
		movies.GET("/genres/:id/stats", middlewares.AuthRequired(), controllers.GetGenreStats)
		movies.POST("/genres", middlewares.AdminRequired(), controllers.CreateGenre)

		// Rutas para géneros de películas específicas
		movies.GET("/:movieId/genres", middlewares.AuthRequired(), controllers.GetMovieGenres)
		movies.POST("/:movieId/genres", middlewares.AdminRequired(), controllers.AddGenreToMovie)
		movies.PUT("/:movieId/genres", middlewares.AdminRequired(), controllers.UpdateMovieGenres)
		movies.DELETE("/:movieId/genres/:genreId", middlewares.AdminRequired(), controllers.RemoveGenreFromMovie)

		// Rutas restringidas a admin
		movies.POST("/", middlewares.AdminRequired(), controllers.CreateMovie)
		movies.PUT("/:movieId", middlewares.AdminRequired(), controllers.UpdateMovie)
		movies.DELETE("/:movieId", middlewares.AdminRequired(), controllers.DeleteMovie)
		movies.POST("/:movieId/poster", middlewares.AdminRequired(), controllers.UploadPoster)
	}

	// Registrar las rutas de recomendaciones
	RegisterRecommendationRoutes(r)
}
