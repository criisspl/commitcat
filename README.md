# commitcat 🐱

CLI en Go que lee `git diff --staged` y propone un mensaje en formato [Conventional Commits](https://www.conventionalcommits.org/es/).

Dos proveedores:

- **Ollama** (por defecto): modelo en local, sin enviar código fuera de tu máquina.
- **Claude** (opcional): API de Anthropic, solo si defines `ANTHROPIC_API_KEY`.

## Instalación

Con Go:

```sh
go install github.com/criisspl/commitcat@latest
```

O descarga un binario desde [Releases](https://github.com/criisspl/commitcat/releases).

Desde el código fuente:

```sh
git clone https://github.com/criisspl/commitcat && cd commitcat
go build -o commitcat .
```

## Uso

```sh
git add .
commitcat                 # imprime el mensaje propuesto
commitcat --commit        # además crea el commit
commitcat -m qwen2.5-coder   # otro modelo de Ollama
```

Con Claude:

```sh
export ANTHROPIC_API_KEY=...
commitcat --provider claude
```

> Con `--provider claude` el diff staged se envía a la API de Anthropic.

### Flags

| Flag | Descripción | Por defecto |
|---|---|---|
| `-p, --provider` | `ollama` o `claude` | `ollama` |
| `-m, --model` | modelo | `llama3.2` (Ollama) / `claude-haiku-5-5` (Claude) |
| `--ollama-url` | URL del servidor Ollama | `http://localhost:11434` |
| `-c, --commit` | crear el commit | `false` |
| `--dir` | directorio del repo | actual |

Variables: `ANTHROPIC_API_KEY` (clave), `ANTHROPIC_BASE_URL` (opcional, para proxies).

Antes de usar Ollama: instálalo, arranca `ollama serve` y descarga el modelo (`ollama pull llama3.2`).

## Desarrollo

```sh
go vet ./...
go test ./...   # usa un proveedor falso: no necesita red ni Ollama
```

Arquitectura: la interfaz `Provider` (`internal/provider`) tiene las implementaciones `Ollama` y `Claude`; `internal/message` construye el prompt y valida el formato; `internal/git` envuelve git.

Ramas `feature/<nombre>`, Conventional Commits y una PR por cambio. Releases con [goreleaser](https://goreleaser.com) al subir un tag `v*`.

## Licencia

Pendiente de definir.
