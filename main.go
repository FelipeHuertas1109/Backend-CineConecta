package main

import (
	"cine_conecta_backend/auth/routes"
	"cine_conecta_backend/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conexión a la base de datos
	config.ConnectDB()

	// Crear router
	r := gin.Default()

	// Registrar rutas
	routes.RegisterRoutes(r)

	// Iniciar servidor local
	r.Run(":8080") // http://localhost:8080
}
