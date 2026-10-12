---
title: Contribuir
description: Cómo agregar un nuevo specialist, mejorar prompts, escribir shared skills y enviar un PR a ASDT.
order: 8
locale: es
---

# Contribuir

Las contribuciones de mayor impacto son los archivos `SKILL.md` de los specialists y las definiciones de los steps del workflow — no necesitas experiencia en Go. La capa de skills ES el producto. Si puedes describir el rol de un specialist, sus workflow steps y el resumen que entrega, puedes shippear un nuevo specialist.

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
specialist: {name}               # el id del rol: nombra el archivo de resumen y los inputs propios
routable: true
steps:
  - name: knowledge-recall       # preludio inline: nombra la referencia que sigue
    skill: ../asdt-core/references/knowledge-recall.md
    description: Una línea.
    execution: inline

  - name: {step}
    skill: steps/{step}.md
    description: Una línea — qué devuelve este step.
    execution: subagent          # o: inline
    model: haiku | sonnet | opus
    agent: analyst | builder     # builder SOLO cuando el step escribe archivos del host
    inputs:                      # opcional — solo en steps subagent
      - pm # opcional — qué extraer del resumen de Product
      - {name}/{step-anterior} # el payload de un step propio de este specialist
    output: summary              # opcional — lo que devuelve este step pasa a ser el resumen del rol
    reference_skills:
      - ../asdt-core/references/{x}.md
      - ../asdt-core/protocol.md
    host_skills:                 # opcional — solo en steps subagent
      - {nombre-de-skill-en-kebab-case}
```

- **`inputs:`** son nombres del scratchpad, nunca paths ni claves de memoria. El orquestador resuelve cada uno a un archivo del scratchpad de la sesión y lo inyecta como `### INPUT {name}: {path}`, o `UNRESOLVED` si el archivo no existe. Una entrada es `{role}` (el resumen de otro specialist para este cambio, `{scratchpad}/asdt/{change}/{role}.md`, donde `{role}` es el id `specialist:` de ese specialist) o `{name}/{step}` (lo que devolvió un step anterior de ESTE specialist, `{scratchpad}/asdt/{change}/{name}/{step}.md`; listarlo es lo que hace que el orquestador escriba ese archivo). Nunca nombres el step de otro specialist: leé su resumen. `feedback` y `previous` están reservados para el orquestador, que los inyecta cuando tu corrección lo obliga a relanzar un step.
- **`output: summary`** significa que lo que devuelve el step se escribe como el resumen del rol, `{scratchpad}/asdt/{change}/{name}.md`. Un step posterior con `output: summary` en la misma corrida lo reemplaza: así un rol por etapas (plan, después construcción) mantiene UN solo resumen con su `Stage`. Un step sin `output: summary` le pasa lo que devuelve solo a los steps que lo listan.
- Un step `inline` de compuerta sin archivo de prompt declara solo `name`, `description` y `execution: inline`, más `output: summary` cuando el propio orquestador reescribe el resumen (el `verify` del Developer). Ningún step inline declara `inputs:`.
- **Techo: cuatro steps `subagent`.** Los preludios y compuertas inline no cuentan. Si necesitás un quinto, el diseño está mal — fusioná dos o dividí el specialist.
- Todo input entre specialists es opcional y degrada — marcalo con un comentario final `# optional — {qué extraer}`. Ningún specialist puede exigir el resumen de otro.
- Ningún step llega a la memoria. El conocimiento duradero del proyecto lo recuerda solo el orquestador, y ningún specialist inventa un segundo almacén para convenciones.

`host_skills:` nombra skills que el **asistente host** puede tener instaladas — un oficio que ASDT no trae, como `frontend-design`. Son nombres, nunca paths, y nunca son obligatorias: si el host tiene la skill (o un equivalente obvio), el orquestador lee el archivo de la skill e inyecta su texto en el prompt del sub-agente como un bloque `### HOST SKILL {name}` — nunca la activa en su propio contexto; si no la tiene, se omite en silencio y no queda registro. Declarala solo en un step `subagent` cuyo archivo diga cuándo aplica y qué tiene prioridad sobre ella — una host skill afina la ejecución dentro del contrato del step, nunca lo reemplaza. Hoy declaran `frontend-design` el `ui-design` de UX/UI y el `implement` del Developer.

### 4. Escribí los archivos de steps

Creá un `.md` por step `subagent` en `skill/{name}/steps/{step}.md`. Cada archivo tiene cuatro secciones — `Purpose`, `Inputs` (una viñeta por cada línea `### INPUT {name}`, cada una un path a leer, con una línea de degradación por input opcional), `Processing` y `Output` (el resumen que devuelve, en el esquema del protocolo compartido). Un step nunca llega a la memoria, nunca escribe un archivo (el `implement` del Developer y el `write` de `/asdt-init` son las únicas excepciones) y nunca le pregunta al usuario: cada hueco se vuelve una entrada `ASSUMED:`. Nunca escribas el nombre de un memory provider ni de sus herramientas en un prompt — `TestPromptsAreProviderNeutral` falla si aparece uno.

### 5. Registrá el specialist

El embed no necesita nada de vos — `//go:embed SKILL.md asdt-*` ya shippea tu directorio. Lo que **no** es automático es el registro: un specialist routable tiene que espejarse a mano en varios lugares.

1. `skill/SKILL.md` — agregá la fila a la tabla de `## Registry` (comando, disciplina, cuándo involucrarlo) y ubicalo en la prosa de dependencias debajo de esa tabla.
2. `internal/installer/assets/agents-template.md` — agregá la fila a la tabla de `## ASDT Specialists`.
3. `internal/installer/registry_gen.go` — agregá el directorio a `registryRenderOrder` y a `inlineStepsDisplayNames`. `TestRegistryDrift` verifica que esos espejos coincidan con los `workflow.yaml`.
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

**2. Inline step.** Un step de `workflow.yaml` con `execution: inline` cuyo `skill:` nombra un archivo compartido — `knowledge-recall.md`, `platform-context.md` (declarado como el step `platform-analysis`). El orquestador lee ese archivo y lo sigue en su propio contexto; no se lanza ningún sub-agente. Lo que produce el preludio — el bloque de plataforma, las notas del proyecto recordadas — le llega a cada step posterior como líneas en su prompt.

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
