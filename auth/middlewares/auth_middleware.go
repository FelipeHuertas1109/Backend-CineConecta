package middlewares

import (
	"cine_conecta_backend/auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Lee el token JWT desde la cookie "cine_token"
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener cookie con el token
		tokenString, err := c.Cookie("cine_token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se encontró el token"})
			c.Abort()
			return
		}

		// Validar token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Validar que el rol sea "admin"
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: solo administradores"})
			c.Abort()
			return
		}

		// Guardar info en el contexto (por si se necesita)
		c.Set("userName", claims.Name)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
