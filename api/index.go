package handler

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/routes"
	"net/http"

	"github.com/gin-contrib/cors"
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

	// 🔐 Middleware CORS para permitir peticiones desde el frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3100", "https://cineconecta.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 horas
	}))

	// Middleware adicional para asegurar que los headers de CORS estén presentes
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Next()
	})

	// Conecta a la base de datos
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
