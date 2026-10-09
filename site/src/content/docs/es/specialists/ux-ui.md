---
title: Diseño UX/UI
description: Diseña cómo se ve y cómo funciona un cambio antes de que se construya una sola pantalla — flujos de usuario con sus estados y su copy, y después cada pantalla y cada componente faltante diseñados en detalle, en el design system del proyecto o en una base visual propuesta, con la accesibilidad que cada uno debe.
order: 25
locale: es
---

# Diseño UX/UI (`/asdt-ux-ui`)

> Diseña cómo se ve y cómo funciona un cambio antes de que se construya una sola pantalla — flujos de usuario con sus estados y su copy, y después cada pantalla y cada componente faltante diseñados en detalle, en el design system del proyecto o en una base visual propuesta, con la accesibilidad que cada uno debe.

## Qué hace

El especialista UX/UI hace el trabajo de un diseñador de producto: decide cómo fluye el cambio y después diseña la UI que necesita, con la precisión suficiente para que el Developer la construya sin adivinar. Corre dos pasos — `ux-spec` escribe los flujos, `ui-design` diseña las pantallas y los componentes — y entrega un solo artefacto, `ux-ui/handoff`. El diseño es texto dentro de ese hand-off; no se escribe ningún archivo.

De `ux-spec`:

- **Brief** — cuatro líneas: el actor, su problema, cómo se ve el éxito para él y la cualidad que debería transmitir la experiencia.
- **Arquitectura de información** — el punto de entrada como ruta completa desde la puerta de la app, la jerarquía de contenido y las acciones principales.
- **Flujos de usuario** — pasos numerados con cada bifurcación, la pantalla en la que ocurre cada uno, los estados vacío, de carga y de error nombrados como pasos, y el copy exacto escrito en línea donde la redacción sostiene la interacción.
- **Mapeo de componentes** — cada paso del flujo mapeado a un componente que ya existe en el proyecto, por su nombre real. Donde nada encaja, el hueco recibe un nombre — y `ui-design` lo diseña.
- **Accesibilidad** — por componente: foco, teclado, etiquetado y el par de contraste que debe cumplir. Lo que no se puede verificar con los valores disponibles queda como recomendación, nunca como un aprobado.

De `ui-design`:

- **Pantallas** — una por cada pantalla que tocan los flujos: layout (ancho máximo, grilla, alineación, y cada región con su tamaño y contenido), jerarquía en orden de lectura con el nivel tipográfico de cada elemento y una sola acción principal, espaciado en pasos de la escala, los componentes que compone, cada estado y exactamente qué cambia en él, cómo maneja los extremos de datos (un ítem, muchísimos, el texto más largo, un campo faltante), qué hace cada región en cada breakpoint, y motion solo donde difiere de la base. El copy queda en los pasos del flujo. Por ejemplo:

  ```
  Pantalla: Solicitar reset
    Layout: columna centrada, max 420px
    Jerarquía: H1 → texto → input → CTA
    Tipografía: H1 Fraunces 28/34 600 (Google Fonts)
    Estados: vacío / enviando / error
  ```

- **Diseño de componentes** — cada componente que los flujos necesitan y el proyecto no tiene, más cualquier variante nueva de un componente existente, diseñado parte por parte: anatomía, props, variantes, tamaños, estados (default, hover, focus-visible, active, disabled, loading, error), el token que consume cada parte y su comportamiento. Un componente que se reutiliza tal cual no lleva diseño — la pantalla que lo usa lo marca como tal cual.

El oficio de diseño viene incluido: una referencia compartida de diseño visual — una intención estética que se hace visible en dos o tres movimientos distintivos, una escala tipográfica sobre fuentes nombradas con su origen, una escala de espaciado y breakpoints, jerarquía, color por rol con familias interactivas y pares de contraste, densidad, elevación e íconos, estados — variantes de vacío y extremos de datos incluidos — como artefactos diseñados, motion que respeta el movimiento reducido, anillos de foco visibles, y anatomía y props de componentes — así que la calidad no depende de ningún plugin. Antes de entregar, `ui-design` contrasta su diseño con esa referencia y corrige lo que falla.

**Primero el design system.** Si el proyecto ya tiene uno — sus propios tokens, tema o librería de componentes — esa es la base: sus tokens y componentes se usan como son, nunca se reestilizan: una variante nueva agrega una opción en el lenguaje del sistema y deja intactas las existentes, un valor que el sistema no tiene se propone como un token con nombre (`--space-7: 28px`) para que el Developer lo agregue, y un componente que le falta se diseña en su lenguaje. Si no tiene ninguno (una herramienta de estilos sin tema propio del proyecto no cuenta), `ui-design` propone una **dirección visual**: la intención y sus movimientos distintivos, las fuentes con su origen y una escala tipográfica, roles de paleta (`surface`, `text`, `accent` con sus valores hover, pressed y on-accent, `focus`, `danger`, …) cada uno con su par de contraste, una escala de espaciado, radios, elevación, íconos, una densidad y una postura de motion. Todo es una propuesta, que el Developer construye como los primeros tokens del proyecto.

**Diseñado para tu superficie.** `/asdt-init` pregunta cuál es la superficie de diseño principal del proyecto — `mobile`, `tablet`, `desktop` o `none`. Los flujos y las pantallas se diseñan primero para esa superficie, con una línea sobre qué cambia en las demás; si nunca se preguntó, se asume mobile. Cuando la respuesta es `none` — un CLI, una librería, un servicio backend — no hay pantallas que diseñar: el especialista se saltea los dos pasos y guarda un hand-off corto que lo dice, para que los roles siguientes lean un "sin UI" explícito en lugar de un hueco silencioso. Un cambio sin ningún paso de cara al usuario — un job en segundo plano, un cambio solo de API — recibe el mismo no explícito en cualquier superficie: los dos pasos corren y no devuelven flujos ni pantallas.

Si el asistente host tiene instalada una skill `frontend-design`, `ui-design` la usa para elevar todavía más el oficio. Nunca pasa por encima del design system del proyecto, de los flujos ni de la accesibilidad, y si no está, no cambia nada.

## Cuándo invocarlo

- Una pantalla, diálogo o UI a nivel de feature nueva necesita ser diseñada
- Los flujos de usuario necesitan mapearse antes de que empiece la arquitectura o la implementación
- Necesitás saber qué componentes existentes cubren un cambio, y que los que faltan queden diseñados
- Querés las pantallas diseñadas — layout, jerarquía, tipografía, estados — antes de que alguien las construya
- El proyecto todavía no tiene design system y las primeras pantallas necesitan una dirección coherente
- Los requisitos de accesibilidad necesitan especificarse explícitamente
- Querés que el Developer reciba una spec en lugar de inferir la UX de los requisitos

## Por su cuenta

Apuntalo a una pantalla o un flujo que ya existe:

```
/asdt-ux-ui "revisá la accesibilidad del checkout"
/asdt-ux-ui "¿qué componentes del design system no estamos usando en el onboarding?"
```

Eso corre `review`: fricción, estados faltantes, desvíos del design system y accesibilidad en lo que ya está en producción.

## Posición en el pipeline

Funciona mejor **después del PM** (lee `pm/handoff` para el requisito y sus criterios de aceptación) y **antes del Arquitecto y del Developer**. El Arquitecto lee los flujos para darle forma a la superficie de API que tiene que servir; el Developer los construye paso a paso. Correrlo después de que una pantalla ya está construida significa que la spec llega demasiado tarde para guiarla.

## Qué produce

`ux-ui/handoff` — brief, superficie, arquitectura de información, flujos, mapeo de componentes, pantallas, diseño de componentes, dirección visual (solo cuando no hay design system) y accesibilidad, como secciones de UN solo artefacto.

Consumido por: **Arquitecto** (los flujos y los huecos de componentes que el diseño tiene que cubrir), **Developer** (construye cada pantalla y componente según su diseño, implementa los flujos y construye la dirección visual como primeros tokens, agrega los tokens que propuso), **QA** (cada bifurcación de flujo, cada estado y extremo de datos de pantalla, y cada estado disabled, loading y error de un componente diseñado se vuelve un caso borde candidato, y el copy en línea se vuelve una aserción de texto exacto).

## Patrones comunes

```
/asdt-ux-ui Diseñar el flujo de onboarding para nuevos usuarios
# → UI multi-paso nueva — IA y flujos antes de cualquier trabajo de componentes
```

```
/asdt-ux-ui Mapear la pantalla de preferencias de notificaciones
# → UI existente a extender — el mapeo nombra qué se reutiliza tal cual, y lo que falta se diseña
```

```
/asdt-ux-ui Diseñar las primeras pantallas del panel de administración
# → Todavía no hay design system — las pantallas se diseñan sobre una dirección visual propuesta
```

## Límites — qué NO hace

- No escribe código de implementación ni ningún archivo — el diseño vive en el hand-off
- No produce decisiones de arquitectura ni planes de prueba
- Nunca reestiliza un design system existente ni inventa un token en su contra — un componente faltante se diseña en su lenguaje, una variante nueva deja intactas las existentes, un token faltante se propone con nombre y valor como ítem abierto
- Nunca escribe un segundo diseño por superficie — uno para la superficie principal y una línea de adaptación para el resto
- Nunca da por aprobada una accesibilidad que no puede verificar con los tokens o la paleta propuesta
