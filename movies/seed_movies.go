package movies

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/movies/services"
	"fmt"
	"time"
)

// SeedMovies crea películas de ejemplo en la base de datos
func SeedMovies() {
	// Primero, verificar si ya hay películas
	var count int64
	config.DB.Model(&models.Movie{}).Count(&count)
	if count > 0 {
		fmt.Println("✅ Ya existen películas en la base de datos, omitiendo seed")
		return
	}

	// Películas de ejemplo con sus géneros
	movies := []struct {
		Movie      models.Movie
		GenreNames []string
	}{
		{
			Movie: models.Movie{
				Title:       "El Padrino",
				Description: "La historia de la familia mafiosa Corleone liderada por Don Vito Corleone.",
				Director:    "Francis Ford Coppola",
				ReleaseDate: parseDate("1972-03-24"),
				Rating:      9.2,
				PosterURL:   "https://m.media-amazon.com/images/M/MV5BM2MyNjYxNmUtYTAwNi00MTYxLWJmNWYtYzZlODY3ZTk3OTFlXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_.jpg",
			},
			GenreNames: []string{"Drama", "Crimen"},
		},
		{
			Movie: models.Movie{
				Title:       "La La Land",
				Description: "Un pianista de jazz y una aspirante a actriz se enamoran mientras persiguen sus sueños.",
				Director:    "Damien Chazelle",
				ReleaseDate: parseDate("2016-12-09"),
				Rating:      8.0,
				PosterURL:   "https://m.media-amazon.com/images/M/MV5BMzUzNDM2NzM2MV5BMl5BanBnXkFtZTgwNTM3NTg4OTE@._V1_.jpg",
			},
			GenreNames: []string{"Musical", "Romance", "Drama"},
		},
		{
			Movie: models.Movie{
				Title:       "Interestelar",
				Description: "Un grupo de astronautas viaja a través de un agujero de gusano en busca de un nuevo hogar para la humanidad.",
				Director:    "Christopher Nolan",
				ReleaseDate: parseDate("2014-11-07"),
				Rating:      8.6,
				PosterURL:   "https://m.media-amazon.com/images/M/MV5BZjdkOTU3MDktN2IxOS00OGEyLWFmMjktY2FiMmZkNWIyODZiXkEyXkFqcGdeQXVyMTMxODk2OTU@._V1_.jpg",
			},
			GenreNames: []string{"Ciencia ficción", "Aventura", "Drama"},
		},
		{
			Movie: models.Movie{
				Title:       "Parásitos",
				Description: "La familia Kim, que vive en la pobreza, se infiltra en la vida de una familia adinerada.",
				Director:    "Bong Joon Ho",
				ReleaseDate: parseDate("2019-05-30"),
				Rating:      8.5,
				PosterURL:   "https://m.media-amazon.com/images/M/MV5BYWZjMjk3ZTItODQ2ZC00NTY5LWE0ZDYtZTI3MjcwN2Q5NTVkXkEyXkFqcGdeQXVyODk4OTc3MTY@._V1_.jpg",
			},
			GenreNames: []string{"Drama", "Comedia", "Thriller"},
		},
		{
			Movie: models.Movie{
				Title:       "Coco",
				Description: "Miguel sueña con ser músico, pero su familia se lo prohíbe. Desesperado por demostrar su talento, se encuentra en la Tierra de los Muertos.",
				Director:    "Lee Unkrich",
				ReleaseDate: parseDate("2017-11-22"),
				Rating:      8.4,
				PosterURL:   "https://m.media-amazon.com/images/M/MV5BYjQ5NjM0Y2YtNjZkNC00ZDhkLWJjMWItN2QyNzFkMDE3ZjAxXkEyXkFqcGdeQXVyODIxMzk5NjA@._V1_.jpg",
			},
			GenreNames: []string{"Animación", "Aventura", "Familia"},
		},
	}

	// Insertar las películas
	for _, item := range movies {
		if err := services.CreateMovieWithGenres(&item.Movie, item.GenreNames); err != nil {
			fmt.Printf("❌ Error al crear la película %s: %v\n", item.Movie.Title, err)
		} else {
			fmt.Printf("✅ Película creada: %s con géneros: %v\n", item.Movie.Title, item.GenreNames)
		}
	}

	fmt.Println("✅ Seed de películas completado")
}

// Función auxiliar para parsear fechas
func parseDate(dateStr string) time.Time {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}
