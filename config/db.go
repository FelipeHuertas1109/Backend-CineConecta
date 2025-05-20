package config

import (
	authModels "cine_conecta_backend/auth/models"
	commentModels "cine_conecta_backend/comments/models"
	movieModels "cine_conecta_backend/movies/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Cargar .env en desarrollo (ignorar si ya está en Vercel)
	if os.Getenv("VERCEL") == "" { // Vercel define esta variable automáticamente
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️  Advertencia: No se pudo cargar el archivo .env, usando variables del sistema")
		}
	}

	// Leer DATABASE_URL de entorno o del .env
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("❌ Error: La variable DATABASE_URL no está configurada")
	}

	// Conectar a la base de datos PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ Error conectando a la base de datos: %v", err))
	}

	fmt.Println("✅ Conectado a PostgreSQL correctamente")
	db.AutoMigrate(
		&authModels.User{},
		&movieModels.Movie{},
		&commentModels.Comment{},
		&commentModels.RecommendationDataset{})

	// Crear índice único para asegurar que un usuario solo pueda comentar una vez por película
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_comments_user_movie ON comments (user_id, movie_id)")

	DB = db
}
