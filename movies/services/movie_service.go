package services

import (
	commentModels "cine_conecta_backend/comments/models"
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// CreateMovie guarda una nueva película en la base de datos.
func CreateMovie(movie *models.Movie) error {
	return config.DB.Create(movie).Error
}

// GetMovies obtiene todas las películas y actualiza sus ratings basados en los comentarios.
func GetMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := config.DB.Preload("Genres").Find(&movies).Error; err != nil {
		return nil, err
	}

	// Actualizar los ratings de todas las películas basado en los comentarios
	for i := range movies {
		if err := updateMovieRating(&movies[i]); err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al actualizar rating de película %d: %v\n", movies[i].ID, err)
		}
	}

	return movies, nil
}

// updateMovieRating actualiza el rating de una película basado en los comentarios
func updateMovieRating(movie *models.Movie) error {
	// Obtener todos los comentarios de la película
	var comments []commentModels.Comment
	if err := config.DB.Where("movie_id = ?", movie.ID).Find(&comments).Error; err != nil {
		return fmt.Errorf("error al obtener comentarios: %w", err)
	}

	// Calcular el nuevo rating basado en el promedio de las puntuaciones de sentimientos
	var newRating float32
	if len(comments) > 0 {
		var totalScore float64
		for _, comment := range comments {
			totalScore += comment.SentimentScore
		}
		newRating = float32(totalScore / float64(len(comments)))
	} else {
		// Si no hay comentarios, dejar el rating en 0
		newRating = 0
	}

	// Solo actualizar si el rating ha cambiado
	if movie.Rating != newRating {
		movie.Rating = newRating
		// Guardar el nuevo rating en la base de datos
		if err := config.DB.Model(&models.Movie{}).Where("id = ?", movie.ID).Update("rating", newRating).Error; err != nil {
			return fmt.Errorf("error al actualizar rating: %w", err)
		}
	}

	return nil
}

// GetRecentMovies obtiene las películas más recientes según su fecha de lanzamiento.
func GetRecentMovies(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Order("release_date DESC").Limit(limit).Find(&movies).Error
	if err != nil {
		return nil, err
	}

	// Actualizar los ratings de las películas
	for i := range movies {
		if err := updateMovieRating(&movies[i]); err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al actualizar rating de película %d: %v\n", movies[i].ID, err)
		}
	}

	return movies, nil
}

// GetMovieByID obtiene una película por su ID.
func GetMovieByID(id uint) (models.Movie, error) {
	var movie models.Movie
	if err := config.DB.Preload("Genres").First(&movie, id).Error; err != nil {
		return movie, err
	}

	// Actualizar el rating de la película
	if err := updateMovieRating(&movie); err != nil {
		fmt.Printf("[DEBUG-SERVICE] Error al actualizar rating de película %d: %v\n", movie.ID, err)
	}

	return movie, nil
}

// UpdateMovie actualiza una película existente.
func UpdateMovie(movie *models.Movie) error {
	return config.DB.Save(movie).Error
}

// DeleteMovie elimina una película por su ID.
func DeleteMovie(id uint) error {
	return config.DB.Delete(&models.Movie{}, id).Error
}

// CreateMovieWithGenres crea una película con su género asociado
func CreateMovieWithGenres(movie *models.Movie, genreNames []string) error {
	// Iniciar transacción
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Si no hay géneros explícitos pero hay un género en el campo texto, usarlo
	if len(genreNames) == 0 && movie.Genre != "" {
		genreNames = models.ParseGenresString(movie.Genre)
	}

	// Guardar el género como texto también (para facilidad de búsqueda)
	if len(genreNames) > 0 {
		movie.Genre = strings.Join(genreNames, ", ")
	}

	// Crear la película primero
	if err := tx.Create(movie).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Procesar géneros
	if len(genreNames) > 0 {
		for _, genreName := range genreNames {
			// Buscar o crear género
			genre, err := findOrCreateGenre(tx, genreName)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Asociar género con película en la tabla movie_genres
			if err := tx.Exec("INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				movie.ID, genre.ID).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit de la transacción
	return tx.Commit().Error
}

// UpdateMovieWithGenres actualiza una película y sus géneros
func UpdateMovieWithGenres(movie *models.Movie, genreNames []string) error {
	// Iniciar transacción
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Si no hay géneros explícitos pero hay un género en el campo texto, usarlo
	if len(genreNames) == 0 && movie.Genre != "" {
		genreNames = models.ParseGenresString(movie.Genre)
	}

	// Guardar el género como texto también (para facilidad de búsqueda)
	if len(genreNames) > 0 {
		movie.Genre = strings.Join(genreNames, ", ")
	} else {
		movie.Genre = ""
	}

	// Actualizar la película
	if err := tx.Save(movie).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Limpiar géneros existentes
	if err := tx.Exec("DELETE FROM movie_genres WHERE movie_id = ?", movie.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Procesar géneros
	if len(genreNames) > 0 {
		for _, genreName := range genreNames {
			// Buscar o crear género
			genre, err := findOrCreateGenre(tx, genreName)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Asociar género con película en la tabla movie_genres
			if err := tx.Exec("INSERT INTO movie_genres (movie_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				movie.ID, genre.ID).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit de la transacción
	return tx.Commit().Error
}

// findOrCreateGenre busca un género por nombre o lo crea si no existe
func findOrCreateGenre(tx *gorm.DB, name string) (*models.Genre, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, nil
	}

	var genre models.Genre

	// Buscar por nombre
	if err := tx.Where("name = ?", name).First(&genre).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Crear si no existe
			genre = models.Genre{Name: name}
			if err := tx.Create(&genre).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &genre, nil
}

var allowedSortFields = map[string]bool{
	"title":  true,
	"rating": true,
}

func GetMoviesSorted(sortBy string, order string) ([]models.Movie, error) {
	var movies []models.Movie

	// Validar columna
	if !allowedSortFields[sortBy] {
		return nil, errors.New("columna inválida para ordenar")
	}

	// Validar orden
	if order != "asc" && order != "desc" {
		return nil, errors.New("orden inválido")
	}

	sortQuery := sortBy + " " + order

	result := config.DB.Order(sortQuery).Find(&movies)

	if result.Error != nil {
		return nil, result.Error
	}

	// Actualizar los ratings de las películas
	for i := range movies {
		if err := updateMovieRating(&movies[i]); err != nil {
			fmt.Printf("[DEBUG-SERVICE] Error al actualizar rating de película %d: %v\n", movies[i].ID, err)
		}
	}

	return movies, nil
}
