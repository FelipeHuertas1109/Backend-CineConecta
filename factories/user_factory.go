package factories

import "cine_conecta_backend/models"

func NewUser(name, email, hashedPassword string) *models.User {
	return &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}
}
