package handler

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// initRouter configura el router de Gin, conecta a la base de datos y registra las rutas.
func initRouter() {
	// Configura Gin en modo Release para producción
	gin.SetMode(gin.ReleaseMode)

	// Crea un nuevo router sin Logger adicional y con middleware de recuperación
	router = gin.New()
	router.Use(gin.Recovery())

	// Conecta a la base de datos (si tienes la función definida en config)
	config.ConnectDB()

	// Registra las rutas definidas en el paquete routes
	routes.RegisterRoutes(router)
}

// Handler es la función de entrada que Vercel invoca para cada request.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Si el router aún no se ha inicializado, se inicializa
	if router == nil {
		initRouter()
	}
	// Delegamos la petición a Gin para que la procese según las rutas definidas
	router.ServeHTTP(w, r)
}
