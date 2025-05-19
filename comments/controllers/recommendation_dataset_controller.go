package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/comments/services"
	"cine_conecta_backend/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /api/recommendations/dataset
func SaveRecommendationsToDataset(c *gin.Context) {
	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID

	// Estructura para recibir los datos
	var input struct {
		Recommendations []models.RecommendationItem `json:"recommendations"`
		UserID          uint                        `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Verificar que el userID del token coincida con el enviado (o que sea admin)
	if userID != input.UserID && claims.(*utils.Claims).Role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "No autorizado para guardar recomendaciones de este usuario")
		return
	}

	// Verificar manualmente qué películas han sido comentadas por este usuario específico
	var moviesWithComments []models.RecommendationItem
	var moviesWithoutComments []string

	for _, rec := range input.Recommendations {
		var count int64
		if err := config.DB.Model(&models.Comment{}).
			Where("movie_id = ? AND user_id = ?", rec.MovieID, input.UserID).
			Count(&count).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al verificar comentarios: "+err.Error())
			return
		}

		if count > 0 {
			moviesWithComments = append(moviesWithComments, rec)
		} else {
			moviesWithoutComments = append(moviesWithoutComments, rec.Title)
		}
	}

	// Si no hay películas comentadas por este usuario, retornar error
	if len(moviesWithComments) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Ninguna de las películas seleccionadas ha sido comentada por este usuario")
		return
	}

	// Guardar en el dataset solo las películas con comentarios
	datasetID, err := services.SaveRecommendationsToDataset(input.UserID, moviesWithComments)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al guardar las recomendaciones: "+err.Error())
		return
	}

	// Preparar respuesta informativa
	response := gin.H{
		"message":        "Recomendaciones guardadas correctamente",
		"dataset_id":     datasetID,
		"saved_count":    len(moviesWithComments),
		"total_received": len(input.Recommendations),
		"user_id":        input.UserID,
		"filtered_count": len(moviesWithoutComments),
	}

	// Añadir información sobre películas guardadas
	if len(moviesWithComments) > 0 {
		var savedTitles []string
		for _, movie := range moviesWithComments {
			savedTitles = append(savedTitles, movie.Title)
		}
		response["saved_movies"] = savedTitles
		response["saved_message"] = "Películas guardadas correctamente (con comentarios del usuario)"
	}

	// Añadir información sobre películas descartadas si las hay
	if len(moviesWithoutComments) > 0 {
		response["skipped_movies"] = moviesWithoutComments
		response["skipped_reason"] = "Estas películas no han sido comentadas por este usuario y fueron excluidas"
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/recommendations/dataset
func GetMyRecommendationDatasets(c *gin.Context) {
	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID

	// Obtener datasets
	datasets, err := services.GetRecommendationDatasets(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener los datasets")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"datasets": datasets,
		"count":    len(datasets),
	})
}

// GET /api/recommendations/dataset/:id
func GetRecommendationDatasetByID(c *gin.Context) {
	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID
	isAdmin := claims.(*utils.Claims).Role == "admin"

	// Obtener ID del dataset
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	// Obtener dataset
	dataset, err := services.GetRecommendationDatasetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Dataset no encontrado")
		return
	}

	// Verificar que el dataset pertenezca al usuario (a menos que sea admin)
	if dataset.UserID != userID && !isAdmin {
		utils.ErrorResponse(c, http.StatusForbidden, "No autorizado para ver este dataset")
		return
	}

	c.JSON(http.StatusOK, dataset)
}

// DELETE /api/recommendations/dataset/:id
func DeleteRecommendationDataset(c *gin.Context) {
	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID
	isAdmin := claims.(*utils.Claims).Role == "admin"

	// Obtener ID del dataset
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}

	// Si es admin, puede eliminar cualquier dataset
	if isAdmin {
		if err := services.DeleteRecommendationDataset(uint(id), 0); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error al eliminar el dataset")
			return
		}
	} else {
		// Si es usuario normal, solo puede eliminar sus propios datasets
		if err := services.DeleteRecommendationDataset(uint(id), userID); err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, "No autorizado para eliminar este dataset o no existe")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dataset eliminado correctamente",
	})
}
