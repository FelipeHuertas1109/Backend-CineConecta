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
	}
}
