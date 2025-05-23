package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	handler "cine_conecta_backend/api"
	"cine_conecta_backend/comments/services"

	"github.com/joho/godotenv"
)

func main() {
	// Procesar flags de línea de comandos
	checkHFToken := flag.Bool("check-hf", false, "Verificar token de HuggingFace")
	flag.Parse()

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

	// Si se solicitó verificar el token de HuggingFace
	if *checkHFToken {
		log.Println("Verificando token de HuggingFace...")
		if services.VerifyHFToken() {
			log.Println("✅ Token de HuggingFace verificado correctamente")
			os.Exit(0)
		} else {
			log.Println("❌ Error en la verificación del token de HuggingFace")
			os.Exit(1)
		}
		return // No continuar con el servidor
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
