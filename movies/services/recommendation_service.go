package services

import (
	"cine_conecta_backend/comments/models"
	commentServices "cine_conecta_backend/comments/services"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
	"math"
	"strings"
)

// MovieWithDetails enriquece la estructura de Movie con información adicional
type MovieWithDetails struct {
	movieModels.Movie
	SentimentScore float64              `json:"sentiment_score,omitempty"`
	SentimentType  models.SentimentType `json:"sentiment_type,omitempty"`
	SentimentText  string               `json:"sentiment_text,omitempty"`
}

// EnrichMovie añade información adicional a un objeto Movie
func EnrichMovie(movie movieModels.Movie) MovieWithDetails {
	enriched := MovieWithDetails{
		Movie: movie,
	}

	// Obtener el sentimiento promedio para esta película
	sentiment, score, _ := commentServices.GetMovieSentiment(movie.ID)
	enriched.SentimentScore = score
	enriched.SentimentType = sentiment

	// Añadir descripción textual del sentimiento
	switch sentiment {
	case models.SentimentPositive:
		enriched.SentimentText = "Muy bien valorada por los usuarios"
	case models.SentimentNeutral:
		enriched.SentimentText = "Valoraciones mixtas por parte de los usuarios"
	case models.SentimentNegative:
		enriched.SentimentText = "No tan bien valorada por los usuarios"
	default:
		enriched.SentimentText = "Sin valoraciones suficientes"
	}

	return enriched
}

// GetFavoriteGenres es una versión pública de getFavoriteGenres para ser usada por el controlador
func GetFavoriteGenres(userID uint) ([]string, error) {
	return getFavoriteGenres(userID)
}

// GetRecommendedMovies devuelve películas recomendadas para un usuario basadas en el análisis de sentimientos
func GetRecommendedMovies(userID uint) ([]movieModels.Movie, error) {
	// 1. Obtener géneros de películas que al usuario le gustan (comentarios positivos)
	favoriteGenres, err := getFavoriteGenres(userID)
	if err != nil {
		return nil, err
	}

	// 2. Obtener películas que el usuario ya ha visto/comentado
	var userWatchedMovies []uint
	if err := config.DB.Model(&models.Comment{}).
		Where("user_id = ?", userID).
		Distinct("movie_id").
		Pluck("movie_id", &userWatchedMovies).Error; err != nil {
		return nil, err
	}

	// 3. Obtener todas las películas disponibles
	var allMovies []movieModels.Movie
	if err := config.DB.Find(&allMovies).Error; err != nil {
		return nil, err
	}

	// 4. Calcular puntuación de relevancia para cada película
	type RankedMovie struct {
		Movie movieModels.Movie
		Score float64
	}

	var rankedMovies []RankedMovie

	for _, movie := range allMovies {
		// Ignorar películas que el usuario ya vio
		alreadyWatched := false
		for _, watchedID := range userWatchedMovies {
			if movie.ID == watchedID {
				alreadyWatched = true
				break
			}
		}

		if alreadyWatched {
			continue
		}

		// Calcular puntuación de relevancia
		relevanceScore := 0.0

		// Factor 1: Género preferido (mayor peso)
		for _, genre := range favoriteGenres {
			if movie.Genre == genre {
				relevanceScore += 3.0
				break
			}
		}

		// Factor 2: Puntuación general de la película
		if movie.Rating > 0 {
			relevanceScore += float64(movie.Rating) * 0.2 // Rating tiene peso moderado
		}

		// Factor 3: Sentimiento colectivo sobre la película
		sentiment, avgScore, err := commentServices.GetMovieSentiment(movie.ID)
		if err == nil && avgScore > 0 {
			// Películas con sentimiento positivo reciben bonus
			if sentiment == models.SentimentPositive {
				relevanceScore += avgScore * 0.3
			}
		}

		// Solo añadir películas con una relevancia mínima
		if relevanceScore > 0 {
			rankedMovies = append(rankedMovies, RankedMovie{
				Movie: movie,
				Score: relevanceScore,
			})
		}
	}

	// 5. Ordenar por relevancia (de mayor a menor)
	for i := 0; i < len(rankedMovies)-1; i++ {
		for j := i + 1; j < len(rankedMovies); j++ {
			if rankedMovies[i].Score < rankedMovies[j].Score {
				rankedMovies[i], rankedMovies[j] = rankedMovies[j], rankedMovies[i]
			}
		}
	}

	// 6. Preparar respuesta (top 5 recomendaciones)
	var recommendations []movieModels.Movie
	for i := 0; i < len(rankedMovies) && i < 5; i++ {
		recommendations = append(recommendations, rankedMovies[i].Movie)
	}

	// 7. Si no hay suficientes recomendaciones específicas, agregar películas populares
	if len(recommendations) < 5 {
		// Evitar duplicados
		existingIDs := make(map[uint]bool)
		for _, movie := range recommendations {
			existingIDs[movie.ID] = true
		}

		// También evitar películas ya vistas
		for _, id := range userWatchedMovies {
			existingIDs[id] = true
		}

		// Obtener películas populares como respaldo
		popularMovies, err := getTopRatedMovies(10, nil)
		if err == nil {
			for _, movie := range popularMovies {
				if !existingIDs[movie.ID] {
					recommendations = append(recommendations, movie)
					existingIDs[movie.ID] = true

					if len(recommendations) >= 5 {
						break
					}
				}
			}
		}
	}

	return recommendations, nil
}

// getFavoriteGenres obtiene los géneros favoritos de un usuario basado en sus comentarios positivos
func getFavoriteGenres(userID uint) ([]string, error) {
	var comments []models.Comment
	var favoriteGenres []string

	// Obtener comentarios con puntuación alta (7 o más en escala 1-10)
	// O comentarios que contienen palabras positivas pero no tienen puntuación aún
	if err := config.DB.Where("user_id = ?", userID).
		Find(&comments).Error; err != nil {
		return favoriteGenres, err
	}

	// Aplicar análisis de sentimientos a todos los comentarios
	var positiveCommentMovieIDs []uint
	for _, comment := range comments {
		// Si el comentario ya tiene puntuación de sentimiento
		if comment.SentimentScore >= 7.0 {
			positiveCommentMovieIDs = append(positiveCommentMovieIDs, comment.MovieID)
		} else if comment.Content != "" {
			// Analizar el sentimiento del comentario si no tiene puntuación o está desactualizado
			sentiment, score := commentServices.AnalyzeSentiment(comment.Content)

			// Si es positivo, agregar a la lista
			if sentiment == models.SentimentPositive || score >= 7.0 ||
				strings.Contains(strings.ToLower(comment.Content), "me gusta") ||
				strings.Contains(strings.ToLower(comment.Content), "encant") ||
				strings.Contains(strings.ToLower(comment.Content), "recomiendo") {
				positiveCommentMovieIDs = append(positiveCommentMovieIDs, comment.MovieID)
			}
		}
	}

	// Si no hay comentarios positivos, devolver lista vacía
	if len(positiveCommentMovieIDs) == 0 {
		return favoriteGenres, nil
	}

	// Obtener géneros de estas películas
	var movies []movieModels.Movie
	if err := config.DB.Preload("Genre").Where("id IN ?", positiveCommentMovieIDs).
		Find(&movies).Error; err != nil {
		return favoriteGenres, err
	}

	// Contar frecuencia de géneros
	genreFrequency := make(map[string]int)
	for _, movie := range movies {
		if movie.Genre != "" {
			genreParts := strings.Split(movie.Genre, ",")
			for _, part := range genreParts {
				genreTrimmed := strings.TrimSpace(part)
				if genreTrimmed != "" {
					genreFrequency[genreTrimmed]++
				}
			}
		}
	}

	// Ordenar géneros por frecuencia (de mayor a menor)
	type GenreCount struct {
		Genre string
		Count int
	}

	var genreCounts []GenreCount
	for genre, count := range genreFrequency {
		genreCounts = append(genreCounts, GenreCount{Genre: genre, Count: count})
	}

	// Ordenar por frecuencia descendente
	for i := 0; i < len(genreCounts)-1; i++ {
		for j := i + 1; j < len(genreCounts); j++ {
			if genreCounts[i].Count < genreCounts[j].Count {
				genreCounts[i], genreCounts[j] = genreCounts[j], genreCounts[i]
			}
		}
	}

	// Seleccionar los géneros más frecuentes (hasta 3)
	for i := 0; i < len(genreCounts) && i < 3; i++ {
		favoriteGenres = append(favoriteGenres, genreCounts[i].Genre)
	}

	return favoriteGenres, nil
}

// getPopularMoviesByGenres obtiene películas populares de ciertos géneros
func getPopularMoviesByGenres(genres []string) ([]movieModels.Movie, error) {
	var movies []movieModels.Movie

	// Si no hay géneros, devolver lista vacía
	if len(genres) == 0 {
		return movies, nil
	}

	// Construir condiciones para la consulta
	var conditions []string
	var values []interface{}

	for _, genreName := range genres {
		conditions = append(conditions, "genre ILIKE ?")
		values = append(values, "%"+genreName+"%")
	}

	// Obtener películas con los géneros especificados y rating alto
	query := config.DB.Order("rating DESC")

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	if err := query.Limit(10).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

// getTopRatedMovies obtiene las películas mejor valoradas
func getTopRatedMovies(limit int, excludeIDs []uint) ([]movieModels.Movie, error) {
	var movies []movieModels.Movie
	query := config.DB.Order("rating DESC").Limit(limit)

	// Excluir películas ya vistas
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	if err := query.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

// GetMoviesByPositiveSentiment obtiene películas con mayor porcentaje de comentarios positivos
func GetMoviesByPositiveSentiment(limit int) ([]movieModels.Movie, error) {
	// Obtener todas las películas
	var allMovies []movieModels.Movie
	if err := config.DB.Find(&allMovies).Error; err != nil {
		return nil, err
	}

	// Calcular ratio de sentimiento para cada película
	type MovieScore struct {
		Movie      movieModels.Movie
		Score      float64
		CommentCnt int
	}

	var scoredMovies []MovieScore

	for _, movie := range allMovies {
		sentiment, score, err := commentServices.GetMovieSentiment(movie.ID)
		if err != nil {
			continue
		}

		// Sólo agregar películas con al menos un comentario
		var commentCount int64
		config.DB.Model(&models.Comment{}).Where("movie_id = ?", movie.ID).Count(&commentCount)

		if commentCount > 0 {
			// Dar más peso a películas con más comentarios pero manteniendo
			// la importancia de la puntuación de sentimiento
			adjustedScore := score

			// Si hay pocos comentarios (1-2), reducir ligeramente la puntuación
			// Si hay muchos comentarios (5+), aumentar ligeramente la puntuación
			if commentCount < 3 {
				adjustedScore *= 0.9 // Reducir un 10%
			} else if commentCount >= 5 {
				// Incrementar hasta un 20% para películas con muchos comentarios
				bonus := 1.0 + (math.Min(float64(commentCount), 20.0)-5.0)/75.0
				adjustedScore *= bonus
			}

			// Películas con sentimiento positivo reciben un bonus adicional
			if sentiment == models.SentimentPositive {
				adjustedScore *= 1.1 // 10% adicional
			}

			scoredMovies = append(scoredMovies, MovieScore{
				Movie:      movie,
				Score:      adjustedScore,
				CommentCnt: int(commentCount),
			})
		}
	}

	// Ordenar películas por puntuación de sentimiento ajustada (de mayor a menor)
	for i := 0; i < len(scoredMovies)-1; i++ {
		for j := i + 1; j < len(scoredMovies); j++ {
			if scoredMovies[i].Score < scoredMovies[j].Score {
				scoredMovies[i], scoredMovies[j] = scoredMovies[j], scoredMovies[i]
			}
		}
	}

	// Preparar resultado
	var result []movieModels.Movie
	for i := 0; i < len(scoredMovies) && i < limit; i++ {
		result = append(result, scoredMovies[i].Movie)
	}

	return result, nil
}
