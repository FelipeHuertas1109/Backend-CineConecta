package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
)

const cineConectaMLURL = "https://cine-conecta-ml.onrender.com/api/score/"

type ScoreResponse struct {
	Score float64 `json:"score"`
}

// SentimentAnalysisResult estructura para parsear el resultado del análisis
type SentimentAnalysisResult struct {
	Sentiment models.SentimentType `json:"sentiment"`
	Score     float64              `json:"score"`
	Reason    string               `json:"reason"`
}

// AnalyzeSentimentWithML devuelve un rating usando la API de Cine Conecta ML
func AnalyzeSentimentWithML(text string) (models.SentimentType, float64, error) {
	fmt.Println("[DEBUG-ML] Iniciando análisis de sentimiento con API externa")

	// Preparar el cuerpo de la solicitud
	requestBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		fmt.Printf("[DEBUG-ML] Error al preparar solicitud JSON: %v\n", err)
		return models.SentimentNeutral, 3.0, fmt.Errorf("error al preparar solicitud: %w", err)
	}

	// Crear la solicitud HTTP
	fmt.Printf("[DEBUG-ML] Enviando solicitud a: %s\n", cineConectaMLURL)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", cineConectaMLURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("[DEBUG-ML] Error al crear solicitud HTTP: %v\n", err)
		return models.SentimentNeutral, 3.0, fmt.Errorf("error al crear solicitud: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar la solicitud
	fmt.Println("[DEBUG-ML] Ejecutando solicitud HTTP...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[DEBUG-ML] Error en la solicitud HTTP: %v\n", err)
		return models.SentimentNeutral, 3.0, err
	}
	defer resp.Body.Close()

	fmt.Printf("[DEBUG-ML] Respuesta recibida: HTTP %d %s\n", resp.StatusCode, resp.Status)

	// Leer el cuerpo de la respuesta para depuración
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("[DEBUG-ML] Cuerpo de la respuesta: %s\n", string(respBody))

	// Crear un nuevo reader con el cuerpo leído para decodificarlo después
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[DEBUG-ML] Error en la API: código %d\n", resp.StatusCode)
		return models.SentimentNeutral, 3.0, fmt.Errorf("API respondió con código %d", resp.StatusCode)
	}

	// Decodificar la respuesta
	var scoreResp ScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&scoreResp); err != nil {
		fmt.Printf("[DEBUG-ML] Error al decodificar respuesta JSON: %v\n", err)
		return models.SentimentNeutral, 3.0, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	fmt.Printf("[DEBUG-ML] Score obtenido: %.2f\n", scoreResp.Score)

	// Convertir el score a un tipo de sentimiento usando la función definida en sentiment_service.go
	sentimentType := scoreToType(scoreResp.Score)
	fmt.Printf("[DEBUG-ML] Tipo de sentimiento determinado: %s\n", sentimentType)

	return sentimentType, scoreResp.Score, nil
}

// GetMovieSentimentWithML obtiene el sentimiento promedio para una película usando ML cuando es posible
func GetMovieSentimentWithML(movieID uint) (models.SentimentType, float64, error) {
	var comments []models.Comment

	if err := config.DB.Where("movie_id = ?", movieID).Find(&comments).Error; err != nil {
		return models.SentimentNeutral, 5.0, err
	}

	if len(comments) == 0 {
		return models.SentimentNeutral, 5.0, nil
	}

	var totalScore float64
	for _, comment := range comments {
		totalScore += comment.SentimentScore
	}

	avgScore := totalScore / float64(len(comments))

	var sentiment models.SentimentType
	if avgScore >= 7 {
		sentiment = models.SentimentPositive
	} else if avgScore <= 4 {
		sentiment = models.SentimentNegative
	} else {
		sentiment = models.SentimentNeutral
	}

	return sentiment, avgScore, nil
}
