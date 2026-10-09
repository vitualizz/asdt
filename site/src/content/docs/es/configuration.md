---
title: Configuración
description: Referencia de .asdt/config.yaml, .asdt/knowledge/knowledge.yaml y todas las opciones de configuración.
order: 2
locale: es
---

# Configuración

## Configuración del proyecto — `.asdt/config.yaml`

Ejecutar `/asdt-init` (dentro de tu asistente) crea `.asdt/config.yaml`. Las keys que agregues a mano sobreviven a un re-init posterior.

```yaml
memory:
  provider: engram
strict_tdd: false              # opcional
code_intelligence: codegraph   # solo si se detectó
primary_design_surface: mobile # solo si se respondió
```

### `memory.provider`

El backend de memoria que usan todos los especialistas para persistir y recuperar hand-offs. Hoy solo se soporta `engram`.

### `strict_tdd`

Cuando es `true`, el paso `implement` del Developer escribe los tests en la misma pasada que el código, dentro de la misma lista de archivos aprobada. No hay un paso de tests aparte, e `implement` nunca los ejecuta: la compuerta `verify` que sigue te ofrece los comandos de chequeo (nunca uno que reescriba archivos) y los corre solo si decís que sí. También podés pedir tests en el pedido sin tocar el flag.

Valor por defecto: ausente, que equivale a `false`.

### `code_intelligence`

`/asdt-init` lo escribe solo cuando detecta un índice de code intelligence (hoy, `codegraph`). Los especialistas lo prefieren entonces a los ciclos de grep y lectura. Nunca se escribe como `none` — si no se detecta nada, la key se elimina.

### `primary_design_surface`

La superficie para la que UX/UI diseña primero: `mobile`, `tablet`, `desktop` o `none`. `/asdt-init` la pregunta una vez (por defecto `mobile`); un re-init te deja cambiarla.

- **Ausente** — la pregunta nunca se hizo (por ejemplo, una corrida no interactiva). Los especialistas que diseñan o construyen UI (UX/UI y el Developer) la tratan como `mobile` y registran una entrada `ASSUMED:` que lo dice.
- **`none`** — el proyecto no tiene ninguna superficie visual (un CLI, una librería, un servicio backend). UX/UI entonces se saltea su spec y guarda un hand-off breve diciendo que no hay nada que diseñar, así los especialistas siguientes leen un "sin UI" explícito en lugar de un hueco silencioso. El Developer no agrega capa responsive.

## Conocimiento del proyecto — `.asdt/knowledge/knowledge.yaml`

Generado por `/asdt-init` a partir de escaneos acotados del repo. Es el único archivo de conocimiento que leen los especialistas: el paso inline `platform-analysis` lo convierte en un resumen corto (stack, convenciones, huella de diseño, más la superficie de diseño de `config.yaml`) para que ningún especialista vuelva a detectar el stack.

```yaml
schema_version: "2"
scanned_at: 2026-05-04T12:00:00Z
stack: [typescript]
file_structure: "src/ for application code, tests/ alongside"
design_fingerprint:
  css_approach: tailwind
  ci_cd: github-actions
is_monorepo: { value: "false", source: "detected", confidence: "high" }
test_runner: { value: "vitest", source: "detected", confidence: "high" }
naming_style: { value: "kebab-case", source: "detected", confidence: "medium" }
architectural_style: { value: "layered", source: "manual", confidence: "high" }
# ASDT:NUANCE:BEGIN
human_nuance:
  - topic: "legacy-css"
    type: repo_practice
    note: "CSS legacy en styles/; lo nuevo usa utilidades de Tailwind"
    source: manual
    origin: user
# ASDT:NUANCE:END
```

Tratalo como propiedad de `/asdt-init`: volvé a correrlo para recalibrar en lugar de editar a mano los campos detectados. Un campo que fijás durante una revisión de recalibración queda marcado `source: manual` y nunca se sobrescribe en silencio.

`human_nuance` es la excepción — guarda lo que ningún escaneo puede detectar, en tus palabras. Se llena con las preguntas opcionales de `/asdt-init`, o cuando le pedís a un especialista que recuerde algo duradero ("recordá que…"). Si un especialista nota un hecho así por su cuenta, lo propone en una línea y lo guarda solo si decís que sí. Estas notas nunca se mezclan con el resumen detectado: llegan a un paso como notas del proyecto aparte y etiquetadas, y solo cuando son relevantes.

`/asdt-init` también escribe `.asdt/knowledge/provenance.yaml`, un archivo auxiliar que registra de dónde salió cada valor detectado. Solo lo lee un `/asdt-init` posterior.

## Variables de entorno

ASDT no lee ninguna variable de entorno directamente. El binario usa el entorno del asistente de IA (Claude Code u OpenCode) para toda la configuración en tiempo de ejecución.

## Múltiples proyectos

Cada proyecto tiene su propio directorio `.asdt/`. Cambiá de proyecto abriendo tu asistente de IA en una carpeta distinta — ASDT lee la configuración desde el `.asdt/config.yaml` más cercano, subiendo desde el directorio de trabajo actual.
