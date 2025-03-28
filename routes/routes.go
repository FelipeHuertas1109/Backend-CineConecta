package routes

import (
	routesAuth "cine_conecta_backend/auth/routes"
	routesMovies "cine_conecta_backend/movies/routes"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	routesMovies.RegisterMovieRoutes(r)
	routesAuth.RegisterAuthRoutes(r)
}
