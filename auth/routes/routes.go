package routes

import (
	"cine_conecta_backend/auth/controllers"
	"cine_conecta_backend/auth/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
		api.POST("/logout", middlewares.AuthRequired(), controllers.Logout)
		api.GET("/profile", middlewares.AuthRequired(), controllers.GetProfile)
		api.POST("/forgot-password", controllers.ForgotPassword) // Solicitar restablecimiento
		api.POST("/reset-password", controllers.ResetPassword)   // Restablecer contraseña
		// Solo accesible para admin
		api.GET("/users", middlewares.AdminRequired(), controllers.GetAllUsers)
		api.DELETE("/users", middlewares.AdminRequired(), controllers.DeleteAllUsers)
		api.GET("/verify-token", middlewares.AuthRequired(), controllers.VerifyToken)
	}
}
