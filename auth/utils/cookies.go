package utils

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetTokenCookie(c *gin.Context, token string) {
	isProduction := os.Getenv("ENV") == "production"

	c.SetCookie(
		"cine_token",                    // nombre
		token,                           // valor
		int((24 * time.Hour).Seconds()), // duración: 1 día
		"/",                             // ruta
		"",                              // dominio (vacío = actual)
		isProduction,                    // secure: solo true en producción
		true,                            // httpOnly: sí o sí
	)
}
