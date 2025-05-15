// movies/services/poster_service.go
package services

import (
	"cine_conecta_backend/config"
	"cine_conecta_backend/movies/models"
	"cine_conecta_backend/storage"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func UploadPoster(movieID uint, fileHeader *multipart.FileHeader) (string, error) {
	// Verificar si la película existe
	var movie models.Movie
	result := config.DB.First(&movie, movieID)
	if result.Error != nil {
		return "", errors.New("película no encontrada")
	}

	// Validaciones básicas
	if fileHeader.Size > 50<<20 { // 50 MB
		return "", errors.New("archivo supera los 50 MB")
	}

	mime := fileHeader.Header.Get("Content-Type")
	if mime != "image/jpeg" && mime != "image/png" && mime != "image/webp" {
		return "", errors.New("formato no permitido (solo JPEG, PNG o WEBP)")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo: %w", err)
	}
	defer file.Close()

	// Obtener extensión del nombre del archivo original
	filename := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(filename))

	// Verificar que la extensión coincida con el tipo MIME
	if !validateExtensionWithMime(ext, mime) {
		return "", errors.New("extensión de archivo no coincide con el tipo de contenido")
	}

	// Si no tiene extensión, usar la del tipo MIME
	if ext == "" {
		ext = extensionFromMime(mime)
	}

	// Construir la clave del archivo
	key := fmt.Sprintf("posters/%d%s", movieID, ext)

	// Borrar póster anterior si existe
	if movie.PosterURL != "" {
		// Podría implementarse la eliminación del archivo antiguo aquí
		fmt.Printf("Reemplazando póster anterior: %s\n", movie.PosterURL)
	}

	url, err := storage.UploadPoster(key, file, mime)
	if err != nil {
		return "", fmt.Errorf("error al subir póster: %w", err)
	}

	// Actualizar BD
	result = config.DB.Model(&movie).Update("poster_url", url)
	if result.Error != nil {
		return "", fmt.Errorf("error al actualizar BD: %w", result.Error)
	}

	return url, nil
}

func extensionFromMime(m string) string {
	switch m {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func validateExtensionWithMime(ext, mime string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return mime == "image/jpeg"
	case ".png":
		return mime == "image/png"
	case ".webp":
		return mime == "image/webp"
	case "": // Si no hay extensión, aceptamos cualquier MIME válido
		return mime == "image/jpeg" || mime == "image/png" || mime == "image/webp"
	default:
		return false
	}
}
