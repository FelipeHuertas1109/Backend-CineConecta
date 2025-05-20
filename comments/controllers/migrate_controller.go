package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/comments/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GET /api/comments/migrate (AdminRequired)
func RecomputeAllSentiments(c *gin.Context) {
	// Verificación adicional de seguridad
	claims, _ := c.Get("claims")
	if claims.(*utils.Claims).Role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "Solo administradores pueden ejecutar esta acción")
		return
	}

	// Iniciar tiempo de ejecución
	startTime := time.Now()

	// Ejecutar la migración
	stats, err := services.RecomputeAllSentiments()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al recalcular sentimientos: "+err.Error())
		return
	}

	// Calcular duración total
	duration := time.Since(startTime)

	// Preparar respuesta
	response := gin.H{
		"message":         "Sentimientos recalculados correctamente",
		"time":            time.Now().Format(time.RFC3339),
		"duration":        duration.String(),
		"total":           stats.Total,
		"processed":       stats.Processed,
		"failed":          stats.Failed,
		"changes_count":   stats.ChangesCount,
		"sentiment_stats": stats.SentimentStats,
		"score_ranges":    stats.ScoreRanges,
	}

	// Solo incluir detalles si hay cambios
	if stats.ChangesCount > 0 && len(stats.Changes) > 0 {
		response["changes_sample"] = stats.Changes
	}

	// Solo incluir IDs fallidos si hubo fallos
	if stats.Failed > 0 && len(stats.FailedIds) > 0 {
		response["failed_ids"] = stats.FailedIds
	}

	c.JSON(http.StatusOK, response)
}
