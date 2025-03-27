package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetTokenCookie(c *gin.Context, token string) {
	cookie := &http.Cookie{
		Name:     "cine_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                  // Permitir acceso desde JavaScript
		Secure:   true,                  // Siempre usar HTTPS
		SameSite: http.SameSiteNoneMode, // Permitir cross-site
		MaxAge:   int((24 * time.Hour).Seconds()),
	}

	http.SetCookie(c.Writer, cookie)
}
