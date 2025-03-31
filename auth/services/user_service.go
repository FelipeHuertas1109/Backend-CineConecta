package services

import (
	"cine_conecta_backend/auth/models"
	"cine_conecta_backend/config"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SaveUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// Enviar correo de restablecimiento de contraseña
func SendResetPasswordEmail(email, token string) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	frontendURL := os.Getenv("FRONTEND_URL")

	// Convertir el puerto de string a int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("Error al convertir el puerto SMTP: %v", err)
	}

	// Configurar el correo
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Restablece tu contraseña")
	mailer.SetBody("text/html", fmt.Sprintf(`
		<p>Haz clic en el siguiente enlace para restablecer tu contraseña:</p>
		<a href="%s/reset-password?token=%s">Restablecer contraseña</a>
	`, frontendURL, token))

	// Configuración del servidor SMTP
	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Enviar el correo
	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}
