package controllers

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/factories"
	"cine_conecta_backend/models"
	"cine_conecta_backend/services"
	"cine_conecta_backend/utils"
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo guardar el usuario")
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

	// Buscar el usuario por email
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email o contraseña incorrectos")
		return
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Email y/o contraseña incorrectos")
		return
	}

	// Generar token JWT
	token, err := utils.GenerateJWT(user.Name, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudo generar el token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesión iniciada correctamente",
		"token":   token,
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
