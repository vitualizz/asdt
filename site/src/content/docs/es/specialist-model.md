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

Un Architecture Decision Record (ADR-006) — la nota de diseño interna que garantiza que cualquier especialista pueda correr de forma independiente, en cualquier orden — formalizó el cambio: un Especialista es una unidad composable e independiente definida por su identidad, sus propios pasos de workflow, su contrato de artefactos y una garantía de independencia — cualquier especialista puede correr primero.

## Qué define a un especialista

Un especialista tiene cuatro partes:

**Identidad** — un `id` estable (p. ej. `developer`), un nombre humano y una descripción que el advisor de pipeline usa para enrutar los pedidos.

**Workflow** — una lista ordenada de pasos específicos de esa disciplina. El Developer corre `explore → spec → design → implement`. El especialista UX/UI corre `feature-brief → information-architecture → user-flows → component-mapping → ux-handoff`. No es el mismo pipeline aplicado a nombres distintos — el workflow de cada especialista refleja cómo funciona realmente esa disciplina.

**Composición de skills** — skills compartidas que se cargan en cada paso (contexto de plataforma, envelope de artefactos, definición de alcance), más skills propias del especialista que se cargan por paso (threat modeling para Security, generación de código para Developer). Las capacidades se mezclan (mixed in) en lugar de heredarse.

**Contrato de artefactos** — qué lee el especialista (`inputs`) y qué escribe (`outputs`). Los inputs son blandos: un input faltante degrada a una nota en `open_items[]`, nunca a un error. Los outputs tienen valores `topic_key` estables para que otros especialistas puedan recuperarlos por clave.

## Agregar un especialista

Agregar un especialista nuevo requiere exactamente dos cosas:

1. Un literal de valor `SpecialistDescriptor` en el registro
2. Un árbol `skill/{id}/SKILL.md` con los archivos de paso

Cero paquetes de Go nuevos, cero ramas de switch nuevas. El glob de embed `asdt-*` en `skill/embedded.go` recoge cualquier directorio que coincida con el patrón y lo incluye en el próximo build. Ver [Contribuir](/asdt/docs/contributing) para el contrato de autoría completo.

## La garantía de independencia

Cualquier especialista puede correr primero — no hay un predecesor requerido. Si el Developer no encuentra ningún artefacto del Architect en Engram, sigue adelante con `open_items: ["architect/adr not found"]` y hace suposiciones razonables. El artefacto de implementación resultante es menos preciso que si el Architect hubiera corrido primero, pero es un output válido.

Esta decisión de diseño prioriza la flexibilidad por sobre las garantías de corrección. Siempre podés correr especialistas fuera de orden. ASDT confía en que vos decidís cuándo involucrar a cada disciplina.
