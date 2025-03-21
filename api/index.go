package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cine_conecta_backend/config"
	"cine_conecta_backend/routes"
)

var router *gin.Engine

// initRouter se llama la primera vez que se invoque Handler
func initRouter() {
	// Modo release para producción
	gin.SetMode(gin.ReleaseMode)

	// Crea un router sin logger extra (opcional)
	router = gin.New()
	router.Use(gin.Recovery())

	// Conectar base de datos
	config.ConnectDB()

	// Registrar tus rutas
	routes.RegisterRoutes(router)
}

// Handler es la función que Vercel llama en cada request
func Handler(w http.ResponseWriter, r *http.Request) {
	if router == nil {
		initRouter()
	}
	router.ServeHTTP(w, r)
}
