package factories

import "cine_conecta_backend/auth/models"

func NewUser(name, email, hashedPassword, role string) *models.User {
	return &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}
}
