package services

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"

	openai "github.com/sashabaranov/go-openai"
)

// SentimentAnalysisResult estructura para parsear el resultado del análisis
type SentimentAnalysisResult struct {
	Sentiment models.SentimentType `json:"sentiment"`
	Score     float64              `json:"score"`
	Reason    string               `json:"reason"`
}

// AnalyzeSentimentWithML analiza sentimientos usando OpenAI
func AnalyzeSentimentWithML(content string) (models.SentimentType, float64, error) {
	// Usar el método léxico tradicional como fallback
	if os.Getenv("OPENAI_API_KEY") == "" || len(content) < 10 {
		sentimentType, score := AnalyzeSentiment(content)
		return sentimentType, score, nil
	}

	// Crear cliente OpenAI
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Preparar sistema y mensaje del usuario
	systemMessage := "Eres un sistema especializado en análisis de sentimientos para comentarios de películas en español."
	userMessage := `Analiza el siguiente comentario y clasifícalo como 'positive', 'neutral' o 'negative'. 
Además, asigna una puntuación de 1 a 10 (donde 1 es extremadamente negativo y 10 es extremadamente positivo).
Responde en formato JSON con los campos 'sentiment', 'score' y 'reason'.

Comentario: ` + content

	// Llamar a la API
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			MaxTokens:   150,
			Temperature: 0.0, // Queremos respuestas consistentes
		},
	)

	if err != nil {
		// Fallback al método léxico tradicional en caso de error
		sentimentType, score := AnalyzeSentiment(content)
		return sentimentType, score, nil
	}

	// Extraer el JSON de la respuesta
	jsonContent := resp.Choices[0].Message.Content
	jsonContent = strings.TrimSpace(jsonContent)

	// Si la respuesta no comienza con '{', buscar el primer '{'
	if !strings.HasPrefix(jsonContent, "{") {
		startBrace := strings.Index(jsonContent, "{")
		if startBrace != -1 {
			jsonContent = jsonContent[startBrace:]
		}
	}

	// Si la respuesta no termina con '}', buscar el último '}'
	if !strings.HasSuffix(jsonContent, "}") {
		endBrace := strings.LastIndex(jsonContent, "}")
		if endBrace != -1 {
			jsonContent = jsonContent[:endBrace+1]
		}
	}

	// Parsear el resultado
	var result SentimentAnalysisResult
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		// Fallback al método léxico tradicional en caso de error
		sentimentType, score := AnalyzeSentiment(content)
		return sentimentType, score, nil
	}

	return result.Sentiment, result.Score, nil
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
