// comments/services/sentiment_service.go
package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
)

/*------------------------------------------------------------
 1. Listas de palabras (fallback heurístico)
------------------------------------------------------------*/

// positivas: +≈60
var positiveWords = []string{
	"bueno", "excelente", "increíble", "fantástico", "maravilloso", "genial",
	"me gusta", "recomiendo", "maravilla", "espectacular", "encantad", "mejor",
	"perfecto", "fenomenal", "magnífico", "impresionante", "positivo",
	"brillante", "obra maestra", "imperdible", "inolvidable", "conmovedor",
	"fascinante", "grandioso", "excepcional", "hermoso", "innovador",
	"impactante", "extraordinario", "sublime", "radiante", "divertida",
	"maravillosa", "emocionante", "apasionante", "deslumbrante", "redonda",
	"formidable", "colosal", "estupenda", "espléndida", "alucinante",
	"notable", "entretenida", "vibrante", "inteligente", "audaz",
	"cautivadora", "sorprendente", "refrescante", "fantástica", "adictiva",
	"épica", "gloriosa", "sensacional",
}

// negativas: +≈60
var negativeWords = []string{
	"malo", "terrible", "horrible", "malísimo", "pésimo", "aburrido",
	"no me gusta", "no recomiendo", "decepcionante", "decepción", "desagrada",
	"peor", "fatal", "desastre", "mediocre", "negativo", "triste", "odio",
	"desagradable", "molest", "fastidios", "insoportable", "inútil",
	"innecesario", "sin sentido", "predecible", "lento", "desperdicio",
	"olvidable", "ridículo", "vergonzoso", "pretencioso", "confuso", "asco",
	"insufrible", "detestable", "espantoso", "espantosa", "lamentable",
	"patético", "chapucero", "floja", "flojo", "innecesaria", "torpe",
	"torpe", "cutre", "pobre", "insulsa", "burda", "desastrosa", "espesa",
	"incoherente", "plana", "forzada", "fría", "superficial",
}

// intensificadores
var emphasisWords = []string{
	"muy", "realmente", "extremadamente", "absolutamente", "totalmente",
	"completamente", "verdaderamente", "genuinamente", "extraordinariamente",
	"sumamente", "increíblemente", "tremendamente", "excepcionalmente",
}

// negativas "fuertes" (peso doble)
var strongNeg = []string{
	"basura", "desastre", "asco", "infame", "detestable", "nefasta",
	"espantoso", "espantosa", "abominable", "horripilante",
}

/*
------------------------------------------------------------
 2. API pública

------------------------------------------------------------
*/
func AnalyzeSentiment(content string) (models.SentimentType, float64) {
	debug := os.Getenv("SENTIMENT_DEBUG") == "true"
	if debug {
		fmt.Println("↓ Analyze:", content)
		fmt.Println("⚠️  Usando análisis heurístico")
	}

	// Usar directamente el análisis heurístico
	return analyzeHeuristic(content)
}

/*
------------------------------------------------------------
 3. Heurístico local

------------------------------------------------------------
*/
func analyzeHeuristic(content string) (models.SentimentType, float64) {
	low := strings.ToLower(content)
	reg := regexp.MustCompile("[^a-záéíóúüñ\\s]")
	low = reg.ReplaceAllString(low, "")
	words := strings.Fields(low)
	if len(words) == 0 {
		return models.SentimentNeutral, 3.0
	}

	pos, neg := 0, 0
	for i, w := range words {
		// multiplicador si hay intensificador previo
		mult := 1.0
		if i > 0 {
			for _, e := range emphasisWords {
				if strings.Contains(words[i-1], e) {
					mult = 1.5
					break
				}
			}
		}
		// positivas
		for _, pw := range positiveWords {
			if strings.Contains(w, pw) {
				pos += int(mult)
				break
			}
		}
		// negativas (con bonus strongNeg)
		for _, nw := range negativeWords {
			if strings.Contains(w, nw) {
				neg += int(mult)
				for _, sn := range strongNeg {
					if strings.Contains(w, sn) {
						neg += int(mult) // doble peso
						break
					}
				}
				break
			}
		}
	}

	score := (float64(pos) - float64(neg)) / float64(len(words)) // -1..1
	rating := ((score + 1) / 2 * 4) + 1                          // 1..5
	if len(words) < 5 {
		rating = rating*0.7 + 3*0.3
	}
	rating = math.Round(math.Max(1, math.Min(5, rating))*10) / 10
	return scoreToType(rating), rating
}

// scoreToType convierte un score a un tipo de sentimiento
func scoreToType(r float64) models.SentimentType {
	switch {
	case r >= 4:
		return models.SentimentPositive
	case r <= 2:
		return models.SentimentNegative
	default:
		return models.SentimentNeutral
	}
}

/*
------------------------------------------------------------
 4. Agrupados por película

------------------------------------------------------------
*/
func GetMovieSentiment(movieID uint) (models.SentimentType, float64, error) {
	var comments []models.Comment
	if err := config.DB.Where("movie_id = ?", movieID).Find(&comments).Error; err != nil {
		return models.SentimentNeutral, 3.0, err
	}
	if len(comments) == 0 {
		return models.SentimentNeutral, 3.0, nil
	}
	var sum float64
	for _, c := range comments {
		sum += c.SentimentScore
	}
	avg := sum / float64(len(comments))
	return scoreToType(avg), avg, nil
}

// GetMovieSentimentByName obtiene el sentimiento promedio para una película por su nombre
func GetMovieSentimentByName(movieName string) (models.SentimentType, float64, error) {
	// Primero buscamos la película por nombre
	var movieIDs []uint
	err := config.DB.Model(&movieModels.Movie{}).
		Where("LOWER(title) LIKE LOWER(?)", "%"+movieName+"%").
		Pluck("id", &movieIDs).Error

	if err != nil {
		return models.SentimentNeutral, 3.0, err
	}

	if len(movieIDs) == 0 {
		return models.SentimentNeutral, 3.0, fmt.Errorf("no se encontraron películas con el nombre: %s", movieName)
	}

	// Si encontramos varias películas, usamos la primera que coincida exactamente o la primera en la lista
	var selectedMovieID uint
	var exactMatch bool

	// Intentar encontrar una coincidencia exacta
	var movies []movieModels.Movie
	if err := config.DB.Where("id IN ?", movieIDs).Find(&movies).Error; err != nil {
		return models.SentimentNeutral, 3.0, err
	}

	for _, movie := range movies {
		if strings.EqualFold(movie.Title, movieName) {
			selectedMovieID = movie.ID
			exactMatch = true
			break
		}
	}

	// Si no hay coincidencia exacta, usar el primer resultado
	if !exactMatch && len(movieIDs) > 0 {
		selectedMovieID = movieIDs[0]
	}

	// Ahora obtenemos el sentimiento para la película seleccionada
	return GetMovieSentiment(selectedMovieID)
}
