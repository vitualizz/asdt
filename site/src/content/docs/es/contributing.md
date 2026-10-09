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
  steps/            # un .md por step subagent
```

No hay directorio `skills/`. Los criterios compartidos viven en `skill/asdt-core/references/` y se declaran por step con `reference_skills:`.

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

## Final Output
...

## Invariants
...
```

Entre el frontmatter y `## Role` van tres bloques textuales de los que depende el instalador — el blockquote `FIRST ACTION`, la región generada vacía (`<!-- ASDT:GENERATED:specialist-header -->` … `<!-- /ASDT:GENERATED:specialist-header -->`) y el blockquote `ORCHESTRATOR GATE`. Copialos de cualquier specialist existente sin tocarlos. Un step `inline` sin archivo de prompt — una compuerta donde el orquestador se detiene a esperar al humano, como `approve` y `verify` del Developer — lleva su propia sección `## {step} — …` antes de `## Final Output`; esa sección es todo su contrato.

`metadata` (`author` + `version`) es obligatorio en todos los `SKILL.md`. `trigger_phrases` es la única key opcional — una lista de descubribilidad que consume el host.

`shared-skills` está **prohibida**. La key está retirada: ningún loader la resolvió nunca, así que declararla no documenta nada ni carga nada. Ver [Cómo se cargan realmente las shared skills](#cómo-se-cargan-realmente-las-shared-skills) para los tres mecanismos reales.

### 3. Escribí workflow.yaml

```yaml
specialist: {name}
routable: true
steps:
  - name: {step}
    skill: steps/{step}.md
    description: Una línea — qué produce este step.
    execution: subagent          # o: inline
    model: haiku | sonnet | opus
    agent: analyst | builder     # builder SOLO cuando el step escribe archivos del host
    inputs:
      - "{project}/{change}/{role}/handoff"  # opcional — decilo, y degradá
    output_topic_key: "{project}/{change}/{name}/handoff"
    reference_skills:
      - ../asdt-core/references/{x}.md
      - ../asdt-core/protocol.md
    host_skills:                 # opcional — solo en steps subagent
      - {nombre-de-skill-en-kebab-case}
```

- Un step cuyo payload solo alimenta al siguiente step de la misma corrida declara `output: context` en lugar de `output_topic_key` — no persiste nada.
- Un step `inline` de compuerta sin archivo de prompt declara solo `name`, `description` y `execution: inline`.
- **Techo: cuatro steps `subagent`.** Los preludios y compuertas inline no cuentan. Si necesitás un quinto, el diseño está mal — fusioná dos o dividí el specialist.
- Todo input entre specialists es opcional y degrada — ningún specialist puede exigir el hand-off de otro.

`host_skills:` nombra skills que el **asistente host** puede tener instaladas — un oficio que ASDT no trae, como `frontend-design`. Son nombres, nunca paths, y nunca son obligatorias: si el host tiene la skill (o un equivalente obvio), el orquestador lee el archivo de la skill e inyecta su texto en el prompt del sub-agente como un bloque `### HOST SKILL {name}` — nunca la activa en su propio contexto; si no la tiene, se omite en silencio y no queda registro. Declarala solo en un step `subagent` cuyo archivo diga cuándo aplica y qué tiene prioridad sobre ella — una host skill afina la ejecución dentro del contrato del step, nunca lo reemplaza. Hoy declaran `frontend-design` el `ui-design` de UX/UI y el `implement` del Developer.

### 4. Escribí los archivos de steps

Creá un `.md` por step `subagent` en `skill/{name}/steps/{step}.md`. Cada archivo contiene las instrucciones para el LLM en ese step — qué leer, qué producir, qué formato debe tener el artifact.

### 5. Registrá el specialist

El embed no necesita nada de vos — `//go:embed SKILL.md asdt-*` ya shippea tu directorio. Lo que **no** es automático es el registro: un specialist routable tiene que espejarse a mano en dos lugares, más un fixture de test.

1. `skill/SKILL.md` — agregá la fila a la tabla de `## Registry` (comando, disciplina, cuándo involucrarlo).
2. `internal/installer/assets/agents-template.md` — agregá la fila a la tabla de `## ASDT Specialists`.
3. `skill/embedded_test.go` — el test de invariantes routed mantiene una lista de specialists hardcodeada que un mantenedor debe actualizar.

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

1. Editá `skill/{specialist}/SKILL.md` o cualquier archivo bajo `skill/{specialist}/steps/`.
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

**2. Inline step.** Un step de `workflow.yaml` con `execution: inline` cuyo `skill:` nombra un archivo compartido — `knowledge-recall.md`, `platform-context.md` (declarado como el step `platform-analysis`). El orquestador lee ese archivo y lo sigue en su propio contexto. No se inyecta nada en ningún lado y no se lanza ningún sub-agente.

**3. `reference_skills:` en un step `subagent`.** Antes de lanzar el step, el orquestador lee cada archivo listado e inyecta su contenido en el prompt del sub-agente como un bloque `### REFERENCE SKILL {path}`. El sub-agente nunca los busca por su cuenta — los sub-agentes corren desde un directorio de trabajo distinto y no pueden resolver esos paths. Cuando una lectura falla, el bloque llega como `### REFERENCE SKILL {path}: UNRESOLVED` y el step sigue en modo best-effort.

**No es una shared skill: `host_skills:`.** Una skill que el asistente host tiene instalada (p. ej. `frontend-design`) llega al sub-agente por la misma vía — el orquestador la inyecta como un bloque `### HOST SKILL {name}` — pero vive fuera de ASDT, es opcional y su ausencia no cambia nada.

## Estándares de código

- Early return: `if err != nil { return err }` — validá los inputs primero.
- Sin estado global — constructor injection en todos lados.
- Interfaces definidas cerca de sus consumidores, no en el paquete que las implementa.
- Nada de paquetes `utils/`, `helpers/`, `common/` o `misc/` — solo domain nouns.
- Table-driven tests para toda lógica con más de dos casos.

## Proceso de PR

- Un cambio lógico por PR.
- `go test ./...` debe pasar.
