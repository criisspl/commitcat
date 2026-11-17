package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	// DefaultOllamaURL es la dirección local por defecto de Ollama.
	DefaultOllamaURL = "http://localhost:11434"
	// DefaultOllamaModel es el modelo por defecto.
	DefaultOllamaModel = "llama3.2"
)

// Ollama usa un servidor Ollama local.
type Ollama struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// Generate implementa Provider.
func (o *Ollama) Generate(ctx context.Context, system, diff string) (string, error) {
	base := strings.TrimRight(o.BaseURL, "/")
	if base == "" {
		base = DefaultOllamaURL
	}
	model := o.Model
	if model == "" {
		model = DefaultOllamaModel
	}
	payload, err := json.Marshal(map[string]any{
		"model":  model,
		"stream": false,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": diff},
		},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/chat", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	body, err := doJSON(o.Client, req)
	if err != nil {
		return "", fmt.Errorf("ollama: %w (¿está corriendo `ollama serve`?)", err)
	}
	var out struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("ollama: respuesta inválida: %w", err)
	}
	if strings.TrimSpace(out.Message.Content) == "" {
		return "", fmt.Errorf("ollama: respuesta vacía")
	}
	return out.Message.Content, nil
}
