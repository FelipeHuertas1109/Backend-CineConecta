package routes

import (
	routesAuth "cine_conecta_backend/auth/routes"
	routesComment "cine_conecta_backend/comments/routes"
	routesMovies "cine_conecta_backend/movies/routes"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	routesMovies.RegisterMovieRoutes(r)
	routesAuth.RegisterAuthRoutes(r)
	routesComment.RegisterCommentRoutes(r)
}
