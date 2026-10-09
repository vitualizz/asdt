---
title: QA Engineer
description: Construye la red de seguridad antes de que el código salga a producción — planes de prueba, validación de criterios de aceptación, análisis de casos borde e informes de calidad — el especialista a invocar cuando "funciona en mi máquina" no es suficiente.
order: 23
locale: es
---

# QA Engineer (`/asdt-qa`)

> Construye la red de seguridad antes de que el código salga a producción — planes de prueba, validación de criterios de aceptación, análisis de casos borde e informes de calidad — el especialista a invocar cuando "funciona en mi máquina" no es suficiente.

## Qué hace

El especialista QA encuentra lo que los criterios de aceptación pasaron por alto y lo convierte en un plan de pruebas con un veredicto go/no-go. Corre un solo paso, `test-plan`, en este orden:

1. **Brechas en los ACs** — cada criterio heredado se juzga por atomicidad, medibilidad y caso negativo. Un criterio que ningún test podría observar es una brecha bloqueante.
2. **Casos borde** — el trabajo de verdad: casos de input, estado, concurrencia y falla de dependencias que los criterios nunca mencionaron. Si corrió UX/UI, cada bifurcación de flujo y cada estado vacío, de carga y de error que nombró es un caso candidato.
3. **Estrategia** — el reparto unitario / integración / e2e para este cambio, en tres líneas.
4. **Casos de prueba** — Given/When/Then; si hay hand-off del Developer, cada uno apunta al archivo construido o planificado que ejercita. Si corrió Security, cada hallazgo recibe un caso que prueba que su mitigación se sostiene; una mitigación que ningún test puede observar se convierte en un chequeo que podés correr vos.
5. **Veredicto** — `go` o `no-go`, con dos líneas de por qué. Una brecha bloqueante en los ACs, un camino crítico sin cubrir o — cuando hay hand-off del Developer — un `high` de Security cuya mitigación no deja rastro en sus archivos (los archivos cambiados una vez construido, los planificados si es un plan) es un `no-go`. Sobre un plan todavía no hay nada construido, así que el veredicto es `no-go — not built yet`, y el razonamiento dice si el plan está listo para construirse.

QA no ejecuta nada. Nunca reporta un pass o un fail sobre algo que no se corrió: un objetivo NFR se convierte en un comando que podés correr para medirlo.

## Cuándo invocarlo

- El código está listo para revisión y necesitas un quality gate antes de que salga a producción
- Los criterios de aceptación existen pero no han sido validados formalmente (atomicidad, mensurabilidad, independencia)
- Querés cobertura sistemática de casos borde, no solo tests del happy path
- Necesitás un plan de pruebas estructurado que un Developer pueda implementar sin adivinar

## Por su cuenta

No necesita que alguien acabe de programar. Apúntalo a lo que ya está:

```
/asdt-qa "¿qué no cubren nuestros tests de auth?"
/asdt-qa "revisa la cobertura del carrito y dame un veredicto"
```

Trabaja con lo que encuentre —hand-offs previos si los hay, el código si no— y cierra igual con go/no-go.

## Posición en el pipeline

Típicamente corre **después del Developer** y es el sign-off final antes de mergear. Lee cada hand-off que exista: `pm/handoff` (criterios de aceptación y objetivos NFR), `developer/handoff` (el plan, o lo construido y si sus chequeos pasaron), `architect/handoff` (el diseño y sus riesgos declarados), `ux-ui/handoff` (bifurcaciones y copy de los flujos, estados y extremos de datos de las pantallas, y los estados disabled, loading y error de los componentes diseñados) y `security/handoff` (las mitigaciones a probar). Puede correr antes — contra el hand-off del PM o un plan guardado del Developer — para detectar problemas en los criterios de aceptación antes de que empiece la implementación. Ese pase temprano ahorra mucho más que encontrar las brechas con el código ya escrito. Toda entrada es opcional: sin ninguna, trabaja desde la petición y el código.

## Qué produce

`qa/handoff` — brechas en los ACs, casos borde, estrategia, casos de prueba, los chequeos que te ofrece y el veredicto, como UN solo artefacto. Auditar una suite existente corre `review` y guarda `{project}/study/{topic}/qa`.

Ningún especialista lo declara como input: es el registro de sign-off. Después de un `no-go`, el reporte cierra proponiendo a quien arregla lo encontrado — normalmente el Developer.

## Patrones comunes

```
/asdt-qa Revisar el flujo de checkout en busca de casos borde
# → El happy path está testeado pero las condiciones límite y las rutas de error necesitan cobertura
```

```
/asdt-qa Validar criterios de aceptación antes de que empiece la implementación
# → Correr QA contra pm/handoff para detectar problemas en los ACs temprano
```

```
/asdt-qa Construir un plan de pruebas para el módulo de autenticación
# → Estrategia completa de pirámide de testing para código sensible a seguridad
```

## Límites — qué NO hace

- No escribe código de implementación
- No escribe decisiones de arquitectura ni specs de UX
- Nunca afirma un pass o un fail sobre algo que no se corrió — no ejecuta nada
- Un plan que solo reformula los criterios de aceptación como tests no agregó nada — los casos borde son el entregable
- Los casos de prueba son especificaciones (Given/When/Then) — no código ejecutable
