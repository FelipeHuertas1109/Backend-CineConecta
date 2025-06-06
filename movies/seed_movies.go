package main

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/movies/services"
	"fmt"
	"time"
)

func main() {
	config.ConnectDB()

	// Crear películas con géneros
	createMovieWithGenres(
		"Inception",
		"Un ladrón que roba secretos a través de los sueños debe realizar una misión inversa: implantar una idea.",
		"Christopher Nolan",
		time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC),
		4.8,
		"",
		[]string{"Ciencia Ficción", "Acción", "Thriller"},
	)

	createMovieWithGenres(
		"The Godfather",
		"La historia de una familia mafiosa italoamericana que lucha por mantener su imperio.",
		"Francis Ford Coppola",
		time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
		4.9,
		"",
		[]string{"Crimen", "Drama"},
	)

	createMovieWithGenres(
		"Parasite",
		"Una familia pobre se infiltra en una casa rica con consecuencias inesperadas.",
		"Bong Joon-ho",
		time.Date(2019, 5, 30, 0, 0, 0, 0, time.UTC),
		4.6,
		"",
		[]string{"Drama", "Thriller", "Comedia Negra"},
	)

	createMovieWithGenres(
		"The Dark Knight",
		"Batman se enfrenta al Joker, un enemigo que desata el caos en Gotham.",
		"Christopher Nolan",
		time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC),
		4.9,
		"",
		[]string{"Acción", "Crimen", "Drama", "Superhéroes"},
	)

	createMovieWithGenres(
		"Pulp Fiction",
		"Historias entrelazadas de crimen, redención y violencia en Los Ángeles.",
		"Quentin Tarantino",
		time.Date(1994, 10, 14, 0, 0, 0, 0, time.UTC),
		4.7,
		"",
		[]string{"Crimen", "Drama", "Comedia Negra"},
	)

	createMovieWithGenres(
		"Spirited Away",
		"Una niña entra en un mundo espiritual y debe rescatar a sus padres.",
		"Hayao Miyazaki",
		time.Date(2001, 7, 20, 0, 0, 0, 0, time.UTC),
		4.8,
		"",
		[]string{"Animación", "Fantasía", "Aventura"},
	)

	fmt.Println("✅ Películas de ejemplo creadas correctamente")
}

// Función auxiliar para crear películas con géneros
func createMovieWithGenres(title, description, director string, releaseDate time.Time, rating float32, posterURL string, genres []string) {
	// Verificar si la película ya existe
	var existingMovie models.Movie
	result := config.DB.Where("title = ?", title).First(&existingMovie)
	if result.Error == nil {
		fmt.Printf("La película '%s' ya existe, omitiendo...\n", title)
		return
	}

	// Crear la película
	movie := models.Movie{
		Title:       title,
		Description: description,
		Director:    director,
		ReleaseDate: releaseDate,
		Rating:      rating,
		PosterURL:   posterURL,
		Genre:       joinGenres(genres), // Para compatibilidad con código existente
	}

	// Guardar la película con sus géneros
	err := services.CreateMovieWithGenres(&movie, genres)
	if err != nil {
		fmt.Printf("Error al crear la película '%s': %v\n", title, err)
		return
	}

	fmt.Printf("Película '%s' creada con %d géneros\n", title, len(genres))
}

// Función auxiliar para unir géneros en una cadena separada por comas (para compatibilidad)
func joinGenres(genres []string) string {
	if len(genres) == 0 {
		return ""
	}
	result := genres[0]
	for i := 1; i < len(genres); i++ {
		result += ", " + genres[i]
	}
	return result
}
