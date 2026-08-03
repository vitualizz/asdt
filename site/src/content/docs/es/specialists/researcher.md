---
title: Researcher
description: Toma un problema difuso y devuelve UNA dirección recomendada con factibilidad detrás — el especialista a invocar antes de que existan los requisitos, cuando todavía no sabés qué construir.
order: 26
locale: es
---

# Researcher (`/asdt-researcher`)

> Toma un problema difuso y devuelve UNA dirección recomendada con factibilidad detrás — el especialista a invocar antes de que existan los requisitos, cuando todavía no sabés qué construir.

## Qué hace

Diverge antes de que el PM converja, en un solo paso:

1. **Encuadra el problema.** Una petición difusa casi siempre esconde dos o tres problemas distintos. Nombra cuál está explorando y cuáles deja de lado a propósito.
2. **Diverge.** Genera de 3 a 5 direcciones genuinamente distintas. Distintas significa que fallan por razones distintas — tres variaciones de una misma idea son una dirección, no tres.
3. **Juzga factibilidad.** Verde, amarillo o rojo por dirección, cada una con **una línea de evidencia**: un fichero que ya hace algo parecido, una dependencia que falta, una restricción que la descarta. Un veredicto sin evidencia detrás no es un veredicto — se marca `ASSUMED:` y queda a la vista.
4. **Converge.** Recomienda exactamente una dirección y dice por qué le ganó a las otras. Cada descartada se lleva su razón: una dirección rechazada sin motivo escrito es la que vuelve el trimestre que viene.

## Cuándo invocarlo

- No sabés todavía qué construir, solo qué duele
- Hay varios caminos plausibles y querés ver los tradeoffs antes de comprometerte
- Alguien propuso una solución y querés saber si es la única

## Cómo invocarlo

Lenguaje natural, sin flags. Si querés que se tome más tiempo o menos, decilo en la petición.

```
/asdt-researcher "los usuarios abandonan el onboarding y no sabemos en qué punto"
```

```
/asdt-researcher "queremos búsqueda semántica — explorá a fondo antes de que decidamos"
```

## Qué produce

Un único hand-off en `{project}/{change}/researcher/handoff`:

| Campo | Qué lleva |
|---|---|
| `what` | La dirección recomendada en una frase |
| `decisions` | La recomendación primero, después cada dirección descartada como `rejected: {dirección} — {porqué}` |
| `constraints` | Lo que cualquier implementación de esa dirección tiene que respetar |
| `files_hint` | El código que realmente leyó para juzgar factibilidad |
| `risks` | `{riesgo, mitigación}` de la dirección recomendada |
| `open_items` | Veredictos que no pudo anclar en evidencia, con prefijo `ASSUMED:` |

En `decisions` es donde sobrevive la exploración: no solo qué eligió, sino qué miró y descartó.

Lo consume el **PM**, como entrada opcional: la dirección recomendada le da el punto de partida y las descartadas alimentan su alcance fuera.

## Su lugar en el pipeline

Es el único especialista **pre-requisitos**. Corre antes que el PM y nunca lo reemplaza: recomienda una dirección, el PM decide qué se construye. También funciona solo, cuando lo único que querés es exploración estructurada.

## Límites

- No escribe requisitos, arquitectura, código ni tests
- Nunca escribe en el sistema de ficheros
- No rankea durante la divergencia — ranquear temprano tira la opción que todavía no terminaste de pensar
