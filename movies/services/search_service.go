package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"sort"
	"strings"
)

// SearchParams contiene los parámetros de búsqueda
type SearchParams struct {
	Title     string  `json:"title"`    // Búsqueda por título
	GenreID   uint    `json:"genre_id"` // Filtro por ID de género
	GenreName string  `json:"genre"`    // Filtro por nombre de género (para compatibilidad)
	Rating    float64 `json:"rating"`   // Puntuación mínima
}

// GenreInfo contiene información sobre un género específico
type GenreInfo struct {
	Name        string  `json:"name"`         // Nombre del género
	Count       int     `json:"count"`        // Cantidad de películas
	TotalRating float64 `json:"total_rating"` // Suma de ratings para calcular promedio
	AvgRating   float64 `json:"avg_rating"`   // Rating promedio de las películas del género
}

// SearchMovies busca películas según criterios de búsqueda
func SearchMovies(params SearchParams) ([]models.Movie, error) {
	var movies []models.Movie
	query := config.DB

	// Filtro por título
	if params.Title != "" {
		query = query.Where("title ILIKE ?", "%"+params.Title+"%")
	}

	// Filtro por género (ID)
	if params.GenreID > 0 {
		query = query.Joins("JOIN movie_genres ON movies.id = movie_genres.movie_id").
			Where("movie_genres.genre_id = ?", params.GenreID)
	}

	// Filtro por género (nombre) - para compatibilidad
	if params.GenreName != "" {
		// Primero intentamos buscar por la relación
		var genre models.Genre
		if err := config.DB.Where("name ILIKE ?", "%"+params.GenreName+"%").First(&genre).Error; err == nil {
			query = query.Joins("JOIN movie_genres ON movies.id = movie_genres.movie_id").
				Where("movie_genres.genre_id = ?", genre.ID)
		} else {
			// Si no existe el género, buscamos en el campo legacy
			query = query.Where("genre ILIKE ?", "%"+params.GenreName+"%")
		}
	}

	// Filtro por puntuación
	if params.Rating > 0 {
		query = query.Where("rating >= ?", params.Rating)
	}

	// Precargar los géneros
	query = query.Preload("Genres")

	// Ejecutar la consulta
	if err := query.Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

// GetAllGenresLegacy obtiene todos los géneros disponibles con información adicional (versión antigua)
func GetAllGenresLegacy() ([]GenreInfo, error) {
	var movies []models.Movie
	if err := config.DB.Find(&movies).Error; err != nil {
		return nil, err
	}

	// Usar un mapa para acumular información de cada género
	genresMap := make(map[string]*GenreInfo)

	for _, movie := range movies {
		if movie.Genre != "" {
			// Algunos géneros pueden ser compuestos (ej: "Acción, Aventura")
			genreParts := strings.Split(movie.Genre, ",")
			for _, part := range genreParts {
				genreTrimmed := strings.TrimSpace(part)
				if genreTrimmed != "" {
					// Si el género no existe en el mapa, lo inicializamos
					if _, exists := genresMap[genreTrimmed]; !exists {
						genresMap[genreTrimmed] = &GenreInfo{
							Name:        genreTrimmed,
							Count:       0,
							TotalRating: 0,
						}
					}

					// Incrementamos el contador y acumulamos el rating
					genresMap[genreTrimmed].Count++
					genresMap[genreTrimmed].TotalRating += float64(movie.Rating)
				}
			}
		}
	}

	// Convertir el mapa a slice y calcular ratings promedio
	var genres []GenreInfo
	for _, info := range genresMap {
		if info.Count > 0 {
			info.AvgRating = info.TotalRating / float64(info.Count)
		}
		genres = append(genres, *info)
	}

	// Ordenar géneros alfabéticamente
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})

	return genres, nil
}

// GetSimpleGenres obtiene solo los nombres de los géneros (mantiene compatibilidad)
func GetSimpleGenres() ([]string, error) {
	genres, err := GetAllGenres()
	if err != nil {
		return nil, err
	}

	var genreNames []string
	for _, genre := range genres {
		genreNames = append(genreNames, genre.Name)
	}

	return genreNames, nil
}

// GetGenreInfoList obtiene la lista de géneros con información estadística
func GetGenreInfoList() ([]GenreInfo, error) {
	genres, err := GetAllGenres()
	if err != nil {
		return nil, err
	}

	var genreInfoList []GenreInfo
	for _, genre := range genres {
		stats, err := GetGenreStats(genre.ID)
		if err != nil {
			continue
		}
		genreInfoList = append(genreInfoList, *stats)
	}

	return genreInfoList, nil
}
