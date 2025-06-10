package routes

import (
	"cine_conecta_backend/auth/middlewares"
	"cine_conecta_backend/comments/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(r *gin.Engine) {
	comments := r.Group("/api/comments")
	{
		comments.GET("/", middlewares.AuthRequired(), controllers.GetComments)
		comments.GET("/:id", middlewares.AuthRequired(), controllers.GetComment)

		// Nueva ruta para listar películas disponibles (sin autenticación para facilitar la depuración)
		comments.GET("/available-movies", controllers.ListAvailableMovies)

		// Nueva ruta para verificar si una película existe
		comments.GET("/check-movie/:id", controllers.CheckMovieExists)

		// Sólo usuarios logueados pueden crear/editar/borrar sus comentarios
		comments.POST("/", middlewares.AuthRequired(), controllers.CreateComment)
		comments.PUT("/:id", middlewares.AuthRequired(), controllers.UpdateComment)
		comments.DELETE("/:id", middlewares.AuthRequired(), controllers.DeleteComment)

		// Ruta de actualización de todos los comentarios (sólo admin)
		comments.POST("/update-sentiments", middlewares.AdminRequired(), controllers.UpdateAllSentiments)

		// Ruta para actualizar todos los ratings de películas (sólo admin)
		comments.POST("/update-ratings", middlewares.AuthRequired(), controllers.UpdateAllMovieRatings)

		// Rutas para configuración del análisis de sentimientos (sólo admin)
		comments.GET("/settings", middlewares.AdminRequired(), controllers.GetSentimentSettings)
		comments.POST("/settings", middlewares.AdminRequired(), controllers.UpdateSentimentSettings)

		// Ruta para eliminar TODOS los comentarios (sólo admin)
		comments.DELETE("/all", middlewares.AdminRequired(), controllers.DeleteAllComments)

		comments.GET("/migrate", middlewares.AdminRequired(), controllers.RecomputeAllSentiments)
	}

	// Rutas para película-comentarios por nombre
	moviesByName := r.Group("/api/movies")
	{
		// Rutas públicas sin autenticación
		moviesByName.GET("/public-comments/:name", controllers.GetPublicMovieCommentsByName)
		moviesByName.GET("/public-sentiment/:name", controllers.GetPublicMovieSentimentByName)
	}

	// Rutas para película-comentarios por ID (protegidas)
	moviesById := r.Group("/api/movies-by-id")
	{
		// Rutas protegidas que requieren autenticación
		moviesById.GET("/:id/comments", middlewares.AuthRequired(), controllers.GetMovieComments)
		moviesById.GET("/:id/sentiment", middlewares.AuthRequired(), controllers.GetMovieSentiment)
	}

	// Rutas para usuario-comentarios
	users := r.Group("/api/users")
	{
		users.GET("/:id/comments", middlewares.AuthRequired(), controllers.GetUserComments)
		users.GET("/:id/recommendations", middlewares.AuthRequired(), controllers.GetUserRecommendations)
	}

	// Rutas para estadísticas (sólo admin)
	sentiment := r.Group("/api/sentiment")
	{
		sentiment.GET("/stats", middlewares.AdminRequired(), controllers.GetSentimentStats)
	}
}
