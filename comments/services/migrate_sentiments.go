// comments/services/migrate_sentiments.go
package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	"fmt"
	"time"
)

// MigrationStats contiene estadísticas sobre la migración
type MigrationStats struct {
	Total          int                          `json:"total"`
	Processed      int                          `json:"processed"`
	Failed         int                          `json:"failed"`
	ChangesCount   int                          `json:"changes_count"`
	StartTime      time.Time                    `json:"start_time"`
	EndTime        time.Time                    `json:"end_time"`
	Duration       string                       `json:"duration"`
	SentimentStats map[models.SentimentType]int `json:"sentiment_stats"`
	ScoreRanges    map[string]int               `json:"score_ranges"`
	FailedIds      []uint                       `json:"failed_ids,omitempty"`
	Changes        []map[string]interface{}     `json:"changes,omitempty"`
}

// RecomputeAllSentiments recalcula el sentimiento para todos los comentarios
func RecomputeAllSentiments() (*MigrationStats, error) {
	// Inicializar estadísticas
	stats := &MigrationStats{
		StartTime: time.Now(),
		SentimentStats: map[models.SentimentType]int{
			models.SentimentPositive: 0,
			models.SentimentNeutral:  0,
			models.SentimentNegative: 0,
		},
		ScoreRanges: map[string]int{
			"1-2":  0, // Muy negativo
			"2-3":  0, // Negativo
			"3-4":  0, // Algo negativo
			"4-5":  0, // Neutral bajo
			"5-6":  0, // Neutral alto
			"6-7":  0, // Algo positivo
			"7-8":  0, // Positivo
			"8-9":  0, // Muy positivo
			"9-10": 0, // Extremadamente positivo
		},
		FailedIds: []uint{},
		Changes:   []map[string]interface{}{},
	}

	// Obtener todos los comentarios
	var comments []models.Comment
	if err := config.DB.Find(&comments).Error; err != nil {
		return stats, err
	}

	stats.Total = len(comments)

	// Procesar cada comentario
	for i := range comments {
		stats.Processed++

		// Guardar valores originales
		originalSentiment := comments[i].Sentiment
		originalScore := comments[i].SentimentScore

		// Recalcular sentimiento
		newSentiment, newScore := AnalyzeSentiment(comments[i].Content)

		// Registrar cambio si hay diferencia
		if originalSentiment != newSentiment || originalScore != newScore {
			stats.ChangesCount++

			// Limitar el número de cambios detallados guardados para evitar respuestas enormes
			if len(stats.Changes) < 50 {
				stats.Changes = append(stats.Changes, map[string]interface{}{
					"id":                comments[i].ID,
					"content":           comments[i].Content,
					"old_sentiment":     originalSentiment,
					"new_sentiment":     newSentiment,
					"old_score":         originalScore,
					"new_score":         newScore,
					"sentiment_changed": originalSentiment != newSentiment,
					"score_change":      newScore - originalScore,
				})
			}
		}

		// Actualizar comentario
		comments[i].Sentiment = newSentiment
		comments[i].SentimentScore = newScore

		// Actualizar estadísticas
		stats.SentimentStats[newSentiment]++

		// Incrementar contador del rango adecuado
		scoreRange := getScoreRange(newScore)
		stats.ScoreRanges[scoreRange]++

		// Guardar en la base de datos
		if err := config.DB.Save(&comments[i]).Error; err != nil {
			stats.Failed++
			stats.FailedIds = append(stats.FailedIds, comments[i].ID)
			fmt.Printf("Error guardando comentario ID %d: %v\n", comments[i].ID, err)
			// Continuamos con el siguiente, no abortamos toda la operación
		}
	}

	// Finalizar estadísticas
	stats.EndTime = time.Now()
	stats.Duration = stats.EndTime.Sub(stats.StartTime).String()

	return stats, nil
}

// getScoreRange devuelve el rango para una puntuación dada
func getScoreRange(score float64) string {
	switch {
	case score < 2:
		return "1-2"
	case score < 3:
		return "2-3"
	case score < 4:
		return "3-4"
	case score < 5:
		return "4-5"
	case score < 6:
		return "5-6"
	case score < 7:
		return "6-7"
	case score < 8:
		return "7-8"
	case score < 9:
		return "8-9"
	default:
		return "9-10"
	}
}
