// Package git envuelve las llamadas a git que necesita commitcat.
package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrNoStaged indica que no hay cambios en el área de staging.
var ErrNoStaged = errors.New("no hay cambios staged: usa `git add` primero")

// StagedDiff devuelve `git diff --staged` en el directorio dir ("" = actual).
func StagedDiff(dir string) (string, error) {
	out, err := run(dir, "diff", "--staged", "--no-color")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(out) == "" {
		return "", ErrNoStaged
	}
	return out, nil
}

// Commit crea un commit con el mensaje dado.
func Commit(dir, msg string) error {
	_, err := run(dir, "commit", "-m", msg)
	return err
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}
