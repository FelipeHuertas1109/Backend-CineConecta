package utils

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func SetTokenCookie(c *gin.Context, token string) {
	isProduction := os.Getenv("ENV") == "production"

	cookie := &http.Cookie{
		Name:     "cine_token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   isProduction,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	}

	if isProduction {
		cookie.Domain = ""
		cookie.SameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, cookie)
}
