---
title: Diseño UX/UI
description: Da forma a cómo las personas experimentan el producto — flujos de usuario con sus estados y su copy, diseñados para la superficie principal del proyecto y mapeados a sus componentes existentes, con la accesibilidad que cada uno debe — el especialista a invocar antes de que se construya una sola pantalla.
order: 25
locale: es
---

# Diseño UX/UI (`/asdt-ux-ui`)

> Da forma a cómo las personas experimentan el producto — flujos de usuario con sus estados y su copy, diseñados para la superficie principal del proyecto y mapeados a sus componentes existentes, con la accesibilidad que cada uno debe — el especialista a invocar antes de que se construya una sola pantalla.

## Qué hace

El especialista UX/UI convierte un requisito en flujos que un developer puede construir sin adivinar. Corre un solo paso, `ux-spec`, y entrega un solo artefacto, `ux-ui/handoff`, con estas secciones:

- **Brief** — cuatro líneas: el actor, su problema, cómo se ve el éxito para él y la cualidad que debería transmitir la experiencia.
- **Arquitectura de información** — el punto de entrada como ruta completa desde la puerta de la app, la jerarquía de contenido y las acciones principales.
- **Flujos de usuario** — el entregable. Pasos numerados con cada bifurcación, los estados vacío, de carga y de error nombrados como pasos, y el copy exacto escrito en línea donde la redacción sostiene la interacción.
- **Mapeo de componentes** — cada paso del flujo mapeado a un componente que ya existe en el proyecto, por su nombre real. Donde nada encaja, lo dice: ese hueco es una decisión del Developer, nunca un componente inventado en silencio.
- **Accesibilidad** — por componente: foco, teclado, etiquetado y el par de contraste que debe cumplir. Lo que no se puede verificar con los valores disponibles queda como recomendación, nunca como un aprobado.

**Primero el design system.** Si el proyecto ya tiene uno — sus propios tokens, tema o librería de componentes — de ahí salen todos los tokens y componentes, y el especialista nunca inventa una paleta, una escala tipográfica o una unidad de espaciado en su contra. Si no tiene ninguno (una herramienta de estilos sin tema propio del proyecto no cuenta), agrega una **dirección visual**: un par tipográfico y su escala, roles de paleta (`surface`, `text`, `accent`, `danger`, …) cada uno con su par de contraste, una escala de espaciado, una densidad y una postura de motion. Todo es una propuesta, que el Developer construye como los primeros tokens del proyecto.

**Diseñado para tu superficie.** `/asdt-init` pregunta cuál es la superficie de diseño principal del proyecto — `mobile`, `tablet`, `desktop` o `none`. Los flujos se diseñan primero para esa superficie, con una línea sobre qué cambia en las demás; si nunca se preguntó, se asume mobile. Cuando la respuesta es `none` — un CLI, una librería, un servicio backend — no hay pantallas que especificar: el especialista se saltea `ux-spec` y guarda un hand-off corto que lo dice, para que los roles siguientes lean un "sin UI" explícito en lugar de un hueco silencioso.

Si el asistente host tiene instalada una skill `frontend-design`, la usa para elevar la calidad de la dirección visual y de las decisiones de layout. Nunca pasa por encima del design system del proyecto, y si no está, no cambia nada.

## Cuándo invocarlo

- Una pantalla, diálogo o UI a nivel de feature nueva necesita ser diseñada
- Los flujos de usuario necesitan mapearse antes de que empiece la arquitectura o la implementación
- Necesitás saber qué componentes existentes cubren un cambio y dónde están los huecos
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

`ux-ui/handoff` — brief, superficie, arquitectura de información, flujos, dirección visual (solo cuando no hay design system), mapeo de componentes y accesibilidad, como secciones de UN solo artefacto.

Consumido por: **Arquitecto** (los flujos y los huecos de componentes que el diseño tiene que cubrir), **Developer** (implementa los flujos, cubre los huecos y construye la dirección visual como primeros tokens), **QA** (cada bifurcación de flujo y cada estado vacío, de carga y de error se vuelve un caso borde candidato, y el copy en línea se vuelve una aserción de texto exacto).

## Patrones comunes

```
/asdt-ux-ui Diseñar el flujo de onboarding para nuevos usuarios
# → UI multi-paso nueva — IA y flujos antes de cualquier trabajo de componentes
```

```
/asdt-ux-ui Mapear la pantalla de preferencias de notificaciones
# → UI existente a extender — el mapeo nombra qué se reutiliza y qué falta
```

```
/asdt-ux-ui Diseñar las primeras pantallas del panel de administración
# → Todavía no hay design system — el hand-off trae una dirección visual propuesta
```

## Límites — qué NO hace

- No escribe código de implementación ni ningún archivo — solo el hand-off
- No produce decisiones de arquitectura ni planes de prueba
- Nunca inventa tokens ni componentes en contra de un design system existente — una necesidad sin cubrir es un hueco con nombre
- Nunca escribe un segundo juego de flujos por superficie — un diseño para la superficie principal y una línea de adaptación para el resto
- Nunca da por aprobada una accesibilidad que no puede verificar con los tokens o la paleta propuesta
