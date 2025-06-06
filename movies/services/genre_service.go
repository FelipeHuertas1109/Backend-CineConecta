package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"sort"
	"strings"
)

// CreateGenre crea un nuevo género si no existe
func CreateGenre(name string) (*models.Genre, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, nil
	}

	var genre models.Genre
	// Verificar si el género ya existe
	result := config.DB.Where("name = ?", name).First(&genre)
	if result.Error == nil {
		// El género ya existe
		return &genre, nil
	}

	// Crear nuevo género
	genre = models.Genre{Name: name}
	if err := config.DB.Create(&genre).Error; err != nil {
		return nil, err
	}

	return &genre, nil
}

// GetAllGenres obtiene todos los géneros disponibles
func GetAllGenres() ([]models.Genre, error) {
	var genres []models.Genre
	if err := config.DB.Find(&genres).Error; err != nil {
		return nil, err
	}

	// Ordenar géneros alfabéticamente
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})

	return genres, nil
}

// GetGenreByID busca un género por su ID
func GetGenreByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	if err := config.DB.First(&genre, id).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

// GetGenreByName busca un género por su nombre
func GetGenreByName(name string) (*models.Genre, error) {
	var genre models.Genre
	if err := config.DB.Where("name = ?", name).First(&genre).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

// GetGenresForMovie obtiene todos los géneros de una película
func GetGenresForMovie(movieID uint) ([]models.Genre, error) {
	var movie models.Movie
	if err := config.DB.Preload("Genres").First(&movie, movieID).Error; err != nil {
		return nil, err
	}
	return movie.Genres, nil
}

// AddGenreToMovie añade un género a una película
func AddGenreToMovie(movieID uint, genreID uint) error {
	var movie models.Movie
	var genre models.Genre

	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	if err := config.DB.First(&genre, genreID).Error; err != nil {
		return err
	}

	// Verificar si la relación ya existe
	var count int64
	config.DB.Table("movie_genres").
		Where("movie_id = ? AND genre_id = ?", movieID, genreID).
		Count(&count)

	if count == 0 {
		return config.DB.Model(&movie).Association("Genres").Append(&genre)
	}

	return nil
}

// RemoveGenreFromMovie elimina un género de una película
func RemoveGenreFromMovie(movieID uint, genreID uint) error {
	var movie models.Movie
	var genre models.Genre

	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	if err := config.DB.First(&genre, genreID).Error; err != nil {
		return err
	}

	return config.DB.Model(&movie).Association("Genres").Delete(&genre)
}

// GetMoviesByGenre obtiene todas las películas de un género específico
func GetMoviesByGenre(genreID uint) ([]models.Movie, error) {
	var genre models.Genre
	if err := config.DB.Preload("Movies").First(&genre, genreID).Error; err != nil {
		return nil, err
	}
	return genre.Movies, nil
}

// UpdateMovieGenres actualiza los géneros de una película
func UpdateMovieGenres(movieID uint, genreIDs []uint) error {
	var movie models.Movie
	if err := config.DB.First(&movie, movieID).Error; err != nil {
		return err
	}

	var genres []models.Genre
	if err := config.DB.Find(&genres, genreIDs).Error; err != nil {
		return err
	}

	return config.DB.Model(&movie).Association("Genres").Replace(&genres)
}

// GetGenreStats obtiene estadísticas de un género
func GetGenreStats(genreID uint) (*GenreInfo, error) {
	var genre models.Genre
	if err := config.DB.Preload("Movies").First(&genre, genreID).Error; err != nil {
		return nil, err
	}

	info := &GenreInfo{
		Name:        genre.Name,
		Count:       len(genre.Movies),
		TotalRating: 0,
		AvgRating:   0,
	}

	if info.Count > 0 {
		for _, movie := range genre.Movies {
			info.TotalRating += float64(movie.Rating)
		}
		info.AvgRating = info.TotalRating / float64(info.Count)
	}

	return info, nil
}
