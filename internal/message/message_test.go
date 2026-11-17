package message

import (
	"strings"
	"testing"
)

func TestNormalize(t *testing.T) {
	casos := []struct {
		nombre, in, want string
		wantErr          bool
	}{
		{"simple", "feat: añade flag --commit", "feat: añade flag --commit", false},
		{"ámbito", "fix(git): corrige diff vacío\n\ncuerpo", "fix(git): corrige diff vacío\n\ncuerpo", false},
		{"bloque de código", "```\nchore: actualiza deps\n```", "chore: actualiza deps", false},
		{"comillas", "\"docs: mejora readme\"", "docs: mejora readme", false},
		{"breaking", "feat(api)!: cambia firma", "feat(api)!: cambia firma", false},
		{"tipo inválido", "update: algo", "", true},
		{"sin descripción", "feat:", "", true},
		{"vacío", "  ", "", true},
		{"prosa", "Aquí tienes el mensaje: feat: x", "", true},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, err := Normalize(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, c.wantErr)
			}
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestTruncateDiff(t *testing.T) {
	if got := TruncateDiff("corto"); got != "corto" {
		t.Errorf("no debería truncar: %q", got)
	}
	largo := strings.Repeat("a", MaxDiffChars+10)
	got := TruncateDiff(largo)
	if !strings.HasSuffix(got, "[... diff truncado ...]") || len(got) > MaxDiffChars+50 {
		t.Errorf("truncado incorrecto (len %d)", len(got))
	}
}
