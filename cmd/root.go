// Package cmd define la CLI de commitcat.
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/criisspl/commitcat/internal/git"
	"github.com/criisspl/commitcat/internal/message"
	"github.com/criisspl/commitcat/internal/provider"
)

type options struct {
	provider  string
	model     string
	ollamaURL string
	commit    bool
	dir       string
}

// Execute ejecuta el comando raíz.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	o := &options{}
	c := &cobra.Command{
		Use:           "commitcat",
		Short:         "Propone un mensaje Conventional Commit a partir de git diff --staged",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			p, err := newProvider(o, os.Getenv)
			if err != nil {
				return err
			}
			return run(c.Context(), o, p, c.OutOrStdout())
		},
	}
	f := c.Flags()
	f.StringVarP(&o.provider, "provider", "p", provider.NameOllama, "proveedor: ollama | claude (claude requiere ANTHROPIC_API_KEY)")
	f.StringVarP(&o.model, "model", "m", "", "modelo (por defecto el del proveedor)")
	f.StringVar(&o.ollamaURL, "ollama-url", provider.DefaultOllamaURL, "URL del servidor Ollama")
	f.BoolVarP(&o.commit, "commit", "c", false, "crear el commit con el mensaje propuesto")
	f.StringVar(&o.dir, "dir", "", "directorio del repo (por defecto el actual)")
	return c
}

func newProvider(o *options, getenv func(string) string) (provider.Provider, error) {
	switch o.provider {
	case provider.NameOllama:
		return &provider.Ollama{BaseURL: o.ollamaURL, Model: o.model}, nil
	case provider.NameClaude:
		key := getenv("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, provider.ErrNoAPIKey
		}
		return &provider.Claude{APIKey: key, BaseURL: getenv("ANTHROPIC_BASE_URL"), Model: o.model}, nil
	default:
		return nil, fmt.Errorf("proveedor desconocido %q (usa ollama o claude)", o.provider)
	}
}

func run(ctx context.Context, o *options, p provider.Provider, out io.Writer) error {
	diff, err := git.StagedDiff(o.dir)
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	raw, err := p.Generate(ctx, message.SystemPrompt, message.TruncateDiff(diff))
	if err != nil {
		return err
	}
	msg, err := message.Normalize(raw)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, msg)
	if o.commit {
		return git.Commit(o.dir, msg)
	}
	return nil
}
