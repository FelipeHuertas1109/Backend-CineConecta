package config

import (
	authModels "cine_conecta_backend/auth/models"
	commentModels "cine_conecta_backend/comments/models"
	movieModels "cine_conecta_backend/movies/models"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Cargar .env en desarrollo (ignorar si ya está en Vercel)
	if os.Getenv("VERCEL") == "" { // Vercel define esta variable automáticamente
		// Obtener directorio actual para depuración
		dir, err := os.Getwd()
		if err != nil {
			log.Printf("❌ [DB] Error al obtener directorio de trabajo: %v", err)
		} else {
			log.Printf("📂 [DB] Directorio de trabajo al cargar .env: %s", dir)
		}

		// Intentar cargar .env desde diferentes ubicaciones
		locations := []string{
			".env",                        // En la raíz del proyecto
			"../.env",                     // Un nivel arriba
			filepath.Join(dir, ".env"),    // Ruta absoluta
			filepath.Join(dir, "../.env"), // Un nivel arriba (absoluto)
		}

		loaded := false
		for _, location := range locations {
			log.Printf("🔍 [DB] Intentando cargar .env desde: %s", location)
			if _, statErr := os.Stat(location); os.IsNotExist(statErr) {
				log.Printf("❌ [DB] No existe archivo en: %s", location)
				continue
			}

			err := godotenv.Load(location)
			if err != nil {
				log.Printf("⚠️ [DB] No se pudo cargar .env desde %s: %v", location, err)
			} else {
				log.Printf("✅ [DB] Archivo .env cargado con éxito desde: %s", location)
				loaded = true
				break
			}
		}

		if !loaded {
			log.Println("⚠️ [DB] No se pudo cargar el archivo .env desde ninguna ubicación, usando variables del sistema")
		}
	}

	// Leer DATABASE_URL de entorno o del .env
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("❌ [DB] La variable DATABASE_URL no está configurada")
		log.Println("💡 [DB] Asegúrate de que tu archivo .env contiene: DATABASE_URL=postgresql://usuario:contraseña@localhost:5432/nombre_db")
		panic("❌ Error: La variable DATABASE_URL no está configurada")
	} else {
		// Mostrar parte de la URL para depuración (ocultando contraseña)
		dsnPreview := dsn
		if len(dsnPreview) > 30 {
			dsnPreview = dsnPreview[:30] + "..."
		}
		log.Printf("✅ [DB] Variable DATABASE_URL encontrada: %s", dsnPreview)
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
		&movieModels.Genre{},
		&movieModels.Like{},
		&commentModels.Comment{},
		&commentModels.RecommendationDataset{})

	// Crear índice único para asegurar que un usuario solo pueda comentar una vez por película
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_comments_user_movie ON comments (user_id, movie_id)")

	// Crear índice único para asegurar que un usuario solo pueda dar me gusta una vez por película
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_movie_likes_user_movie ON movie_likes (user_id, movie_id)")

	DB = db
}

// migrateGenres ya no es necesaria porque ahora los géneros son simplemente strings
func migrateGenres(db *gorm.DB) {
	// Esta función ya no hace nada porque los géneros son cadenas de texto simples
	// Se mantiene por compatibilidad con código existente
}
