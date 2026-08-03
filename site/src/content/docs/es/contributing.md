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

`metadata` (`author` + `version`) es obligatorio en todos los `SKILL.md`. `trigger_phrases` es la única key opcional — una lista de descubribilidad que consume el host.

`shared-skills` está **prohibida**. La key está retirada: ningún loader la resolvió nunca, así que declararla no documenta nada ni carga nada. Ver [Cómo se cargan realmente las shared skills](#cómo-se-cargan-realmente-las-shared-skills) para los tres mecanismos reales.

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

### 5. Registrá el specialist

El embed no necesita nada de vos — `//go:embed SKILL.md asdt-*` ya shippea tu directorio. Lo que **no** es automático es el registro: un specialist routable tiene que espejarse a mano en tres lugares, más un fixture de test.

1. `skill/SKILL.md` — agregá la fila de `Specialist Registry` (comando, disciplina, cuándo involucrarlo).
2. `skill/SKILL.md` — agregá la fila de la tabla por specialist en `Tailored Workflow Generation`. Nunca edites a mano dentro de los marcadores generados de inline steps de esa tabla; esa subregión se regenera al instalar y tus ediciones quedan sobrescritas.
3. `internal/installer/assets/agents-template.md` — agregá la fila de ASDT Specialists.
4. `skill/embedded_test.go` — el test de invariantes routed mantiene una lista de specialists hardcodeada que un mantenedor debe actualizar.

**No te saltees estos pasos.** El directorio se shippea igual, así que nada falla en tiempo de build — el specialist simplemente nunca aparece en el routing, ni en el archivo de agents instalado, ni en la cobertura del test de invariantes.

`/asdt-init` es la excepción: es un specialist de clase setup, deliberadamente no routable y deliberadamente ausente de las tablas de routing. No "arregles" esa omisión.

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

## Mejorar el prompt de un specialist

1. Editá `skill/{specialist}/SKILL.md` o cualquier archivo bajo `skill/{specialist}/steps/` o `skill/{specialist}/skills/`.
2. Corré `go test ./skill/...` para confirmar que el embed registry recoge los cambios.
3. Abrí un PR. Los PRs que solo tocan prompts son contribuciones de primera clase.

## Agregar una shared skill

Las shared skills son fragmentos de capacidad reutilizados entre múltiples specialists — detección de platform context, knowledge recall, definición de scope.

1. Creá `skill/asdt-core/references/{name}.md` con las instrucciones de la capacidad.
2. Conectala a través de uno de los tres mecanismos de carga de abajo. Una shared skill que nadie declara nunca se lee — no hay carga implícita ni ambiental.
3. Abrí un PR.

### Cómo se cargan realmente las shared skills

Tres mecanismos, y solo tres. Los paths siempre se resuelven desde el directorio propio del specialist.

**1. Splice en tiempo de instalación.** El instalador injerta `asdt-core/specialist-header.md` en una región generada de cada `SKILL.md` routado, así el orquestador lee el header inline en lugar de ir a buscar un archivo aparte. Esto aplica solo a ese archivo. El blockquote de FIRST ACTION ya no indica leer `specialist-header.md` — el único archivo al que te manda es `./workflow.yaml`. Nunca edites a mano entre los marcadores de la región; el splice sobrescribe lo que haya ahí.

**2. Inline step.** Un step de `workflow.yaml` con `execution: inline` cuyo `skill:` nombra un archivo compartido — `knowledge-recall.md`, `platform-context.md` (declarado como el step `platform-analysis`), `decision-preservation.md`. El orquestador lee ese archivo y lo sigue en su propio contexto. No se inyecta nada en ningún lado y no se lanza ningún sub-agente.

**3. `reference_skills:` en un step `subagent`.** Antes de lanzar el step, el orquestador lee cada archivo listado e inyecta su contenido en el prompt del sub-agente como un bloque `### REFERENCE SKILL {path}`. El sub-agente nunca los busca por su cuenta — los sub-agentes corren desde un directorio de trabajo distinto y no pueden resolver esos paths. Cuando una lectura falla, el bloque llega como `### REFERENCE SKILL {path}: UNRESOLVED` y el step sigue en modo best-effort.

## Estándares de código

- Early return: `if err != nil { return err }` — validá los inputs primero.
- Sin estado global — constructor injection en todos lados.
- Interfaces definidas cerca de sus consumidores, no en el paquete que las implementa.
- Nada de paquetes `utils/`, `helpers/`, `common/` o `misc/` — solo domain nouns.
- Table-driven tests para toda lógica con más de dos casos.

## Proceso de PR

- Un cambio lógico por PR.
- `go test ./...` debe pasar.
