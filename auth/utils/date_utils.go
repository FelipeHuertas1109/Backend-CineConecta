package utils

import (
	"fmt"
	"time"
)

// ParseDate intenta analizar una fecha en formato ISO (YYYY-MM-DD) o RFC3339
func ParseDate(dateStr string) (time.Time, error) {
	// Intentar primero con formato ISO
	date, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		return date, nil
	}

	// Intentar con formato RFC3339
	date, err = time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return date, nil
	}

	// Intentar con formato RFC3339Nano
	date, err = time.Parse(time.RFC3339Nano, dateStr)
	if err == nil {
		return date, nil
	}

	// Intentar con formato con hora
	date, err = time.Parse("2006-01-02 15:04:05", dateStr)
	if err == nil {
		return date, nil
	}

	// Intentar con formato con hora sin segundos
	date, err = time.Parse("2006-01-02 15:04", dateStr)
	if err == nil {
		return date, nil
	}

	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", dateStr)
}
