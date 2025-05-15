package routes

import (
	"cine_conecta_backend/auth/middlewares"
	"cine_conecta_backend/comments/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(r *gin.Engine) {
	comments := r.Group("/api/comments")
	{
		comments.GET("/", controllers.GetComments)
		comments.GET("/:id", controllers.GetComment)

		// Sólo usuarios logueados pueden crear/editar/borrar sus comentarios
		comments.POST("/", middlewares.AuthRequired(), controllers.CreateComment)
		comments.PUT("/:id", middlewares.AuthRequired(), controllers.UpdateComment)
		comments.DELETE("/:id", middlewares.AuthRequired(), controllers.DeleteComment)

		// Ruta de actualización de todos los comentarios (sólo admin)
		comments.POST("/update-sentiments", middlewares.AdminRequired(), controllers.UpdateAllSentiments)
	}

	// Rutas para película-comentarios
	movies := r.Group("/api/movies")
	{
		movies.GET("/:id/comments", controllers.GetMovieComments)
		movies.GET("/:id/sentiment", controllers.GetMovieSentiment)
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
