package controllers

import (
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/recommendations/me
func GetMyRecommendations(c *gin.Context) {
	// Obtener userID del token
	claims, _ := c.Get("claims")
	userID := claims.(*utils.Claims).UserID

	recommendations, err := services.GetRecommendedMovies(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener recomendaciones")
		return
	}

	// Obtener géneros favoritos del usuario
	favoriteGenres, _ := services.GetFavoriteGenres(userID)

	// Preparar respuesta enriquecida
	type RecommendationItem struct {
		Movie       services.MovieWithDetails `json:"movie"`
		Explanation string                    `json:"explanation"`
	}

	var response struct {
		Recommendations []RecommendationItem `json:"recommendations"`
		FavoriteGenres  []string             `json:"favorite_genres"`
		Message         string               `json:"message"`
	}

	// Generar explicaciones para cada recomendación
	for _, movie := range recommendations {
		explanation := "Recomendado porque "

		// Comprobar si el género coincide con los favoritos
		isGenreFavorite := false
		for _, genre := range favoriteGenres {
			if movie.Genre == genre {
				explanation += "te gustan películas de " + genre
				isGenreFavorite = true
				break
			}
		}

		// Si no es un género favorito, explicar por valoración
		if !isGenreFavorite {
			explanation += "es una película bien valorada"
		}

		// Añadir información sobre puntuación
		if movie.Rating > 0 {
			explanation += " y tiene una puntuación de " + strconv.FormatFloat(float64(movie.Rating), 'f', 1, 64) + "/10"
		}

		response.Recommendations = append(response.Recommendations, RecommendationItem{
			Movie:       services.EnrichMovie(movie),
			Explanation: explanation,
		})
	}

	response.FavoriteGenres = favoriteGenres

	// Personalizar mensaje según si tenemos o no recomendaciones
	if len(recommendations) > 0 {
		response.Message = "Hemos encontrado " + strconv.Itoa(len(recommendations)) + " películas que podrían interesarte basadas en tus comentarios"
	} else {
		response.Message = "Aún no tenemos suficiente información para recomendarte películas personalizadas. Comenta más películas para mejorar tus recomendaciones."
	}

	c.JSON(http.StatusOK, response)
}

// GET /api/recommendations/popular
func GetPopularRecommendations(c *gin.Context) {
	limit := 10
	limitParam := c.DefaultQuery("limit", "10")
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}

	recommendations, err := services.GetMoviesByPositiveSentiment(limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al obtener recomendaciones populares")
		return
	}

	// Preparar respuesta enriquecida
	type PopularRecommendationItem struct {
		Movie      services.MovieWithDetails `json:"movie"`
		Popularity string                    `json:"popularity"`
	}

	var response struct {
		Recommendations []PopularRecommendationItem `json:"popular_recommendations"`
		Message         string                      `json:"message"`
	}

	// Enriquecer cada recomendación
	for i, movie := range recommendations {
		enrichedMovie := services.EnrichMovie(movie)

		var popularity string
		switch {
		case i < 3:
			popularity = "Muy popular entre los usuarios"
		case i < 6:
			popularity = "Bastante popular entre los usuarios"
		default:
			popularity = "Popular entre los usuarios"
		}

		response.Recommendations = append(response.Recommendations, PopularRecommendationItem{
			Movie:      enrichedMovie,
			Popularity: popularity,
		})
	}

	// Mensaje personalizado
	if len(recommendations) > 0 {
		response.Message = "Estas son las " + strconv.Itoa(len(recommendations)) + " películas mejor valoradas por la comunidad"
	} else {
		response.Message = "Aún no hay suficientes valoraciones para generar recomendaciones populares"
	}

	c.JSON(http.StatusOK, response)
}
