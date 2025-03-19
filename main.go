package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ruta básica
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "¡Hola, mundo desde Gin!",
		})
	})

	r.Run(":8080") // Inicia el servidor en el puerto 8080
}
