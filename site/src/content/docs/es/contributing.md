---
title: Contribuir
description: Cómo agregar un nuevo specialist, mejorar prompts, escribir shared skills y enviar un PR a ASDT.
order: 8
locale: es
---

# Contribuir

Las contribuciones de mayor impacto son los archivos `SKILL.md` de los specialists y las definiciones de los steps del workflow — no necesitas experiencia en Go. La capa de skills ES el producto. Si puedes describir el rol de un specialist, sus workflow steps y sus contratos de artifacts, puedes shippear un nuevo specialist.

## Agregar un nuevo specialist

### 1. Creá la estructura de directorios

```
skill/asdt-{name}/
  SKILL.md          # definición del specialist y su workflow
  workflow.yaml     # secuencia de steps y metadata
  steps/            # un .md por workflow step
  skills/           # fragmentos de skills específicos del specialist (opcional)
```

El nombre del directorio **debe** empezar con `asdt-`. El binario embebe el skill tree vía `//go:embed SKILL.md asdt-*` en `skill/embedded.go` — cualquier directorio que matchee `asdt-*` se shippea automáticamente en el próximo build.

### 2. Escribí SKILL.md

```markdown
---
name: asdt-{name}
description: "Una oración: qué produce este specialist."
user-invocable: true
specialist-id: {name}
shared-skills: # solo metadata — documenta skills relacionadas; no es un loader
  - platform-context
  - artifact-envelope
metadata:
  author: "Your Name"
  version: "1.0"
---

# {Name} Specialist

## Role
...

## Orchestration Plan
...

## Invariants
...
```

`metadata` (`author` + `version`) es obligatorio en todos los `SKILL.md`.

### 3. Escribí workflow.yaml

```yaml
specialist: {name}
steps:
  - id: step-one
    name: Step One
  - id: step-two
    name: Step Two
```

### 4. Escribí los archivos de steps

Creá un `.md` por step en `skill/{name}/steps/{step-id}.md`. Cada archivo contiene las instrucciones para el LLM en ese step — qué leer, qué producir, qué formato debe tener el artifact.

### 5. Conectalo al embed

Abrí `skill/embedded.go` y agregá el nombre de tu directorio a la directiva `//go:embed`:

```go
//go:embed SKILL.md asdt-shared asdt-developer asdt-ux-ui asdt-architect asdt-qa asdt-security asdt-init asdt-{name}
var skillFS embed.FS
```

**No te saltees este paso.** Un directorio que existe en disco pero no está nombrado en la directiva queda excluido del binario en silencio — sin error de compilación, sin error en runtime. La skill simplemente no aparece cuando los usuarios instalan.

### 6. Verificá con el sandbox

```sh
mkdir -p /tmp/asdt-sandbox
HOME=/tmp/asdt-sandbox go run ./cmd/asdt-tui
```

Instala en un directorio descartable. Confirmá que tu specialist aparece como su propio sibling de primer nivel bajo `/tmp/asdt-sandbox/.claude/skills/{name}/`.

### 7. Corré los embed tests

```sh
go test ./skill/...
```

`skill/embedded_test.go` verifica que todo directorio `asdt-*` en disco esté presente en el embedded FS y tenga un `SKILL.md`. Falla ruidosamente si tu specialist falta.

### 8. Actualizá este README

Agregá una fila a la tabla de specialists en el README del proyecto con el comando, el rol y qué produce.

## Mejorar el prompt de un specialist

1. Editá `skill/{specialist}/SKILL.md` o cualquier archivo bajo `skill/{specialist}/steps/` o `skill/{specialist}/skills/`.
2. Corré `go test ./skill/...` para confirmar que el embed registry recoge los cambios.
3. Abrí un PR. Los PRs que solo tocan prompts son contribuciones de primera clase.

## Agregar una shared skill

Las shared skills son fragmentos de capacidad reutilizados entre múltiples specialists — detección de platform context, formato de artifact envelope, definición de scope.

1. Creá `skill/asdt-shared/skills/{name}.md` con las instrucciones de la capacidad.
2. Documentala en el frontmatter `shared-skills` de cualquier specialist que se relacione con ella — esa key es solo metadata; la carga real ocurre vía el FIRST ACTION Read en el body del `SKILL.md` (`specialist-header`) y los `reference_skills` por step en `workflow.yaml`.
3. Abrí un PR.

## Estándares de código

- Early return: `if err != nil { return err }` — validá los inputs primero.
- Sin estado global — constructor injection en todos lados.
- Interfaces definidas cerca de sus consumidores, no en el paquete que las implementa.
- Nada de paquetes `utils/`, `helpers/`, `common/` o `misc/` — solo domain nouns.
- Table-driven tests para toda lógica con más de dos casos.

## Proceso de PR

- Un cambio lógico por PR.
- `go test ./...` debe pasar.
