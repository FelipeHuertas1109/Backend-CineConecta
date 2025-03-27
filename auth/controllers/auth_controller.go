package controllers

import (
	"cine_conecta_backend/auth/factories"
	"cine_conecta_backend/auth/models"
	"cine_conecta_backend/auth/services"
	"cine_conecta_backend/auth/utils"
	"cine_conecta_backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Registrar usuario
func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al encriptar la contraseña")
		return
	}

	// Asignar rol por defecto o como admin
	role := "user"
	if input.Email == "fhuertas@unillanos.edu.co" {
		role = "admin"
	}

	// Crear el usuario con el factory
	user := factories.NewUser(input.Name, input.Email, string(hashedPassword), role)

	// Guardar el usuario con el servicio
	if err := services.SaveUser(user); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Correo ya utilizado en otra cuenta")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado correctamente"})
}

// Iniciar sesión
func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Buscar usuario
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email o contraseña incorrectos")
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email o contraseña incorrectos")
		return
	}

	// Generar token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo generar el token")
		return
	}

	// Comentado: Ya no establecemos la cookie desde el backend
	// utils.SetTokenCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesión iniciada correctamente",
		"token":   token,
	})
}

func Logout(c *gin.Context) {
	// Expirar la cookie 'cine_token'
	cookie := &http.Cookie{
		Name:     "cine_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1, // Eliminar la cookie
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesión cerrada correctamente",
	})
}

// Obtener todos los usuarios (solo admin)
func GetAllUsers(c *gin.Context) {
	var users []models.User

	result := config.DB.Find(&users)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudieron obtener los usuarios")
		return
	}

	// Ocultar contraseñas
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

func DeleteAllUsers(c *gin.Context) {
	result := config.DB.Exec("DELETE FROM users WHERE role != ?", "admin")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudieron eliminar los usuarios",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todos los usuarios no admin han sido eliminados"})
}

func GetProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "No autorizado")
		return
	}

	// Extraer ID del token
	userClaims := claims.(*utils.Claims)

	var user models.User
	if err := config.DB.First(&user, userClaims.UserID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Usuario no encontrado")
		return
	}

	// Devolver solo los campos necesarios
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func VerifyToken(c *gin.Context) {
	// La cookie ya está siendo verificada por el middleware de autenticación
	// Si llegamos aquí, significa que el token es válido
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
	})
}
