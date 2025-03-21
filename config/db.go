package config

import (
	"cine_conecta_backend/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "postgres://postgres:lolandia1@db.zufjxpgxyhphoygtxqit.supabase.co:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Error conectando a PostgreSQL: %v", err))
	}

	// Migrar modelos
	db.AutoMigrate(&models.User{})

	DB = db
	fmt.Println("✅ Conectado a PostgreSQL exitosamente")
}
