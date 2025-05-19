package services

import (
	"cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	movieModels "cine_conecta_backend/movies/models"
)

// Estructura para almacenar el resultado con la película y su puntuación media
type MovieWithAverageScore struct {
	Movie        movieModels.Movie `json:"movie"`
	AverageScore float64           `json:"average_score"`
	CommentCount int               `json:"comment_count"`
}

// GetTopRatedMovies obtiene las 5 películas mejor valoradas según la puntuación media de comentarios
func GetTopRatedMovies() ([]MovieWithAverageScore, error) {
	// Consulta SQL para obtener películas con su puntuación media de comentarios
	// ordenadas por puntuación media descendente y limitadas a 5
	query := `
		SELECT 
			m.id, m.title, m.description, m.genre, m.director, 
			m.release_date, m.rating, m.poster_url, m.created_at, m.updated_at,
			AVG(c.sentiment_score) as average_score,
			COUNT(c.id) as comment_count
		FROM 
			movies m
		JOIN 
			comments c ON m.id = c.movie_id
		GROUP BY 
			m.id
		HAVING 
			COUNT(c.id) >= 3 -- Mínimo 3 comentarios para considerar una película
		ORDER BY 
			average_score DESC
		LIMIT 5
	`

	rows, err := config.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MovieWithAverageScore
	for rows.Next() {
		var movie movieModels.Movie
		var avgScore float64
		var commentCount int

		err := rows.Scan(
			&movie.ID, &movie.Title, &movie.Description, &movie.Genre, &movie.Director,
			&movie.ReleaseDate, &movie.Rating, &movie.PosterURL, &movie.CreatedAt, &movie.UpdatedAt,
			&avgScore, &commentCount,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, MovieWithAverageScore{
			Movie:        movie,
			AverageScore: avgScore,
			CommentCount: commentCount,
		})
	}

	return result, nil
}

// GetMovieWithAverageScore obtiene una película con su puntuación media de comentarios
func GetMovieWithAverageScore(movieID uint) (MovieWithAverageScore, error) {
	var result MovieWithAverageScore

	// Obtener la película
	if err := config.DB.First(&result.Movie, movieID).Error; err != nil {
		return result, err
	}

	// Obtener comentarios de la película
	var comments []models.Comment
	if err := config.DB.Where("movie_id = ?", movieID).Find(&comments).Error; err != nil {
		return result, err
	}

	// Calcular puntuación media
	if len(comments) > 0 {
		var totalScore float64
		for _, comment := range comments {
			totalScore += comment.SentimentScore
		}
		result.AverageScore = totalScore / float64(len(comments))
		result.CommentCount = len(comments)
	}

	return result, nil
}
