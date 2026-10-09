---
title: Arquitecto
description: Toma decisiones de arquitectura y produce ADRs, diseño de sistema y artefactos de diseño de API — el especialista a invocar cuando una decisión va a dar forma a los límites de servicios, modelos de datos o escalabilidad a largo plazo.
order: 21
locale: es
---

# Arquitecto (`/asdt-architect`)

> Toma decisiones de arquitectura y produce ADRs, diseño de sistema y artefactos de diseño de API — el especialista a invocar cuando una decisión va a dar forma a los límites de servicios, modelos de datos o escalabilidad a largo plazo.

## Qué hace

El especialista Arquitecto toma las decisiones técnicas sobre las que se construye todo lo demás. Evalúa enfoques en competencia, documenta el camino elegido como un Architecture Decision Record (ADR) y produce un diseño de sistema concreto con modelos de datos, superficies de API y límites de servicios — todo antes de que se escriba una sola línea de código de implementación.

Cada decisión nombra dos o tres enfoques viables y registra, en el mismo bloque, por qué ganó el elegido y por qué perdió cada alternativa — ese bloque es el registro de decisión; no hay un artefacto de ADR aparte. Por defecto elige el enfoque más simple que satisface los criterios de aceptación, y solo invierte rigor extra donde una decisión es difícil de revertir o visible para otros. Esto fuerza un análisis honesto de trade-offs en lugar de justificaciones post-hoc.

El especialista Arquitecto nunca escribe código de implementación, specs de UX ni planes de prueba. Su único trabajo es tomar la decisión estructural que el Developer puede implementar sin ambigüedad.

## Cuándo invocarlo

- Una decisión dará forma a los límites de servicios, modelos de datos o escalabilidad más allá del feature actual
- El enfoque técnico no es obvio y hay trade-offs significativos entre al menos dos opciones viables
- Una preocupación transversal (estrategia de caché, modelo de auth, event bus) necesita una decisión documentada
- Querés un ADR formal para explicar a futuros ingenieros por qué el código es como es

## Por su cuenta

No hace falta un cambio en marcha. Apúntalo a lo que ya existe y lo juzga en vez de rediseñarlo — te devuelve hallazgos priorizados con evidencia, y también lo que está bien:

```
/asdt-architect "¿escala esta estructura si triplicamos el tráfico?"
/asdt-architect "audita los boundaries del módulo de pagos"
/asdt-architect "¿qué decisiones de este diseño ya son caras de revertir?"
```

Lo que encuentre queda guardado, así que la próxima corrida sobre esa área arranca sabiéndolo.

## Posición en el pipeline

Típicamente corre **después del PM** y **antes del Developer** (el Developer lee `architect/handoff`). Lee lo que exista antes que él: `pm/handoff` (los requisitos que el diseño satisface), `ux-ui/handoff` (los flujos y huecos de componentes que la superficie de API tiene que servir), `security/handoff` (un hallazgo que redefine un límite se vuelve una restricción de diseño) y `researcher/handoff` — que enmarca el problema solo cuando se salteó el PM; si el PM corrió, manda el PM. Todos son opcionales. En cambios simples no se invoca — de eso se encarga el Developer directamente. Cuando corre, corre un solo paso, `design`, y qué tan profundo va lo decide él mismo.

## Qué produce

`architect/handoff` — la decisión y el diseño de sistema que se desprende de ella, en UN solo hand-off:

- **La decisión** — primero el enfoque elegido, después cada alternativa descartada con el porqué
- **El diseño** — el modelo de datos y la superficie de API (o por qué el cambio no tiene ninguno), las restricciones que la implementación tiene que respetar, dónde cae en el codebase y los riesgos con sus mitigaciones

Juzgar lo que ya existe corre `review` y guarda `{project}/study/{topic}/architect`.

Consumido por: **Developer** (la decisión ya está tomada — su spec la reformula al nivel de implementación), **QA** (el diseño y sus riesgos declarados), **Security** (la superficie de API y los límites de confianza).

Los presupuestos de NFR, si el PM fijó alguno, llegan dentro de `pm/handoff.constraints` y el diseño tiene que caber en ellos. Si el PM no corrió, el diseño sigue igual y anota el faltante en vez de inventar un presupuesto.

## Patrones comunes

```
/asdt-architect Diseñar la estrategia de rate-limiting para la API pública
# → Preocupación transversal que afectará cada endpoint
```

```
/asdt-architect Elegir el enfoque de event sourcing para el pipeline de órdenes
# → Decisión estructural no reversible con trade-offs significativos
```

```
/asdt-architect ADR para migrar de REST a GraphQL en el cliente móvil
# → Cambio de contrato externo que necesita rationale documentado
```

## Límites — qué NO hace

- No escribe código de implementación
- No escribe specs de UX ni wireframes
- No produce planes de prueba ni criterios de aceptación
- Nunca omite alternativas — cada registro de decisión las requiere
- No diseña en aislamiento — siempre tiene en cuenta las restricciones de la plataforma existente
- El diseño de sistema siempre está incompleto sin un modelo de datos Y una superficie de API
