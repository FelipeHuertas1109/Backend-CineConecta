package routes

import (
	"cine_conecta_backend/auth/middlewares"
	commentsControllers "cine_conecta_backend/comments/controllers"
	"cine_conecta_backend/movies/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRecommendationRoutes(r *gin.Engine) {
	recommendations := r.Group("/api/recommendations")
	{
		// Recomendaciones personalizadas (requiere autenticación)
		recommendations.GET("/me", middlewares.AuthRequired(), controllers.GetMyRecommendations)

		// Recomendaciones populares basadas en el análisis de sentimientos (disponible para todos)
		recommendations.GET("/popular", middlewares.AuthRequired(), controllers.GetPopularRecommendations)

		// Rutas para el dataset de recomendaciones
		recommendations.POST("/dataset", middlewares.AuthRequired(), commentsControllers.SaveRecommendationsToDataset)
		recommendations.GET("/dataset", middlewares.AuthRequired(), commentsControllers.GetMyRecommendationDatasets)
		recommendations.GET("/dataset/:id", middlewares.AuthRequired(), commentsControllers.GetRecommendationDatasetByID)
		recommendations.DELETE("/dataset/:id", middlewares.AuthRequired(), commentsControllers.DeleteRecommendationDataset)
	}
}
