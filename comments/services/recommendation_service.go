package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"math"
	"sort"
)

// ---------- tipos ----------

type RecommendationResponse struct {
	MovieID         uint    `json:"movie_id"`
	Title           string  `json:"title"`
	PredictedRating float64 `json:"predicted_rating"`
	RatingText      string  `json:"rating_text"`
	Reason          string  `json:"reason"`
}

// ---------- API pública ----------

// GetRecommendationsForUser genera recomendaciones basadas
// en los comentarios (sentiment_score 1-10) del usuario.
func GetRecommendationsForUser(userID uint, limit int) ([]RecommendationResponse, error) {
	// Obtener todos los comentarios del usuario (independientemente del score)
	var userComments []models.Comment
	if err := config.DB.
		Where("user_id = ?", userID).
		Preload("Movie").
		Find(&userComments).Error; err != nil {
		return nil, err
	}

	// Si no hay comentarios, retornar lista vacía
	if len(userComments) == 0 {
		return []RecommendationResponse{}, nil
	}

	// Preparar películas comentadas
	commentedMovies := make(map[uint]*movieModels.Movie)
	for _, c := range userComments {
		commentedMovies[c.MovieID] = &c.Movie
	}

	// Preparar recomendaciones solo de películas que el usuario ha comentado
	var recommendations []RecommendationResponse
	for movieID, movie := range commentedMovies {
		// Obtener sentimiento asociado a esta película para este usuario
		var userSentiment float64
		var userCommentForThisMovie models.Comment
		if err := config.DB.
			Where("user_id = ? AND movie_id = ?", userID, movieID).
			Order("updated_at DESC"). // Comentario más reciente
			First(&userCommentForThisMovie).Error; err == nil {
			userSentiment = userCommentForThisMovie.SentimentScore
		}

		// Convertir sentimiento a puntuación (0-10)
		score := userSentiment
		if score == 0 {
			score = 5 // Valor neutral si no hay sentimiento detectado
		}

		// Determinar razón de la recomendación
		reason := "Película que has comentado"
		if score >= 7 {
			reason = "Película que has valorado positivamente"
		} else if score < 5 {
			reason = "Película que has valorado negativamente"
		}

		// Añadir a recomendaciones
		recommendations = append(recommendations, RecommendationResponse{
			MovieID:         movieID,
			Title:           movie.Title,
			PredictedRating: round1(score),
			RatingText:      getRatingDescription(score),
			Reason:          reason,
		})
	}

	// Ordenar por puntuación (descendente)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].PredictedRating > recommendations[j].PredictedRating
	})

	// Limitar resultados si es necesario
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// ---------- helpers ----------

func keys(m map[uint]struct{}) []uint {
	out := make([]uint, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
func round1(v float64) float64 { return math.Round(v*10) / 10 }

func getRatingDescription(score float64) string {
	switch {
	case score >= 9.5:
		return "Obra maestra"
	case score >= 9.0:
		return "Excepcional"
	case score >= 8.0:
		return "Excelente"
	case score >= 7.0:
		return "Muy buena"
	case score >= 6.0:
		return "Buena"
	case score >= 5.0:
		return "Aceptable"
	case score >= 4.0:
		return "Regular"
	case score >= 3.0:
		return "Mala"
	case score >= 2.0:
		return "Muy mala"
	default:
		return "Pésima"
	}
}

// ranking global para usuarios sin historial
func getTopMovies(limit int) ([]RecommendationResponse, error) {
	var movies []movieModels.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	type scored struct {
		movie movieModels.Movie
		score float64
	}
	var tmp []scored
	for _, m := range movies {
		_, avg, err := GetMovieSentiment(m.ID)
		if err != nil {
			return nil, err
		}
		tmp = append(tmp, scored{movie: m, score: avg})
	}
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].score > tmp[j].score })
	if len(tmp) > limit {
		tmp = tmp[:limit]
	}

	resp := make([]RecommendationResponse, len(tmp))
	for i, sm := range tmp {
		resp[i] = RecommendationResponse{
			MovieID:         sm.movie.ID,
			Title:           sm.movie.Title,
			PredictedRating: round1(sm.score),
			RatingText:      getRatingDescription(sm.score),
			Reason:          "Alta puntuación media",
		}
	}
	return resp, nil
}
