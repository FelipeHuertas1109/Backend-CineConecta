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

// SentimentScoreML devuelve un rating 1-5 usando la media ponderada.
// Si el modelo está calentando reintenta hasta 3 veces.
func SentimentScoreML(text string) (float64, error) {
	debug := os.Getenv("SENTIMENT_DEBUG") == "true"
	tk := os.Getenv("HF_TOKEN")
	if tk == "" {
		return 0, fmt.Errorf("HF_TOKEN vacío")
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

		client := &http.Client{Timeout: 25 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if debug {
			fmt.Printf("[HF] HTTP %d %s (try %d)\n", resp.StatusCode, resp.Status, attempt+1)
		}
		defer resp.Body.Close()

		// ---- leer cuerpo bruto ----
		var raw json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			lastErr = fmt.Errorf("decode raw: %w", err)
			continue
		}

		// ---- 1) ¿error de HF? ----
		var apiErr hfErr
		if json.Unmarshal(raw, &apiErr) == nil && apiErr.Error != "" {
			lastErr = fmt.Errorf("HF error: %s", apiErr.Error)
			if debug {
				fmt.Println("[HF] ❗", apiErr.Error)
			}
			// esperar y reintentar mientras el modelo carga
			continue
		}

		// ---- 2) lista de listas (batch) ----
		var batch [][]hfStar
		if json.Unmarshal(raw, &batch) == nil && len(batch) > 0 && len(batch[0]) > 0 {
			return ponderate(batch[0], debug), nil
		}

		// ---- 3) lista simple ----
		var single []hfStar
		if json.Unmarshal(raw, &single) == nil && len(single) > 0 {
			return ponderate(single, debug), nil
		}

		lastErr = fmt.Errorf("respuesta HF inesperada: %s", string(raw))
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
