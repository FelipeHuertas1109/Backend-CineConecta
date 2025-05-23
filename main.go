package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	handler "cine_conecta_backend/api"

	"github.com/joho/godotenv"
)

func main() {
	// Mostrar directorio de trabajo actual
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("❌ Error al obtener directorio de trabajo: %v", err)
	} else {
		log.Printf("📂 Directorio de trabajo actual: %s", dir)
	}

	// Verificar si el archivo .env existe
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Printf("❌ El archivo .env no existe en la ruta: %s", envPath)
	} else {
		log.Printf("✅ Archivo .env encontrado en: %s", envPath)

		// Intentar cargar manualmente el archivo .env
		err := godotenv.Load(envPath)
		if err != nil {
			log.Printf("❌ Error al cargar el archivo .env: %v", err)

			// Leer el contenido del archivo para depuración
			content, readErr := os.ReadFile(envPath)
			if readErr != nil {
				log.Printf("❌ No se pudo leer el archivo .env: %v", readErr)
			} else {
				log.Printf("📄 Contenido del archivo .env (primeros 100 caracteres): %s", string(content[:min(len(content), 100)]))
			}
		} else {
			log.Printf("✅ Archivo .env cargado manualmente con éxito")

			// Verificar si la variable DATABASE_URL está definida
			if dbURL := os.Getenv("DATABASE_URL"); dbURL == "" {
				log.Printf("❌ La variable DATABASE_URL no está definida en el archivo .env")
			} else {
				log.Printf("✅ Variable DATABASE_URL encontrada (primeros 20 caracteres): %s...", dbURL[:min(len(dbURL), 20)])
			}
		}
	}

	// Conexión a la base de datos
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", http.HandlerFunc(handler.Handler))
}

// Función auxiliar para obtener el mínimo de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
