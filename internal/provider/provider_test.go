package provider

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("ruta inesperada %s", r.URL.Path)
		}
		w.Write([]byte(`{"message":{"role":"assistant","content":"feat: hola"}}`))
	}))
	defer srv.Close()
	got, err := (&Ollama{BaseURL: srv.URL}).Generate(context.Background(), "sys", "diff")
	if err != nil || got != "feat: hola" {
		t.Fatalf("got %q, err %v", got, err)
	}
}

func TestOllamaError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := (&Ollama{BaseURL: srv.URL}).Generate(context.Background(), "s", "d"); err == nil {
		t.Fatal("esperaba error")
	}
}

func TestClaudeGenerate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "k" || r.Header.Get("anthropic-version") == "" {
			t.Errorf("cabeceras incorrectas: %v", r.Header)
		}
		w.Write([]byte(`{"content":[{"type":"text","text":"fix: arregla "},{"type":"text","text":"bug"}]}`))
	}))
	defer srv.Close()
	got, err := (&Claude{APIKey: "k", BaseURL: srv.URL}).Generate(context.Background(), "s", "d")
	if err != nil || got != "fix: arregla bug" {
		t.Fatalf("got %q, err %v", got, err)
	}
}

func TestClaudeSinKey(t *testing.T) {
	_, err := (&Claude{}).Generate(context.Background(), "s", "d")
	if !errors.Is(err, ErrNoAPIKey) {
		t.Fatalf("err = %v", err)
	}
}
