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

¿Corro esta cadena? Hago una pausa después de cada especialista.
```

Cada comando lleva la misma versión afinada de tu pedido. Si decís que sí, `/asdt` corre la cadena de a un especialista por vez y **hace una pausa después de cada uno**: podés debatir el resultado, y solo sigue cuando le das el visto bueno (o pasa al especialista que el debate volvió más urgente). O ejecutás vos los comandos, en orden. **Nada corre sin tu sí.**

## Ejecutar especialistas directamente

Cuando ya sabes lo que necesitas, salteate `/asdt` e invocá el especialista directamente:

```
/asdt-architect Diseñar la estrategia de rate-limiting para la API
/asdt-qa Revisar el flujo de checkout en busca de casos borde
/asdt-security Auditar la integración OAuth
```

Cada especialista juzga qué pasos propios necesita el pedido — el Developer, por ejemplo, se queda en `explore` ante una pregunta, en `spec` ante un plan, y solo sigue a `approve → implement → verify` cuando le pedís construir — y devuelve su único resumen, que el orquestador escribe en el scratchpad de la sesión para el siguiente especialista. Su reporte arranca contando sobre qué trabajó y cierra proponiendo el siguiente especialista.

## Retomar en mitad del pipeline

Si ejecutaste algunos especialistas y querés continuar, invocá el siguiente. En la misma sesión, lee los resúmenes anteriores desde el scratchpad:

```
/asdt-developer Implementar basándose en el diseño del Arquitecto
```

El Developer recibe el resumen del Arquitecto como input. No pasás contexto manualmente. En una sesión posterior puede que el scratchpad ya no esté, y es normal: el Developer trabaja a partir del código, de git y de las decisiones clave que recuerda la memoria, y te pide lo que le falte.

Lo mismo vale para un plan del Developer. Si le pedís un plan, la corrida termina ahí; volvé más tarde con `/asdt-developer "implementá el plan que aprobamos"` y retoma en la compuerta de aprobación en lugar de empezar de cero. Si la sesión terminó en el medio, reconstruye el plan a partir de git, la memoria y tus palabras, y vos aprobás el plan reconstruido.

## Memoria y continuidad

Dentro de una sesión, los especialistas se pasan el trabajo mediante resúmenes en el scratchpad de la sesión. Entre sesiones, un memory provider guarda solo el conocimiento clave del proyecto: decisiones no obvias, conclusiones de debates, convenciones, lecciones y hallazgos de auditorías. El provider actual es [Engram](https://github.com/Gentleman-Programming/engram). Sin un provider conectado, ASDT igual funciona dentro de la sesión; simplemente no pasa nada a la siguiente. Ver [Scratchpad y memoria](/asdt/es/docs/memory-and-engram).

## Asistentes de IA compatibles

Todos los slash commands funcionan tanto en **Claude Code** como en **OpenCode**. La sintaxis es idéntica en ambas herramientas.
