package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	"math"
	"regexp"
	"strings"
)

// Lista de palabras positivas y negativas para un análisis básico
var (
	positiveWords = []string{
		"bueno", "excelente", "increíble", "fantástico", "maravilloso", "genial",
		"me gusta", "recomiendo", "maravilla", "espectacular", "encantado", "mejor",
		"buena", "perfecto", "fenomenal", "magnífico", "impresionante", "positivo",
		"feliz", "alegre", "agradable", "divertido", "entretenido", "bien", "amor",
		"brillante", "obra maestra", "imperdible", "inolvidable", "conmovedor",
		"fascinante", "grandioso", "excepcional", "hermoso", "innovador",
	}

	negativeWords = []string{
		"malo", "terrible", "horrible", "malísimo", "pésimo", "aburrido",
		"no me gusta", "no recomiendo", "decepcionante", "decepción", "desagrada", "peor",
		"mala", "fatal", "desastre", "mediocre", "negativo", "triste", "odio", "mal",
		"desagradable", "molesto", "fastidioso", "insoportable", "inútil", "innecesario",
		"pobre", "sin sentido", "predecible", "lento", "desperdicio", "olvidable",
		"horrible", "ridículo", "vergonzoso", "pretencioso", "confuso",
	}

	// Palabras con peso adicional (enfáticas)
	emphasisWords = []string{
		"muy", "realmente", "extremadamente", "absolutamente", "totalmente",
		"completamente", "verdaderamente", "genuinamente", "extraordinariamente",
		"sumamente", "increíblemente", "tremendamente", "excepcionalmente",
	}
)

// AnalyzeSentiment analiza el sentimiento de un comentario y devuelve tipo de sentimiento y puntuación en escala 1-10
func AnalyzeSentiment(content string) (models.SentimentType, float64) {
	// Convertir a minúsculas para facilitar la comparación
	normalizedContent := strings.ToLower(content)

	// Eliminar caracteres especiales excepto espacios
	reg := regexp.MustCompile("[^a-záéíóúüñ\\s]")
	normalizedContent = reg.ReplaceAllString(normalizedContent, "")

	// Contar palabras positivas, negativas y enfáticas
	positiveCount := 0
	negativeCount := 0
	emphasisCount := 0

	// Palabras en el comentario
	words := strings.Fields(normalizedContent)

	// Verificar si hay palabras enfáticas cerca de las palabras positivas/negativas
	for i, word := range words {
		// Buscar palabras de énfasis
		for _, emphWord := range emphasisWords {
			if strings.Contains(word, emphWord) {
				emphasisCount++
				break
			}
		}

		// Buscar palabras positivas
		for _, positiveWord := range positiveWords {
			if strings.Contains(word, positiveWord) {
				// Dar más peso si hay una palabra enfática justo antes
				multiplier := 1.0
				if i > 0 {
					for _, emphWord := range emphasisWords {
						if strings.Contains(words[i-1], emphWord) {
							multiplier = 1.5
							break
						}
					}
				}
				positiveCount += int(1.0 * multiplier)
				break
			}
		}

		// Buscar palabras negativas
		for _, negativeWord := range negativeWords {
			if strings.Contains(word, negativeWord) {
				// Dar más peso si hay una palabra enfática justo antes
				multiplier := 1.0
				if i > 0 {
					for _, emphWord := range emphasisWords {
						if strings.Contains(words[i-1], emphWord) {
							multiplier = 1.5
							break
						}
					}
				}
				negativeCount += int(1.0 * multiplier)
				break
			}
		}
	}

	// Calcular puntaje de sentimiento: rango de -1 (muy negativo) a 1 (muy positivo)
	totalWords := len(words)
	if totalWords == 0 {
		return models.SentimentNeutral, 5.0 // Puntuación neutral en escala 1-10
	}

	// Calcular puntuación relativa a la longitud del texto
	positiveScore := float64(positiveCount) / float64(totalWords)
	negativeScore := float64(negativeCount) / float64(totalWords)

	// Puntuación final entre -1 y 1
	finalScore := positiveScore - negativeScore

	// Convertir el rango de -1 a 1 a una escala de 1 a 10
	// Donde -1 = 1, 0 = 5.5, y 1 = 10
	rating := ((finalScore+1)/2)*9 + 1

	// Ajustar la puntuación considerando el conteo de palabras
	// (comentarios cortos podrían dar resultados extremos)
	if totalWords < 5 {
		// Para comentarios muy cortos, acercamos un poco la puntuación hacia el punto neutral (5.5)
		rating = rating*0.7 + 5.5*0.3
	}

	// Limitar a un rango de 1 a 10
	rating = math.Max(1, math.Min(10, rating))
	rating = math.Round(rating*10) / 10 // Redondear a un decimal

	// Determinar el tipo de sentimiento
	var sentiment models.SentimentType
	if rating >= 7 {
		sentiment = models.SentimentPositive
	} else if rating <= 4 {
		sentiment = models.SentimentNegative
	} else {
		sentiment = models.SentimentNeutral
	}

	return sentiment, rating
}

// GetMovieSentiment obtiene el sentimiento promedio para una película en escala 1-10
func GetMovieSentiment(movieID uint) (models.SentimentType, float64, error) {
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
