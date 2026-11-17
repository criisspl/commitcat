// Package provider define la interfaz de generación y sus implementaciones.
package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Provider genera texto a partir de un prompt de sistema y un diff.
type Provider interface {
	// Generate devuelve la respuesta cruda del modelo.
	Generate(ctx context.Context, system, diff string) (string, error)
}

// Nombres de proveedor válidos.
const (
	NameOllama = "ollama"
	NameClaude = "claude"
)

// doJSON ejecuta la petición y devuelve el cuerpo si el estado es 2xx.
func doJSON(c *http.Client, req *http.Request) ([]byte, error) {
	if c == nil {
		c = http.DefaultClient
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("respuesta HTTP %d: %s", resp.StatusCode, body)
	}
	return body, nil
}
