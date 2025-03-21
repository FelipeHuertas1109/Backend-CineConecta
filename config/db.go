package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Error conectando a PostgreSQL: %v", err))
	}
	// Migrar tus modelos aquí, por ejemplo:
	// db.AutoMigrate(&models.User{}, &models.Movie{})
	DB = db
}
