package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"errors"
)

// CreateMovie guarda una nueva película en la base de datos.
func CreateMovie(movie *models.Movie) error {
	return config.DB.Create(movie).Error
}

// GetMovies obtiene todas las películas.
func GetMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Preload("Genres").Find(&movies).Error
	return movies, err
}

// GetRecentMovies obtiene las películas más recientes según su fecha de lanzamiento.
func GetRecentMovies(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := config.DB.Preload("Genres").Order("release_date DESC").Limit(limit).Find(&movies).Error
	return movies, err
}

// GetMovieByID obtiene una película por su ID.
func GetMovieByID(id uint) (models.Movie, error) {
	var movie models.Movie
	err := config.DB.Preload("Genres").First(&movie, id).Error
	return movie, err
}

// UpdateMovie actualiza una película existente.
func UpdateMovie(movie *models.Movie) error {
	return config.DB.Save(movie).Error
}

// DeleteMovie elimina una película por su ID.
func DeleteMovie(id uint) error {
	return config.DB.Delete(&models.Movie{}, id).Error
}

// CreateMovieWithGenres crea una película con sus géneros asociados
func CreateMovieWithGenres(movie *models.Movie, genreNames []string) error {
	// Iniciar transacción
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Crear la película primero
	if err := tx.Create(movie).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Procesar géneros
	if len(genreNames) > 0 {
		for _, genreName := range genreNames {
			// Buscar o crear género
			genre, err := CreateGenre(genreName)
			if err != nil {
				tx.Rollback()
				return err
			}

			if genre != nil {
				// Asociar género con película
				if err := tx.Model(movie).Association("Genres").Append(genre); err != nil {
					tx.Rollback()
					return err
				}
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

	// Actualizar la película
	if err := tx.Save(movie).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Si se especificaron géneros, actualizar la relación
	if genreNames != nil {
		// Limpiar géneros existentes
		if err := tx.Model(movie).Association("Genres").Clear(); err != nil {
			tx.Rollback()
			return err
		}

		// Agregar nuevos géneros
		for _, genreName := range genreNames {
			genre, err := CreateGenre(genreName)
			if err != nil {
				tx.Rollback()
				return err
			}

			if genre != nil {
				if err := tx.Model(movie).Association("Genres").Append(genre); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// Commit de la transacción
	return tx.Commit().Error
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

	result := config.DB.Preload("Genres").Order(sortQuery).Find(&movies)

	if result.Error != nil {
		return nil, result.Error
	}

	return movies, nil
}
