package controllers

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/factories"
	"cine_conecta_backend/models"
	"cine_conecta_backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Registrar usuario
func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	// Crear el usuario con el factory
	user := factories.NewUser(input.Name, input.Email, string(hashedPassword))

	// Guardar el usuario con el servicio
	if err := services.SaveUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado correctamente"})
}

// Iniciar sesión
func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar por email
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
		return
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email o contraseña incorrectos"})
		return
	}

	// Inicio de sesión exitoso
	c.JSON(http.StatusOK, gin.H{"message": "Sesión iniciada correctamente", "user": user})
}

func GetAllUsers(c *gin.Context) {
	var users []models.User

	result := config.DB.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los usuarios"})
		return
	}

	// Evitar devolver las contraseñas hasheadas
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}
