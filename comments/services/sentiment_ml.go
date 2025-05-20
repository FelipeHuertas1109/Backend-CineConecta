// comments/services/sentiment_ml.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Modelo español entrenado en 500 M tuits
const hfURL = "https://api-inference.huggingface.co/models/pysentimiento/robertuito-sentiment-analysis"

// respuesta HF → [{"label":"POS","score":0.83}, {"label":"NEU",…}, {"label":"NEG",…}]
type hfItem struct {
	Label string  `json:"label"` // POS, NEU, NEG
	Score float64 `json:"score"` // probabilidad
}

// SentimentScoreML devuelve una nota 1-5
func SentimentScoreML(text string) (float64, error) {
	debug := os.Getenv("SENTIMENT_DEBUG") == "true"

	token := os.Getenv("HF_TOKEN")
	if token == "" {
		if debug {
			fmt.Println("[HF] token VACÍO – se usará heurístico")
		}
		return 0, fmt.Errorf("HF_TOKEN vacío")
	}

	// ----- llamada a la API HF -----
	payload, _ := json.Marshal(map[string]string{"inputs": text})
	req, _ := http.NewRequest("POST", hfURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			fmt.Println("[HF] error de red:", err)
		}
		return 0, err
	}
	defer resp.Body.Close()

	if debug {
		fmt.Printf("[HF] HTTP %d %s\n", resp.StatusCode, resp.Status)
	}

	// ----- decodificar -----
	var out []hfItem
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || len(out) != 3 {
		return 0, fmt.Errorf("respuesta HF inesperada")
	}

	var pPos, pNeu, pNeg float64
	for _, it := range out {
		switch it.Label {
		case "POS":
			pPos = it.Score
		case "NEU":
			pNeu = it.Score
		case "NEG":
			pNeg = it.Score
		}
	}

	// media ponderada en escala 1-5
	score := 1*pNeg + 3*pNeu + 5*pPos

	if debug {
		fmt.Printf("[HF] POS %.2f  NEU %.2f  NEG %.2f  → score %.2f\n",
			pPos, pNeu, pNeg, score)
	}

	return score, nil
}
