package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	// DefaultClaudeURL es la URL base de la API de Anthropic.
	DefaultClaudeURL = "https://api.anthropic.com"
	// DefaultClaudeModel es el modelo por defecto.
	DefaultClaudeModel = "claude-haiku-5-5"
)

// ErrNoAPIKey se devuelve si se pide Claude sin clave.
var ErrNoAPIKey = errors.New("falta ANTHROPIC_API_KEY para usar el proveedor claude")

// Claude usa la API de Anthropic (opcional, requiere API key).
type Claude struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

// Generate implementa Provider.
func (c *Claude) Generate(ctx context.Context, system, diff string) (string, error) {
	if c.APIKey == "" {
		return "", ErrNoAPIKey
	}
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = DefaultClaudeURL
	}
	model := c.Model
	if model == "" {
		model = DefaultClaudeModel
	}
	payload, err := json.Marshal(map[string]any{
		"model":      model,
		"max_tokens": 512,
		"system":     system,
		"messages":   []map[string]string{{"role": "user", "content": diff}},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	body, err := doJSON(c.Client, req)
	if err != nil {
		return "", fmt.Errorf("claude: %w", err)
	}
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("claude: respuesta inválida: %w", err)
	}
	var sb strings.Builder
	for _, b := range out.Content {
		if b.Type == "text" {
			sb.WriteString(b.Text)
		}
	}
	if strings.TrimSpace(sb.String()) == "" {
		return "", fmt.Errorf("claude: respuesta vacía")
	}
	return sb.String(), nil
}
