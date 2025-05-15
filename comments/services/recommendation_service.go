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
	// 1. comentarios “me gusta” (≥ 7)
	var userComments []models.Comment
	if err := config.DB.
		Where("user_id = ? AND sentiment_score >= ?", userID, 7.0).
		Preload("Movie").
		Find(&userComments).Error; err != nil {
		return nil, err
	}

	if len(userComments) == 0 {
		return getTopMovies(limit) // sin historial → ranking global
	}

	// 2. gustos
	genreLikes := map[string]int{}
	directorLikes := map[string]int{}
	commentedIDs := map[uint]struct{}{}
	for _, c := range userComments {
		genreLikes[c.Movie.Genre]++
		directorLikes[c.Movie.Director]++
		commentedIDs[c.MovieID] = struct{}{}
	}

	// 3. candidatas sin comentar
	var candidates []movieModels.Movie
	if err := config.DB.Where("id NOT IN ?", keys(commentedIDs)).Find(&candidates).Error; err != nil {
		return nil, err
	}

	// 4. puntuar
	type scored struct {
		movie  movieModels.Movie
		score  float64
		reason string
	}
	var all []scored
	for _, m := range candidates {
		_, avg, _ := GetMovieSentiment(m.ID) // función ya existente
		s := avg
		reason := ""

		if genreLikes[m.Genre] > 0 {
			s += 2
			reason = "Mismo género que sueles puntuar alto"
		}
		if directorLikes[m.Director] > 0 {
			s += 1
			if reason == "" {
				reason = "Mismo director que sueles puntuar alto"
			}
		}
		if reason == "" {
			reason = "Alta puntuación media"
		}
		all = append(all, scored{movie: m, score: s, reason: reason})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].score > all[j].score })
	if len(all) > limit {
		all = all[:limit]
	}

	// 5. respuesta
	resp := make([]RecommendationResponse, len(all))
	for i, sm := range all {
		resp[i] = RecommendationResponse{
			MovieID:         sm.movie.ID,
			Title:           sm.movie.Title,
			PredictedRating: round1(sm.score),
			RatingText:      getRatingDescription(sm.score),
			Reason:          sm.reason,
		}
	}
	return resp, nil
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
