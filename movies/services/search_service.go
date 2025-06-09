package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"fmt"
	"sort"
	"strings"
)

// SearchParams contiene los parámetros de búsqueda
type SearchParams struct {
	Title  string  `json:"title"`  // Búsqueda por título
	Genre  string  `json:"genre"`  // Filtro por nombre de género
	Rating float64 `json:"rating"` // Puntuación mínima
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

	fmt.Printf("[DEBUG-SEARCH] Iniciando búsqueda con parámetros: título=%s, género=%s, rating=%.1f\n",
		params.Title, params.Genre, params.Rating)

	// Filtro por título
	if params.Title != "" {
		// Usar LOWER para hacer la búsqueda case-insensitive de manera más compatible
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+params.Title+"%")
		fmt.Printf("[DEBUG-SEARCH] Aplicando filtro de título: LOWER(title) LIKE LOWER('%%%s%%')\n", params.Title)
	}

	// Filtro por género
	if params.Genre != "" {
		query = query.Where("LOWER(genre) LIKE LOWER(?)", "%"+params.Genre+"%")
		fmt.Printf("[DEBUG-SEARCH] Aplicando filtro de género: LOWER(genre) LIKE LOWER('%%%s%%')\n", params.Genre)
	}

	// Filtro por puntuación
	if params.Rating > 0 {
		query = query.Where("rating >= ?", params.Rating)
		fmt.Printf("[DEBUG-SEARCH] Aplicando filtro de puntuación: rating >= %.1f\n", params.Rating)
	}

	// Ejecutar la consulta sin precargar para evitar problemas
	if err := query.Find(&movies).Error; err != nil {
		fmt.Printf("[DEBUG-SEARCH] Error en la consulta: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG-SEARCH] Búsqueda completada. Encontradas %d películas.\n", len(movies))

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
	return GetUniqueGenres()
}

// GetGenreInfoList obtiene la lista de géneros con información estadística
func GetGenreInfoList() ([]GenreInfo, error) {
	genres, err := GetUniqueGenres()
	if err != nil {
		return nil, err
	}

	var genreInfoList []GenreInfo
	for _, genreName := range genres {
		stats, err := GetGenreStats(genreName)
		if err != nil {
			continue
		}
		genreInfoList = append(genreInfoList, *stats)
	}

	return genreInfoList, nil
}
