package cmd

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/criisspl/commitcat/internal/git"
	"github.com/criisspl/commitcat/internal/provider"
)

// falso implementa provider.Provider sin red.
type falso struct {
	resp string
	err  error
	diff string
}

func (f *falso) Generate(_ context.Context, _, diff string) (string, error) {
	f.diff = diff
	return f.resp, f.err
}

func repoTemporal(t *testing.T, staged bool) string {
	t.Helper()
	dir := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "t@t.test"},
		{"config", "user.name", "t"},
		{"config", "commit.gpgsign", "false"},
	} {
		if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hola\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if staged {
		if out, err := exec.Command("git", "-C", dir, "add", "a.txt").CombinedOutput(); err != nil {
			t.Fatalf("git add: %v: %s", err, out)
		}
	}
	return dir
}

func TestRunImprimeMensaje(t *testing.T) {
	dir := repoTemporal(t, true)
	f := &falso{resp: "feat: añade a.txt"}
	var out bytes.Buffer
	if err := run(context.Background(), &options{dir: dir}, f, &out); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != "feat: añade a.txt" {
		t.Errorf("salida: %q", out.String())
	}
	if !strings.Contains(f.diff, "a.txt") {
		t.Errorf("el proveedor no recibió el diff: %q", f.diff)
	}
}

func TestRunCommit(t *testing.T) {
	dir := repoTemporal(t, true)
	if err := run(context.Background(), &options{dir: dir, commit: true}, &falso{resp: "feat: añade a.txt"}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", dir, "log", "-1", "--format=%s").Output()
	if err != nil || strings.TrimSpace(string(out)) != "feat: añade a.txt" {
		t.Errorf("log: %q, err %v", out, err)
	}
}

func TestRunSinStaged(t *testing.T) {
	dir := repoTemporal(t, false)
	err := run(context.Background(), &options{dir: dir}, &falso{resp: "feat: x"}, &bytes.Buffer{})
	if !errors.Is(err, git.ErrNoStaged) {
		t.Fatalf("err = %v", err)
	}
}

func TestRunRespuestaInvalida(t *testing.T) {
	dir := repoTemporal(t, true)
	if err := run(context.Background(), &options{dir: dir}, &falso{resp: "hola mundo"}, &bytes.Buffer{}); err == nil {
		t.Fatal("esperaba error")
	}
}

func TestNewProvider(t *testing.T) {
	noEnv := func(string) string { return "" }
	if _, err := newProvider(&options{provider: "claude"}, noEnv); !errors.Is(err, provider.ErrNoAPIKey) {
		t.Errorf("claude sin key: %v", err)
	}
	if p, err := newProvider(&options{provider: "claude"}, func(string) string { return "k" }); err != nil || p == nil {
		t.Errorf("claude con key: %v", err)
	}
	if _, err := newProvider(&options{provider: "ollama"}, noEnv); err != nil {
		t.Errorf("ollama: %v", err)
	}
	if _, err := newProvider(&options{provider: "x"}, noEnv); err == nil {
		t.Error("esperaba error por proveedor desconocido")
	}
}
