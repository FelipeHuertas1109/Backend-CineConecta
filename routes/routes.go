package routes

import (
	"cine_conecta_backend/controllers"
	"cine_conecta_backend/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		// Solo accesible para admin
		api.GET("/users", middlewares.AdminRequired(), controllers.GetAllUsers)
	}
}
