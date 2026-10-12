---
title: Configuración
description: Referencia de .asdt/config.yaml, el único archivo que guarda los hechos de tu proyecto.
order: 2
locale: es
---

# Configuración

## Configuración del proyecto — `.asdt/config.yaml`

`.asdt/config.yaml` es el único lugar donde viven los hechos de tu proyecto. Ejecutar `/asdt-init` (dentro de tu asistente) los detecta con escaneos acotados del repo, resuelve con vos los ambiguos y escribe este archivo. Todos los especialistas lo leen: el paso inline `platform-analysis` lo convierte en un bloque de plataforma corto (stack, convenciones, huella de diseño, superficie de diseño), así ningún especialista vuelve a detectar el stack.

```yaml
memory:
  provider: engram
strict_tdd: false              # opcional — /asdt-init nunca lo agrega
code_intelligence: codegraph   # solo si se detecta
primary_design_surface: mobile # solo si se respondió
stack: [typescript]
file_structure: "src/ for application code, tests/ alongside"
is_monorepo: { value: "false", source: detected, confidence: high }
test_runner: { value: "vitest", source: detected, confidence: high }
naming_style: { value: "kebab-case", source: detected, confidence: medium }
architectural_style: { value: "layered", source: manual, confidence: high }
design_fingerprint:
  css_approach: { value: "tailwind", source: detected, confidence: high }
  ci_cd: { value: "github-actions", source: detected, confidence: high }
```

Tratá los campos detectados como propiedad de `/asdt-init`: volvé a ejecutarlo para recalibrar en vez de editarlos a mano. Una recalibración te muestra una tabla con cada valor que cambiaría, más cada valor que ya está marcado `source: manual`, y te pregunta una vez si aceptar todo o revisar key por key. Un valor que fijes distinto del ofrecido queda como `source: manual` y nunca se sobrescribe en silencio. Reconfirmar un valor nunca lo marca como manual: elegir el valor detectado guarda la entrada detectada, y aceptar un valor que habías fijado a mano lo deja como está. Para una key que fijaste a mano antes, el valor por defecto que ofrece `/asdt-init` es ese valor, no el detectado. Las keys que agregues a mano y que `/asdt-init` no reconoce se conservan intactas, y un re-init que no cambia nada deja el archivo idéntico byte a byte.

### `memory.provider`

El memory provider al que está vinculada tu instalación. Hoy solo se soporta `engram`. `/asdt-init` escribe el valor del provider de tu instalación, incluso sobre un valor existente. Cuando no hay ningún provider vinculado (una instalación **Sin memoria**), no escribe ningún bloque `memory` y deja como está el que exista. Si el proyecto nombra un provider distinto al de tu instalación, los especialistas tratan la memoria como no disponible: igual corren, trabajando desde el scratchpad de la sesión, y no se consulta ni se recuerda nada hasta que vuelvas a ejecutar `/asdt-init`. Ver [Scratchpad y memoria](/asdt/es/docs/memory-and-engram).

### `strict_tdd`

Cuando es `true`, el paso `implement` del Developer escribe los tests en la misma pasada que el código, dentro de la misma lista de archivos aprobada. No hay un paso de tests aparte, e `implement` nunca los corre: el gate `verify` que sigue te ofrece los comandos de chequeo (nunca uno que reescriba archivos) y los corre solo si decís que sí. También podés pedir tests en el propio pedido sin activar el flag.

Por defecto: ausente, que equivale a `false`. `/asdt-init` lo conserva pero nunca lo agrega.

### `code_intelligence`

`/asdt-init` lo escribe solo cuando detecta un índice de code intelligence (hoy, `codegraph`). Los especialistas lo prefieren entonces a los bucles de grep y lectura. Nunca se escribe como `none` — si no se detecta nada, la key se elimina.

### `primary_design_surface`

La superficie para la que UX/UI diseña primero: `mobile`, `tablet`, `desktop` o `none`. `/asdt-init` la pregunta una vez (por defecto `mobile`); un re-init te permite cambiarla.

- **Ausente** — la pregunta nunca se hizo (por ejemplo, una corrida no interactiva). Los especialistas que diseñan o construyen UI (UX/UI y el Developer) la tratan como `mobile` y registran una entrada `ASSUMED:` que lo dice.
- **`none`** — el proyecto no tiene ninguna superficie visual (una CLI, una librería, un servicio de backend). UX/UI entonces se saltea sus pasos de diseño y escribe un resumen corto diciendo que no hay nada que diseñar, para que los especialistas siguientes lean un "sin UI" explícito en lugar de un hueco silencioso. El Developer no agrega capa responsive.

### Hechos detectados

`stack`, `file_structure`, los cuatro campos de decisión (`is_monorepo`, `test_runner`, `naming_style`, `architectural_style`) y los aspectos de `design_fingerprint` (`i18n`, `css_approach`, `orm`, `state_management`, `ci_cd`, `lint`) son lo que detecta `/asdt-init`. Los campos de decisión y los aspectos de la huella se guardan como `{value, source, confidence}`; los especialistas leen solo el `value`. Un aspecto de la huella que se sondeó y no se encontró se guarda como `none`, para que una corrida posterior distinga "sondeado, ausente" de "nunca sondeado".

## Lo que no va en la configuración

**Lo que vos decís del proyecto** — una convención, una trampa que alguien nuevo no vería — no es configuración. Va a la memoria como notas, nunca a un archivo: `/asdt-init` hace unas pocas preguntas opcionales sobre código no obvio y recuerda tus respuestas, y cualquier especialista recuerda lo que le digas ("acordate de que…"). Las notas le llegan a un paso como notas del proyecto separadas y etiquetadas, nunca mezcladas con el bloque de plataforma detectado. Ver [Scratchpad y memoria](/asdt/es/docs/memory-and-engram).

**Los archivos de trabajo** — el resumen de cada especialista — viven en el scratchpad de la sesión, fuera de tu repositorio, nunca dentro de `.asdt/`.

## Variables de entorno

ASDT no lee ninguna variable de entorno directamente. El binario usa el entorno del asistente de IA (Claude Code u OpenCode) para toda la configuración en tiempo de ejecución.

## Múltiples proyectos

Cada proyecto tiene su propio directorio `.asdt/`. Cambiá de proyecto abriendo tu asistente de IA en una carpeta distinta — ASDT lee la configuración desde el `.asdt/config.yaml` más cercano, subiendo desde el directorio de trabajo actual. La memoria también queda separada por proyecto: las notas tienen como alcance el nombre del directorio del proyecto más un hash de su remoto de git (ver [Scratchpad y memoria](/asdt/es/docs/memory-and-engram#el-scratchpad-resúmenes-de-trabajo)).
