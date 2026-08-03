---
title: Product Manager
description: Convierte una petición cruda en requisitos que el resto del equipo puede construir sin adivinar — historias en orden de entrega, alcance explícito y criterios de aceptación.
order: 20
locale: es
---

# Product Manager (`/asdt-pm`)

> Convierte una petición cruda en requisitos que el resto del equipo puede construir sin adivinar — historias en orden de entrega, alcance explícito y criterios de aceptación.

## Qué hace

El PM corre un solo paso y entrega un solo artefacto. Toma lo que pediste y lo convierte en de 3 a 6 historias de usuario, cada una con uno o dos criterios de aceptación en lenguaje llano.

**El orden de las historias es la prioridad.** La primera es la que se entrega primero. No hay ratings MoSCoW, ni un campo de prioridad aparte, ni una matriz que mantener: si querés cambiar la prioridad, cambiás el orden.

**El alcance fuera es obligatorio.** Un hand-off que no dice explícitamente qué queda afuera está incompleto — no "no había nada que excluir". La ambigüedad de alcance es de donde sale casi todo el scope creep, y nombrar lo adyacente es lo que la corta.

Los NFRs solo entran si el feature realmente los implica y si son medibles. "p95 bajo 200ms, medido con k6" es un NFR; "que sea rápido" no es nada.

## Cuándo invocarlo

- La petición está en lenguaje de usuario y hace falta cerrarla antes de diseñar o programar
- El alcance necesita quedar acordado antes de empezar, para que no se expanda en el medio
- Hay varias necesidades que conciliar y alguien tiene que decidir el orden

No lo llames para un refactor, un cambio cosmético o una petición que ya viene con alcance técnico definido. Para eso, andá directo al Developer.

## Cómo invocarlo

Se le habla en lenguaje natural. No hay flags: si querés más o menos profundidad, decilo dentro de la petición y el especialista la ajusta.

```
/asdt-pm "agregar autenticación con email y contraseña"
```

```
/asdt-pm "rediseñar las notificaciones — hay varios equipos involucrados, tomate el tiempo"
```

```
/asdt-pm "exportar CSV en el panel de reportes, algo acotado"
```

## Qué produce

Un único hand-off en `{project}/{change}/pm/handoff`:

| Campo | Qué lleva |
|---|---|
| `what` | El cambio en una frase |
| `decisions` | Las historias en orden de entrega |
| `constraints` | Alcance dentro, alcance fuera, y los NFRs medibles |
| `acceptance_criteria` | Given/When/Then, máximo 5 |
| `risks` | `{riesgo, mitigación}`, una línea cada uno |
| `open_items` | Huecos reales, con prefijo `ASSUMED:` cuando el PM asumió una respuesta |

**El PM es la autoridad de los criterios de aceptación.** Quien viene después los refina —el Developer los lleva a granularidad de implementación, QA busca lo que les falta— pero nadie los reescribe desde cero.

Lo consumen: **Arquitecto**, **Developer**, **QA** y **UX/UI**. Todos lo leen como entrada opcional: si el PM no corrió, cada uno trabaja desde la petición cruda y lo anota.

## Qué consume

La petición cruda, las convenciones detectadas del proyecto, y —si el Researcher corrió antes— su `researcher/handoff`: la dirección recomendada le da el punto de partida, y las direcciones descartadas alimentan el alcance fuera.

Si no hay nada de eso, arranca igual desde la petición y lo deja registrado.

## Límites

- No escribe decisiones de arquitectura ni ADRs
- No escribe código ni diseños técnicos
- No escribe flujos de UX ni specs de componentes
- No produce un hand-off sin alcance fuera explícito
