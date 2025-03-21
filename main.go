package main

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
