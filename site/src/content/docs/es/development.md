---
title: Desarrollo
description: Cómo compilar, testear y correr ASDT localmente — incluyendo la trampa del cableado de embed y el flujo de testing en sandbox.
order: 9
locale: es
---

# Desarrollo

Flujos de trabajo prácticos del día a día para trabajar en ASDT mismo: compilar, lintear, agregar skills, y probar el TUI del instalador sin tocar la configuración real de tu asistente de IA.

## Prerrequisitos

- Go (versión fijada en `go.mod`)
- golangci-lint **v2** — instalá específicamente la ruta del módulo v2:

```sh
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

> La configuración de `.golangci.yml` declara `version: "2"`. Instalar la ruta del módulo v1 te da un binario v1 que rechaza la config.

- lefthook (opcional) — `make hooks` lo instala y registra el hook de lint pre-push

## Comandos del día a día

```sh
make build   # go build ./...
make test    # go test ./...
make lint    # golangci-lint run ./...
make hooks   # instala lefthook + registra el hook pre-push
```

Corré el TUI del instalador desde el código fuente sin un paso de build:

```sh
go run ./cmd/asdt-tui
```

## Testing en un sandbox

`asdt-tui` escribe directamente en el directorio real de skills de tu asistente de IA (`~/.claude/skills`, `~/.config/opencode/skills`). Para evitar tocar tu configuración real, sobreescribí `$HOME`:

```sh
mkdir -p /tmp/asdt-sandbox
HOME=/tmp/asdt-sandbox go run ./cmd/asdt-tui
```

Esto instala en `/tmp/asdt-sandbox/.claude/skills` y `/tmp/asdt-sandbox/.config/opencode/skills`. Tu `~/.claude` real nunca se toca.

Inspeccioná el resultado:

```sh
eza --tree /tmp/asdt-sandbox/.claude/skills
```

Volvé a correrlo contra el mismo `$HOME` para probar el camino de detección de "ya instalado / actualizar". Limpiá cuando termines:

```sh
rm -rf /tmp/asdt-sandbox
```

## Agregar un nuevo especialista — checklist de cableado

Los detalles completos de autoría están en [Contribuir](/asdt/docs/contributing). El paso que la guía de contribución no destaca:

**Cablealo en el embed.** Abrí `skill/embedded.go`:

```go
//go:embed SKILL.md asdt-shared asdt-developer asdt-ux-ui asdt-architect asdt-qa asdt-security asdt-init
var skillFS embed.FS
```

Agregá el nombre de tu directorio a esta lista. Un directorio de skill que existe en el disco pero **no** está listado acá queda excluido silenciosamente del binario — `go build` funciona bien, el TUI corre sin problemas, y la skill nunca aparece en el resultado instalado. No hay ninguna advertencia en tiempo de compilación ni en tiempo de ejecución.

Después de cablearlo:

1. Corré el flujo de sandbox para confirmar que la skill aparece en `/tmp/asdt-sandbox/.claude/skills/{name}/`
2. Corré `go test ./skill/...` — el test del registro embebido verifica que la skill esté presente

## Verificar ediciones de prompts

Editar el `SKILL.md` o los `steps/*.md` de un especialista existente no requiere cambios en la lista de embed. Simplemente:

```sh
go test ./internal/prompt/...
```

`go:embed` vuelve a leer los archivos en tiempo de compilación, así que `go test` y `go run` siempre reflejan tus últimas ediciones — sin caché de la cual preocuparse.

## Estructura del proyecto

```
cmd/asdt-tui/       # punto de entrada del TUI del instalador
internal/
  installer/        # detección de skills, instalación, lógica de actualización
  setup/            # vistas del TUI y máquina de estados
  prompt/           # ensamblado y embebido de prompts
  i18n/             # catálogo de strings del TUI (en + es)
skill/
  embedded.go       # registro go:embed
  asdt-shared/      # fragmentos de skill compartidos
  asdt-architect/   # especialista Architect
  asdt-developer/   # especialista Developer
  asdt-qa/          # especialista QA
  asdt-security/    # especialista Security
  asdt-ux-ui/       # especialista UX/UI
  asdt-init/        # especialista de inicialización de proyecto
site/               # Este sitio de documentación (Astro)
docs/               # ADRs y guías para contribuidores
```
