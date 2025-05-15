package routes

import (
	"cine_conecta_backend/auth/middlewares"
	"cine_conecta_backend/movies/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRecommendationRoutes(r *gin.Engine) {
	recommendations := r.Group("/api/recommendations")
	{
		// Recomendaciones personalizadas (requiere autenticación)
		recommendations.GET("/me", middlewares.AuthRequired(), controllers.GetMyRecommendations)

		// Recomendaciones populares basadas en el análisis de sentimientos (disponible para todos)
		recommendations.GET("/popular", controllers.GetPopularRecommendations)
	}
}
