package main

import (
	"log"
	"net/http"

	handler "cine_conecta_backend/api"
)

func main() {
	// Conexión a la base de datos
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", http.HandlerFunc(handler.Handler))
}
