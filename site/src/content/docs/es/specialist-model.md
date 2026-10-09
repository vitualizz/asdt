---
title: Modelo de especialistas
description: Cómo ASDT modela la entrega de software como un equipo de especialistas independientes, cada uno dueño de una disciplina.
order: 6
locale: es
---

# Modelo de especialistas

## Por qué especialistas, no un pipeline

La primera versión de ASDT modelaba la entrega de software como una FSM (finite-state machine — máquina de estados finitos, un flujo rígido que solo se mueve a través de un conjunto fijo de pasos en un orden fijo) de cuatro fases fija: `requirements → plan → implement → review`. Agregar un rol nuevo requería un paquete de Go nuevo, un struct nuevo y una rama de switch nueva — código, no redacción de prompts. La FSM tenía hardcodeado a `requirements` como el único punto de entrada válido, así que un ingeniero de seguridad o un diseñador UX no tenían un lugar válido en el modelo sin reestructurar todo el grafo.

Ese es el modelo equivocado. La entrega real de software la hace un equipo de especialistas, cada uno dueño de una disciplina independiente. Un ingeniero de seguridad no espera a que el desarrollador termine antes de revisar el código de auth. Un diseñador UX no sigue un flujo requirements → plan — sigue su propio proceso creativo.

Por eso ASDT reemplazó la FSM por una unidad distinta: un Especialista es una unidad composable e independiente definida por su identidad, sus propios pasos de workflow, su contrato de artefactos y una garantía de independencia — cualquier especialista puede correr primero, sin ningún predecesor requerido.

## Qué define a un especialista

Un especialista tiene cuatro partes:

**Identidad** — un `id` estable (p. ej. `developer`), un nombre humano y una descripción que el advisor de pipeline usa para enrutar los pedidos.

**Workflow** — una lista corta de pasos propios de esa disciplina, declarada en su `workflow.yaml`. El especialista juzga cuáles necesita el pedido; la profundidad cambia qué tan exhaustivo es el output de cada paso, nunca qué pasos existen. El Developer elige su cadena según lo que pidas: una pregunta corre `explore`; un plan corre `explore → spec` y guarda el plan para retomarlo después; construir corre `explore → spec → approve → implement → verify`, y se detiene a pedir tu aprobación antes de escribir cualquier archivo; y "implementá el plan que aprobamos" retoma en `approve → implement → verify` — o solo en `verify`, cuando el plan ya se construyó pero nunca se chequeó. El especialista UX/UI corre un único `ux-spec` — flujos para la superficie de diseño del proyecto, mapeados a sus componentes — o `review` cuando le pedís auditar lo que ya está en producción. No es el mismo pipeline aplicado a nombres distintos — el workflow de cada especialista refleja cómo funciona realmente esa disciplina.

**Composición de skills** — referencias compartidas (contexto de plataforma, knowledge recall, definición de alcance, OWASP, accesibilidad) declaradas por paso, más skills opcionales que el asistente host ya puede tener instaladas (`host_skills:`, p. ej. `frontend-design` para UX/UI y para los archivos de UI del Developer). Nada se carga de forma ambiental: una referencia compartida se lee únicamente donde está declarada, ya sea como un paso `inline` en `workflow.yaml` o en la lista `reference_skills:` de un paso. Las capacidades se mezclan (mixed in) en lugar de heredarse.

**Contrato de artefactos** — qué hand-offs de sus compañeros lee el especialista (`inputs`) y el único hand-off que escribe, en la clave estable `{project}/{change}/{role}/handoff`, para que otros especialistas lo recuperen por clave. Los inputs son blandos: un input faltante degrada a una nota `ASSUMED:` en `open_items`, nunca a un error.

## Agregar un especialista

Agregar un especialista nuevo requiere dos cosas:

1. Un directorio `skill/asdt-{id}/` — `SKILL.md`, `workflow.yaml` y un archivo por cada paso subagent
2. Su fila en las tablas de routing — la tabla `## Registry` de `skill/SKILL.md` y la tabla `## ASDT Specialists` del template de agents instalado — más la lista de especialistas routed en `skill/embedded_test.go`

Cero paquetes de Go nuevos, cero ramas de switch nuevas. El glob de embed `asdt-*` en `skill/embedded.go` recoge cualquier directorio que coincida con el patrón y lo incluye en el próximo build. Ver [Contribuir](/asdt/docs/contributing) para el contrato de autoría completo.

## La garantía de independencia

Cualquier especialista puede correr primero — no hay un predecesor requerido. Si el Developer no encuentra ningún hand-off del PM en Engram, escribe él mismo los criterios de aceptación y registra `ASSUMED: no PM hand-off — acceptance criteria authored from the exploration` en `open_items`. Su reporte lo dice en la primera línea, que nombra sobre qué trabajó y qué faltaba. El resultado es menos preciso que si el PM hubiera corrido primero, pero es un output válido.

Esta decisión de diseño prioriza la flexibilidad por sobre las garantías de corrección. Siempre puedes correr especialistas fuera de orden. ASDT confía en que tú decides cuándo involucrar a cada disciplina.
