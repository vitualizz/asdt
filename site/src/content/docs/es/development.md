---
title: Desarrollo
description: Cómo compilar, testear y correr ASDT localmente — incluyendo el checklist de registro de especialistas y el flujo de testing en sandbox.
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

## Agregar un nuevo especialista — checklist de registro

Los detalles completos de autoría están en [Contribuir](/asdt/docs/contributing). Lo que conviene saber de entrada:

**El embed no necesita cableado.** `skill/embedded.go` usa un glob:

```go
//go:embed SKILL.md asdt-*
var skillFS embed.FS
```

Cualquier directorio llamado `asdt-{name}` se shippea automáticamente en el próximo build — no hay lista a la cual agregarlo ni forma de olvidarse. Lo que sí hay que hacer a mano es el registro: la fila de la tabla `## Registry` en `skill/SKILL.md`, la fila de `## ASDT Specialists` en `internal/installer/assets/agents-template.md`, y la lista de specialists hardcodeada en `skill/embedded_test.go`. Si te salteás eso, la skill igual se shippea, solo que nunca se enruta.

Después de agregar el directorio:

1. Corré el flujo de sandbox para confirmar que la skill aparece en `/tmp/asdt-sandbox/.claude/skills/{name}/`
2. Corré `go test ./skill/...` — el test del registro embebido verifica que la skill esté presente

## Agregar un memory provider

Los prompts nunca nombran un memory provider. Hablan tres verbos — **save** (upsert por clave), **search** (por prefijo de clave o texto libre) y **get** (un registro completo por id) — definidos una sola vez en `skill/asdt-core/protocol.md` §0. Cada archivo que llega a la memoria lleva una región `<!-- ASDT:GENERATED:memory-binding -->` commiteada vacía: `protocol.md` (que se empalma, con ella, en cada `SKILL.md` enrutado), `executor-header.md`, `asdt-init/SKILL.md` y el `SKILL.md` raíz. Al instalar, el instalador renderiza el binding del provider elegido en esas regiones y otorga sus herramientas a `asdt-analyst` y `asdt-builder`.

Un provider nuevo es un cambio solo en Go — sin editar prompts:

1. Agregá un `ProviderDescriptor` a `Providers` en `internal/installer/providers.go`:
   - `ID`, `Name`, `Description` — la pantalla de provider de la TUI lo lista a partir de estos.
   - `ConfigValue` — lo que `/asdt-init` escribe en `memory.provider`, y lo que cada especialista verifica antes de correr.
   - `Tools` — un `MemoryTool` por verbo: el nombre simple de la herramienta y un `Usage` de una línea que le dice al modelo cómo llamarla (por ejemplo, qué parámetro de la herramienta lleva la clave). `Usage` se empalma tal cual en el binding. El `Usage` de search además tiene que decir dónde muestra cada resultado su clave — el protocolo hace que un paso elija un resultado por clave exacta, así que necesita saber qué campo la lleva. El de Engram, por ejemplo, dice que un resultado muestra la clave como su `title`, sin el `{project}/` inicial.
   - `ClaudeToolPrefixes` — los namespaces MCP bajo los que Claude Code puede exponer las herramientas; al menos uno. La allowlist de Claude Code de los agentes ejecutores es cada prefijo × la herramienta de cada verbo, derivada de `Tools`, nunca una segunda lista.
   - `Detect` — cómo lo busca el chequeo de entorno de la TUI: `Binary`, que se busca en tu `PATH`, e `InstallURL`, adonde se envía a quien no lo tiene. Dejá `Binary` vacío para un provider sin nada que chequear localmente: su fila es informativa y nunca bloquea. El chequeo lista todos los providers; *Continue* necesita al menos uno disponible, y la pantalla de provider solo avanza con uno disponible.
2. Corré `go test ./...`. `TestProviders_BindEveryVerb` falla ante cualquier hueco que encuentre el `Validate` del descriptor — un verbo sin vincular, un `Usage` faltante, ningún prefijo; `TestInstall_MemoryBindingReachesInstalledFiles` muestra lo que renderiza un binding; `TestPromptsAreProviderNeutral` falla si algún prompt, agente generado o persona nombra un provider incluido o sus herramientas — los nombres prohibidos salen de `Providers`, así que los tuyos quedan cubiertos apenas lo agregás.

## Verificar ediciones de prompts

Editar el `SKILL.md` o los `steps/*.md` de un especialista existente no requiere ningún cambio de cableado. Simplemente:

```sh
go test ./skill/... ./internal/installer/...
```

`go:embed` vuelve a leer los archivos en tiempo de compilación, así que `go test` y `go run` siempre reflejan tus últimas ediciones — sin caché de la cual preocuparse.

## Estructura del proyecto

```
cmd/asdt-tui/       # punto de entrada del TUI del instalador
internal/
  grader/           # califica payloads de artefactos contra probe sets
  i18n/             # catálogo de strings del TUI (en + es)
  installer/        # detección de skills, instalación, lógica de actualización
  setup/            # vistas del TUI y máquina de estados
  tui/panels/       # primitivas puras de render del TUI (badge, hero, spinner)
skill/
  SKILL.md          # orquestador raíz + registro de especialistas
  TEMPLATE.md       # contrato de autoría de especialistas (solo repo)
  README.md         # panorama de la capa de skills
  embedded.go       # registro go:embed
  embedded_test.go  # tests de invariantes routed + embed
  asdt-core/        # protocolo, headers y referencias opcionales
  asdt-architect/   # especialista Architect
  asdt-developer/   # especialista Developer
  asdt-pm/          # especialista Product Manager
  asdt-qa/          # especialista QA
  asdt-researcher/  # especialista Researcher
  asdt-security/    # especialista Security
  asdt-ux-ui/       # especialista UX/UI
  asdt-init/        # especialista de inicialización de proyecto (clase setup)
site/               # Este sitio de documentación (Astro)
```
