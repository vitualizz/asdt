---
title: Configuración
description: Referencia de .asdt/config.yaml, platform.yaml y todas las opciones de configuración.
order: 2
locale: es
---

# Configuración

## Configuración del proyecto — `.asdt/config.yaml`

Ejecutar `/asdt-init` (dentro de tu asistente) crea `.asdt/config.yaml` con valores por defecto razonables. Editalo para que se ajuste a tu proyecto.

```yaml
memory:
  provider: engram

specialists:
  strict_tdd: false
```

### `memory.provider`

El backend de memoria que usan todos los especialistas para persistir y recuperar artefactos. Hoy solo se soporta `engram`.

### `specialists.strict_tdd`

Cuando es `true`, el workflow del especialista Developer incluye un paso `test` después de `implement`. Los tests se producen como artefactos estructurados — no se ejecutan — así que este flag controla la generación, no el CI.

Valor por defecto: `false`.

## Contexto de plataforma — `.asdt/knowledge/platform.yaml`

Generado por `/asdt-init`. Describe el contexto técnico del proyecto para que cada especialista pueda ajustar su output sin tener que redetectar el stack.

```yaml
language: typescript
framework: astro
package_manager: npm
test_runner: vitest
lint: eslint
```

Los campos son de formato libre — los especialistas leen lo que está presente y se saltean lo que falta. Podés agregar campos personalizados como `deployment: vercel` o `database: postgres` y referenciarlos en los prompts de los especialistas a través del *shared skill* `platform-context` (un fragmento de instrucciones que cada paso de cada especialista carga automáticamente).

## Variables de entorno

ASDT no lee ninguna variable de entorno directamente. El binario usa el entorno del asistente de IA (Claude Code u OpenCode) para toda la configuración en tiempo de ejecución.

## Múltiples proyectos

Cada proyecto tiene su propio directorio `.asdt/`. Cambiá de proyecto abriendo tu asistente de IA en una carpeta distinta — ASDT lee la configuración desde el `.asdt/config.yaml` más cercano, subiendo desde el directorio de trabajo actual.
