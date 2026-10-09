---
title: Flujos de usuario
description: Patrones comunes y workflows de punta a punta con ASDT.
order: 10
locale: es
---

# Flujos de usuario

## Obtener una sugerencia de pipeline

Usá `/asdt` cuando no estás seguro qué especialistas necesita una feature:

```
/asdt Agregar login sin contraseña con magic links
```

ASDT analiza la petición, evalúa complejidad y superficie de riesgo, y presenta un plan de routing — por ejemplo:

```
Especialistas recomendados:
  PM — definir alcance y criterios de aceptación
  Arquitecto — diseñar el flujo de tokens y contratos de API
  Developer — implementar el handler de magic link
  Seguridad — revisar el mecanismo de autenticación

Orden sugerido:
  /asdt-pm → /asdt-architect → /asdt-developer → /asdt-security

¿Continuar con este plan? (yes / modify / no)
```

Confirmá el plan. ASDT te da los comandos exactos a ejecutar, en orden, cada uno con la misma versión afinada de tu pedido. **Los ejecutas tú — ASDT no corre los especialistas automáticamente.**

## Ejecutar especialistas directamente

Cuando ya sabes lo que necesitas, salteate `/asdt` e invocá el especialista directamente:

```
/asdt-architect Diseñar la estrategia de rate-limiting para la API
/asdt-qa Revisar el flujo de checkout en busca de casos borde
/asdt-security Auditar la integración OAuth
```

Cada especialista juzga qué pasos propios necesita el pedido — el Developer, por ejemplo, se queda en `explore` ante una pregunta, en `spec` ante un plan, y solo sigue a `approve → implement → verify` cuando le pedís construir — y guarda su único hand-off en la base de conocimiento. Su reporte arranca contando sobre qué trabajó y cierra proponiendo el siguiente especialista.

## Retomar en mitad del pipeline

Si ejecutaste algunos especialistas y quieres continuar después, invocá el siguiente. Lee los artefactos previos de la base de conocimiento automáticamente — incluso de una sesión anterior:

```
/asdt-developer Implementar basándose en el ADR del Arquitecto
```

El Developer lee el hand-off del Arquitecto desde la base de conocimiento. No pasas contexto manualmente.

Lo mismo vale para un plan del Developer. Si le pedís un plan, queda guardado; volvé más tarde con `/asdt-developer "implementá el plan que aprobamos"` y retoma en la compuerta de aprobación en lugar de empezar de cero.

## Memoria y continuidad

ASDT usa un memory provider para persistir artefactos entre especialistas y a través de sesiones. Un memory provider es **requerido** para que el pipeline funcione. La implementación por defecto es [Engram](https://github.com/Gentleman-Programming/engram). Más providers están planificados.

## Asistentes de IA compatibles

Todos los slash commands funcionan tanto en **Claude Code** como en **OpenCode**. La sintaxis es idéntica en ambas herramientas.
