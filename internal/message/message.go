// Package message construye el prompt y valida el Conventional Commit generado.
package message

import (
	"fmt"
	"regexp"
	"strings"
)

// MaxDiffChars limita el diff enviado al modelo.
const MaxDiffChars = 20000

// SystemPrompt instruye al modelo.
const SystemPrompt = `Eres un asistente que escribe mensajes de commit en formato Conventional Commits.
Recibes un "git diff --staged". Responde SOLO con el mensaje, sin explicaciones ni comillas ni bloques de código.
Formato: <tipo>(<ámbito opcional>): <descripción>
Tipos válidos: feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert.
La primera línea tiene máximo 72 caracteres, en imperativo y minúsculas. Añade un cuerpo solo si aporta contexto.`

var header = regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([^)\n]+\))?!?: \S.*$`)

// TruncateDiff recorta el diff a MaxDiffChars.
func TruncateDiff(diff string) string {
	if len(diff) <= MaxDiffChars {
		return diff
	}
	return diff[:MaxDiffChars] + "\n[... diff truncado ...]"
}

// Normalize limpia la respuesta del modelo y comprueba que es un Conventional Commit.
func Normalize(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		if i := strings.Index(s, "\n"); i >= 0 {
			s = s[i+1:]
		}
		s = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(s), "```"))
	}
	s = strings.Trim(s, "`\"'")
	s = strings.TrimSpace(s)
	first, _, _ := strings.Cut(s, "\n")
	if !header.MatchString(first) {
		return "", fmt.Errorf("la respuesta no es un Conventional Commit: %q", first)
	}
	return s, nil
}
