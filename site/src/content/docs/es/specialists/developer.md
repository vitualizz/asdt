---
title: Developer
description: Convierte specs y diseños en código funcional — planes de implementación, código de producción y suites de tests — el especialista a invocar una vez que la forma de la solución está definida y es momento de construirla.
order: 22
locale: es
---

# Developer (`/asdt-developer`)

> Convierte specs y diseños en código funcional — planes de implementación, código de producción y suites de tests — el especialista a invocar una vez que la forma de la solución está definida y es momento de construirla.

## Qué hace

El Developer trabaja spec-first. Después de cargar las convenciones y la superficie de diseño del proyecto (el preludio inline `platform-analysis`), antes que nada lee el código afectado, escribe un spec — qué entra y qué queda fuera del alcance, los criterios de aceptación, el enfoque técnico y los archivos exactos que tiene permitido tocar — y recién entonces, con tu aprobación, escribe código.

Decide su propia cadena según lo que le pidas:

| Le pedís | Corre |
|---|---|
| una pregunta o un chequeo rápido | `explore` |
| un plan — "¿cómo lo harías?" | `explore → spec` — el plan pasa a ser su resumen y la corrida termina; no se escribe código |
| el cambio en sí | `explore → spec → approve → implement → verify` |
| "implementá el plan que aprobamos" | `approve → implement → verify` — vuelve a mostrar el plan del scratchpad de la sesión en lugar de rehacerlo (o lo reconstruye, si el scratchpad ya no está); un plan ya construido pero nunca verificado se retoma directo en `verify` |

Si el pedido es ambiguo entre un plan y una construcción, produce el plan, y construirlo después es retomarlo, con una nueva aprobación en la compuerta. Una reanudación arranca desde el scratchpad, nunca desde una aprobación anterior: una vez confirmado el cambio (si varios podrían coincidir, una pregunta decide), lee su propio resumen. Un plan vuelve a la compuerta de aprobación; una construcción nunca verificada va a `verify`; en cualquier otro caso, dice qué contiene ya el resumen y se detiene. Cuando el scratchpad ya no está, algo normal en una sesión nueva, reconstruye el plan a partir de git, de las notas recordadas y de tus palabras. Como los resúmenes anteriores también se perdieron, te pregunta qué especialistas dieron forma al plan y te ofrece volver a correrlos primero; también te pide que repitas el plan si hace falta. Después la compuerta te muestra ese plan reconstruido para que lo apruebes. Una aprobación dada en otra sesión nunca se arrastra.

**La compuerta de aprobación.** Antes de escribir cualquier archivo de tu repo, te muestra el plan en prosa: qué entra en el alcance y qué queda explícitamente fuera, los criterios de aceptación, la dirección visual cuando la hay, los nombres de las pantallas y componentes que diseñó UX/UI, los archivos exactos que va a crear y modificar, y todo lo que asumió o marcó — como un archivo de UI en un proyecto sin superficie visual. Si el cambio ya se entregó una vez, también te avisa que esta construcción se suma sobre esa entrega anterior, que sigue siendo parte del cambio. Respondés aprobar, ajustar o frenar. Ajustar vuelve a correr el spec una vez con tus palabras; frenar termina la corrida en el plan, listo para retomar. Si no hay un humano que pueda responder, no se escribe nada.

**Alcance de escritura.** `implement` solo escribe dentro de los archivos que declaró el spec. Si una edición necesaria cae fuera, se detiene y reporta el path en lugar de escribirlo por su cuenta. Nada bajo `.asdt/` es nunca un destino de edición del Developer: un cambio que haga falta ahí se reporta como un ítem abierto que te deriva a `/asdt-init`, dueño de ese directorio. Un spec que no declara archivos es una corrida plan-only: el código vuelve como snippets en su resumen y el repo no se toca.

**Tests.** No hay un paso de test aparte. Cuando `strict_tdd: true` está en `.asdt/config.yaml`, o cuando pedís tests, `implement` los escribe en la misma pasada y con el mismo alcance de archivos que el código.

**La compuerta de verificación.** `implement` escribe código y tests, pero nunca ejecuta nada. Después, el Developer te muestra los comandos de chequeo que sugiere — build, lint, tests — y cómo se ve una corrida sana, y te pregunta si los corre. Un comando que escribiría tus archivos fuente — un flag `--fix` o `--write`, una actualización de snapshots, un formateador que reescribe, codegen — nunca se ofrece; el output de build y los reportes de cobertura no cuentan, y se prefiere una variante que no escribe (`--ci`, `tsc --noEmit`) cuando existe. Con tu sí, corre exactamente esos comandos y nada más; sin un sí, el resultado queda registrado como no corrido, nunca como aprobado. Si algo falla, tiene como máximo dos rondas de corrección, dentro de los mismos archivos que declaró el spec, volviendo a correr los mismos comandos; si después de eso sigue fallando, se detiene y registra qué falla. Una reanudación que entra directo a la verificación no tiene el spec a mano, así que una falla ahí se registra sin rondas de corrección.

**Archivos de UI.** Cuando un cambio escribe vistas, componentes o estilos, los diseña primero para la superficie de diseño del proyecto. Cuando corrió UX/UI, construye cada pantalla y componente según el diseño del resumen de UX/UI — layout, jerarquía, tipografía, cada estado y cada breakpoint. Con design system, agrega solo los tokens que UX/UI propuso por nombre, en el archivo de tokens del sistema, y una variante nueva deja intactas las existentes. Un design system existente — o, si no lo hay, la dirección visual que propuso UX/UI — está por encima de todo salvo del alcance aprobado. Si el asistente host tiene instalada una skill `frontend-design`, la usa para dar forma a la ejecución visual, y esa skill cede ante todo lo anterior.

## Cuándo invocarlo

- La forma de la solución ya está definida (requisitos, arquitectura o UX están definidos)
- Querés un plan con targets a nivel de archivo para aprobarlo antes de que se escriba nada
- Querés código de producción escrito en el codebase, dentro de un alcance que aprobaste
- Estás retomando un plan de una sesión anterior (lo volvés a aprobar en la compuerta)

## Por su cuenta

Apuntalo a código que ya existe y responde sin tocarlo:

```
/asdt-developer "¿cómo está armado el flujo de login hoy?"
/asdt-developer "¿cómo encararías migrar esto a la nueva API?"
/asdt-developer "code review del módulo de facturación"
```

Una pregunta se queda en exploración; pedirle un plan llega hasta el spec; solo pedirle que lo construya escribe archivos — y solo después de que apruebes.

## Posición en el pipeline

Típicamente corre **después del Arquitecto**, y lee cada resumen previo que exista: el del PM (la autoridad sobre los criterios de aceptación), el del Arquitecto (la decisión de diseño — no se reabre), el de UX/UI (los flujos, el diseño de pantallas y componentes, y la dirección visual, si la hay), el de Security (las mitigaciones que este cambio tiene que incluir) y el de QA (después de un `no-go`, los casos borde y de prueba que el arreglo tiene que cerrar). Todos son opcionales — sin ninguno, explora y especifica el problema por su cuenta y registra lo que asumió. En cambios simples el Arquitecto ni se invoca.

## Qué produce

El resumen del Developer (`developer`) — un solo resumen, escrito por etapas. `spec` escribe el plan (`Stage: spec`); `implement` lo reemplaza con lo construido (`Stage: implemented`) — los archivos cambiados, o los snippets de código en una corrida plan-only; `verify` agrega si los chequeos corrieron y pasaron. Una entrega anterior del mismo cambio nunca se pierde: cada plan y construcción nuevos arrastran sus archivos bajo `## Prior delivery`. Revisar código que ya existe corre `review` y escribe su resumen bajo `study-{topic}`.

Consumido por: **QA** (prueba contra el plan o contra lo construido, y nunca da por aprobado algo que el resumen no muestra), **Security** (la superficie planificada o el código que cambió) y el propio **Developer** cuando una corrida posterior lo retoma o itera sobre él.

## Patrones comunes

```
/asdt-developer "¿cómo agregarías exportación CSV al panel de reportes?"
# → Solo plan — el plan es el resumen, no se escribe nada
```

```
/asdt-developer "implementá el plan que aprobamos"
# → Vuelve a mostrar el plan (reconstruido si la sesión terminó), pide aprobación, lo construye y ofrece los chequeos
```

```
/asdt-developer "agregá exportación CSV al panel de reportes"
# → Construcción standalone — explora, especifica, te pide aprobación e implementa
```

## Límites — qué NO hace

- No produce decisiones de arquitectura ni ADRs
- No escribe specs de UX ni planes de prueba — los tests, cuando están activos, son código
- Nunca escribe un archivo en tu repo sin tu aprobación en la misma corrida
- Nunca escribe fuera de los archivos que declaró el spec aprobado — se detiene y reporta, también durante las rondas de corrección
- Nunca edita nada bajo `.asdt/` — eso lo escribe `/asdt-init`
- El paso `implement` nunca corre build, lint ni tests; los comandos corren solo en la compuerta de verificación, solo con tu sí y solo los que te mostró
