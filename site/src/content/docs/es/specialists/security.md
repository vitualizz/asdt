---
title: Seguridad
description: Busca los huecos que encontraría un atacante primero — modelos de amenazas, revisiones OWASP y checklists de hardening — el especialista a invocar siempre que auth, manejo de datos o integraciones externas estén sobre la mesa, en cualquier punto del pipeline.
order: 24
locale: es
---

# Seguridad (`/asdt-security`)

> Busca los huecos que encontraría un atacante primero — modelos de amenazas, revisiones OWASP y checklists de hardening — el especialista a invocar siempre que auth, manejo de datos o integraciones externas estén sobre la mesa, en cualquier punto del pipeline.

## Qué hace

El especialista de Seguridad realiza modelado de amenazas y análisis de seguridad usando STRIDE y el OWASP Top 10. Mapea la superficie de ataque, identifica amenazas sistemáticamente y produce un checklist de hardening priorizado donde cada hallazgo tiene una mitigación concreta y accionable — no "monitorearlo" ni "agregar logging."

El invariante crítico: **Seguridad no tiene predecesor requerido.** Puede correr en cualquier etapa — en un proyecto nuevo sin artefactos previos, a mitad del desarrollo o después del lanzamiento. Si existen artefactos previos (decisiones de arquitectura, implementación), los lee. Si no, trabaja desde el contexto de plataforma y la petición sola, notando las brechas en `open_items` y continuando.

La profundidad está controlada por la superficie de riesgo, no por complejidad — el especialista la juzga según lo que toca el cambio: autenticación, secretos, manejo de datos, integraciones externas. Este es el único especialista donde la pregunta no es "¿qué tan compleja es la feature?" sino "¿qué tan grande es la superficie de ataque?" La superficie de riesgo define qué tan profundo va el análisis; los dos pasos, `assess` y `harden`, corren siempre.

## Cuándo invocarlo

- Autenticación, gestión de sesiones o autorización están involucrados
- La feature maneja o almacena información de identificación personal
- Hay integraciones externas, webhooks o URLs controladas por el usuario
- Se están exponiendo nuevos endpoints de API públicamente
- En cualquier momento antes de salir a producción cuando la seguridad no ha sido revisada

## Por su cuenta

Es el especialista que más se usa solo. No requiere que haya un cambio en curso:

```
/asdt-security "audita el módulo de pagos"
/asdt-security "revisa cómo manejamos las sesiones"
/asdt-security "¿qué expone nuestro endpoint de webhooks?"
```

Mapea la superficie, la evalúa y te deja hallazgos priorizados con mitigación concreta.

## Posición en el pipeline

**Sin predecesor requerido** — invocalo en cualquier punto. Para máximo impacto, correlo después del Arquitecto (lee `architect/handoff` para la superficie de API y los límites de confianza, y `developer/handoff` para el código que cambió, cuando existen). Si recibe un plan del Developer todavía sin construir, mapea la superficie planificada a partir de los archivos que declara el plan y registra que todavía no hay nada construido. Para un modelo de amenazas temprano en el diseño, correlo antes de que la arquitectura esté finalizada para exponer riesgos a nivel de diseño antes de que queden incorporados.

Su hand-off lo leen el Arquitecto, el Developer y QA.

## Qué produce

`security/handoff` — hallazgos y checklist de hardening como secciones de UN solo artefacto (`{project}/study/{topic}/security` cuando la corrida audita algo que ya existe):

- **Hallazgos** — de mayor a menor severidad, cada uno con una severidad de una palabra (`high`, `medium`, `low`), qué gana un atacante, la evidencia que lo respalda y la mitigación concreta que lo cierra
- **Checklist de hardening** — ítems ordenados, cada uno verificable como hecho, más dónde caen las mitigaciones

Consumido por: **Arquitecto** (un hallazgo que redefine un límite es una restricción de diseño), **Developer** (cada mitigación que cae en los archivos del cambio se vuelve una restricción del spec), **QA** (cada hallazgo recibe un caso de prueba que demuestra que su mitigación se sostiene; un `high` cuya mitigación no deja rastro en los archivos cambiados o planificados del Developer es un `no-go`).

## Patrones comunes

```
/asdt-security Auditar la integración OAuth
# → Flujo de auth externo con manejo de tokens — superficie de riesgo alta
```

```
/asdt-security Modelar amenazas para el nuevo handler de webhooks de pago
# → Input controlado por el usuario que impacta lógica financiera
```

```
/asdt-security Revisión de seguridad rápida antes del lanzamiento de v2
# → No se necesitan artefactos previos — corre solo desde el contexto de plataforma
```

## Límites — qué NO hace

- No escribe código de implementación
- No produce decisiones de arquitectura ni specs de UX
- No produce planes de prueba (aunque sus hallazgos informan qué debería cubrir QA)
- Cada hallazgo debe tener una mitigación concreta — "agregar monitoring" no es una mitigación
- La severidad es una palabra — `high`, `medium` o `low`. Sin CVSS, sin puntajes numéricos
- Nunca corre scanners, auditorías de dependencias ni ningún otro comando — razona sobre el cambio y lee el repositorio
