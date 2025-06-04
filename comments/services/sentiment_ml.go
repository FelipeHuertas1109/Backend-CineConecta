package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const hfURL = "https://api-inference.huggingface.co/models/nlptown/bert-base-multilingual-uncased-sentiment"

type hfStar struct {
	Label string  `json:"label"` // "1 star" … "5 stars"
	Score float64 `json:"score"`
}

type hfErr struct {
	Error string `json:"error"`
}

// VerifyHFToken verifica que el token de HuggingFace esté configurado correctamente
// y muestra información útil para depuración
func VerifyHFToken() bool {
	fmt.Println("\n🔍 VERIFICACIÓN DEL TOKEN DE HUGGINGFACE 🔍")

	// Verificar token de HuggingFace
	tk := os.Getenv("HF_TOKEN")
	if tk == "" {
		fmt.Println("❌ ERROR: La variable HF_TOKEN no está definida o está vacía")
		fmt.Println("Por favor, agrega la siguiente línea a tu archivo .env:")
		fmt.Println("HF_TOKEN=hf_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
		fmt.Println("Donde 'hf_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' es tu token de HuggingFace")
		return false
	}

	fmt.Println("--- Token de HuggingFace ---")
	fmt.Printf("Longitud: %d caracteres\n", len(tk))

	// Mostrar primeros y últimos caracteres del token
	if len(tk) > 10 {
		fmt.Printf("Primeros 10 caracteres: %s\n", tk[:10])
		fmt.Printf("Últimos 5 caracteres: %s\n", tk[len(tk)-5:])
	} else {
		fmt.Printf("Token completo (muy corto): %s\n", tk)
	}

	// Verificar formato típico de token HF (hf_...)
	if len(tk) > 3 && tk[:3] == "hf_" {
		fmt.Println("✅ El formato del token parece correcto (comienza con 'hf_')")
	} else {
		fmt.Println("⚠️ El token no tiene el formato esperado (debería comenzar con 'hf_')")
		fmt.Println("Un token válido de HuggingFace debe comenzar con 'hf_'")
	}

	// Realizar una prueba de conexión
	fmt.Println("\n🔄 Realizando prueba de conexión a HuggingFace...")

	body, _ := json.Marshal(map[string]string{"inputs": "Prueba de conexión"})
	req, _ := http.NewRequest("POST", hfURL, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tk)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("❌ Error de conexión: %v\n", err)
		return false
	}

	fmt.Printf("📡 Respuesta HTTP: %d %s\n", resp.StatusCode, resp.Status)

	// Verificar código de respuesta
	if resp.StatusCode == 200 {
		fmt.Println("✅ Conexión exitosa a HuggingFace API")
		return true
	} else if resp.StatusCode == 401 {
		fmt.Println("❌ Error de autenticación: Token inválido o expirado")
		fmt.Println("Por favor, verifica que tu token sea correcto y esté vigente")
		return false
	} else {
		fmt.Printf("⚠️ Respuesta inesperada del servidor: %d\n", resp.StatusCode)

		// Intentar leer el cuerpo de la respuesta para más información
		var raw json.RawMessage
		if json.NewDecoder(resp.Body).Decode(&raw) == nil {
			fmt.Printf("📄 Detalles: %s\n", string(raw))
		}
		resp.Body.Close()

		return false
	}
}

// SentimentScoreML devuelve un rating 1-5 usando la media ponderada.
// Si el modelo está calentando reintenta hasta 3 veces.
func SentimentScoreML(text string) (float64, error) {
	debug := os.Getenv("SENTIMENT_DEBUG") == "true"
	tk := os.Getenv("HF_TOKEN")

	// Validar token
	if tk == "" {
		if debug {
			fmt.Println("[HF] ❌ Error: HF_TOKEN no está configurado en las variables de entorno")
		}
		return 0, fmt.Errorf("HF_TOKEN vacío")
	}

	// Mostrar primeros caracteres del token para depuración
	if debug {
		tokenPreview := tk
		if len(tokenPreview) > 10 {
			tokenPreview = tokenPreview[:10] + "..."
		}
		fmt.Printf("[HF] 🔑 Usando token: %s\n", tokenPreview)
	}

	body, _ := json.Marshal(map[string]string{"inputs": text})
	var lastErr error

	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(1<<attempt) * time.Second) // back-off 1 s, 2 s, 4 s
		}

		req, _ := http.NewRequest("POST", hfURL, bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tk)
		req.Header.Set("Content-Type", "application/json")

		if debug {
			fmt.Printf("[HF] 🔄 Intento %d: Enviando solicitud a %s\n", attempt+1, hfURL)
			fmt.Printf("[HF] 📝 Texto a analizar: '%s'\n", text)
		}

		client := &http.Client{Timeout: 25 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if debug {
				fmt.Printf("[HF] ❌ Error de red: %v\n", err)
			}
			continue
		}

		if debug {
			fmt.Printf("[HF] 📡 HTTP %d %s (intento %d)\n", resp.StatusCode, resp.Status, attempt+1)
		}
		defer resp.Body.Close()

		// ---- leer cuerpo bruto ----
		var raw json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			lastErr = fmt.Errorf("decode raw: %w", err)
			if debug {
				fmt.Printf("[HF] ❌ Error al decodificar respuesta: %v\n", err)
			}
			continue
		}

		if debug {
			fmt.Printf("[HF] 📄 Respuesta: %s\n", string(raw))
		}

		// ---- 1) ¿error de HF? ----
		var apiErr hfErr
		if json.Unmarshal(raw, &apiErr) == nil && apiErr.Error != "" {
			lastErr = fmt.Errorf("HF error: %s", apiErr.Error)
			if debug {
				fmt.Printf("[HF] ❗ Error de API: %s\n", apiErr.Error)
			}
			// esperar y reintentar mientras el modelo carga
			continue
		}

		// ---- 2) lista de listas (batch) ----
		var batch [][]hfStar
		if json.Unmarshal(raw, &batch) == nil && len(batch) > 0 && len(batch[0]) > 0 {
			result := ponderate(batch[0], debug)
			if debug {
				fmt.Printf("[HF] ✅ Análisis completado: %.1f/5.0\n", result)
			}
			return result, nil
		}

		// ---- 3) lista simple ----
		var single []hfStar
		if json.Unmarshal(raw, &single) == nil && len(single) > 0 {
			result := ponderate(single, debug)
			if debug {
				fmt.Printf("[HF] ✅ Análisis completado: %.1f/5.0\n", result)
			}
			return result, nil
		}

		lastErr = fmt.Errorf("respuesta HF inesperada: %s", string(raw))
		if debug {
			fmt.Printf("[HF] ❌ Formato de respuesta inesperado\n")
		}
	}

	if debug {
		fmt.Printf("[HF] ❌ Todos los intentos fallaron, último error: %v\n", lastErr)
	}
	return 0, lastErr
}

// ponderate convierte las probabilidades en nota 1-5
func ponderate(stars []hfStar, debug bool) float64 {
	var sum, prob float64
	for _, s := range stars {
		if len(s.Label) == 0 {
			continue
		}
		n := int(s.Label[0] - '0') // '1'..'5'
		sum += float64(n) * s.Score
		prob += s.Score
	}
	rating := sum / prob
	if debug {
		fmt.Printf("[HF] ponderado %.2f  (prob %.2f)\n", rating, prob)
	}
	return rating
}
