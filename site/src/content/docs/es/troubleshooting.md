---
title: Solución de problemas
description: 'Solucioná errores comunes de ASDT: command not found, problemas de PATH, Engram sin conectar, fallo de autenticación, especialista que no aparece, cambiar de especialista.'
order: 11
locale: es
---

# Solución de problemas

## Problemas de instalación

### `command not found: asdt-tui` después de instalar

El script de instalación coloca el binario en `~/.local/bin/asdt-tui`. Si tu shell no lo encuentra, tu PATH no incluye ese directorio.

**Solución:**

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Para que sea permanente, agregá la línea de arriba a tu `~/.bashrc`, `~/.zshrc` o `~/.profile`, y luego reiniciá tu terminal o ejecutá `source ~/.zshrc`.

Verificá la instalación abriendo el menú:

```bash
asdt-tui
```

### El script de instalación falla en silencio

Si el script de instalación termina sin imprimir una versión, ejecutalo con salida de depuración:

```bash
bash -x <(curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh)
```

El flag `-x` imprime cada comando a medida que se ejecuta. Buscá la primera línea que falla. Causas comunes: `~/.local/bin` no existe (solución: `mkdir -p ~/.local/bin`), o la descarga fue bloqueada por un firewall.

## Problemas del memory provider

### Engram no está conectado

ASDT recuerda el conocimiento clave del proyecto (decisiones, convenciones, lecciones) en Engram. Si Engram no está conectado, los especialistas igual corren: cada uno avisa en una línea que la memoria no está conectada, trabaja desde el scratchpad de la sesión y no consulta ni recuerda nada, así que nada pasa a la sesión siguiente.

Si instalaste con **Sin memoria**, primero volvé a ejecutar `asdt-tui` y elegí Engram en el paso de memoria.

**Solución:**
1. Confirmá que Engram esté conectado a la configuración MCP (Model Context Protocol — el protocolo conector que permite que tu asistente hable con Engram; ver [Scratchpad y memoria](/asdt/es/docs/memory-and-engram) para entender cómo encaja todo) de tu asistente. Claude Code y OpenCode configuran los servidores MCP de forma distinta, así que seguí el setup de MCP del asistente que uses (y la [guía de setup de Engram](https://github.com/Gentleman-Programming/engram)). Ambos se conectan al mismo servidor de Engram; solo cambia la ubicación de la configuración.
2. Reiniciá tu asistente después de editar la configuración MCP — tanto Claude Code como OpenCode inician los servidores MCP configurados al arrancar.
3. Ejecutá `/asdt-init` dentro de tu asistente — escribe la configuración del memory provider en `.asdt/config.yaml`.
4. Revisá `.asdt/config.yaml` y confirmá que `memory.provider: engram` esté presente. Un proyecto que nombra un provider distinto al de tu instalación se trata como si no tuviera memoria.

### Conexión rechazada desde Engram

Si Engram reporta un error de conexión, el proceso del servidor MCP no está corriendo.

**Solución:** Reiniciá tu asistente — tanto Claude Code como OpenCode inician los servidores MCP configurados al arrancar. Si el error persiste, revisá los logs de tu asistente en busca de errores de arranque de MCP.

## Problemas del asistente (Claude Code / OpenCode)

### Fallo de autenticación al correr un especialista

Si un comando de especialista devuelve un error de autenticación, tu sesión del asistente expiró.

**Solución:** Volvé a autenticarte con el flujo de login de tu asistente y luego reinicialo:

- **Claude Code** — iniciá sesión de nuevo con la CLI `claude` (`claude auth login`).
- **OpenCode** — iniciá sesión de nuevo con la CLI `opencode` (`opencode auth login`).

Seguí las indicaciones para volver a autenticarte y luego reiniciá el asistente.

### Modelo no disponible

Si tu asistente reporta que un modelo no está disponible, los pasos de especialista o los agentes ejecutores instalados pueden referenciar un ID de modelo que no puede servir. Los modelos se fijan al instalar, no en `.asdt/config.yaml`.

**Solución:** Volvé a ejecutar `asdt-tui` y elegí otro preset de modelos, o fijá vos el modelo de cada paso. Los IDs de modelo disponibles dependen de tu asistente y su provider configurado — para Claude Code, los modelos soportados están listados en la [documentación de Claude](https://docs.anthropic.com/en/docs/about-claude/models); para OpenCode, usá un ID de modelo expuesto por tu provider configurado. También puedes elegir el preset **Chameleon** para quitar el campo `model:` por completo y dejar que cada asistente use su propio default.

## Problemas de especialistas

### Cambiar de especialista

No existe un especialista equivocado del que tengas que recuperarte — cada uno trabaja de forma independiente y lee los resúmenes que los especialistas anteriores dejaron en el scratchpad de la sesión, así que cambiar no te cuesta nada. Si una petición no encaja con el especialista que invocaste, simplemente detené la corrida y ejecutá otro directamente; no se pierde nada.

**Ejemplo:** Si corriste `/asdt-developer` cuando todavía no había un diseño, ejecutá `/asdt-architect` para tomar la decisión. Luego volvé a correr `/asdt-developer` con el mismo pedido — recibe el resumen del Arquitecto como input.

Si no estás seguro de qué especialista se ajusta a tu petición, ejecutá `/asdt` — recomienda uno más adecuado cuando tu petición no encaja. Ver [Especialistas](/asdt/es/docs/specialists) para elegir el especialista correcto para tu escenario.

### El comando del especialista no aparece en tu asistente

Si al escribir `/asdt-pm` (o cualquier especialista) no aparece el autocompletado, los archivos de skill no están instalados.

**Solución:**
1. Ejecutá `asdt-tui` y (re)instalá las skills para tu asistente. Esto instala en `~/.claude/skills` para Claude Code, o en `~/.config/opencode/skills` (más los wrappers de comandos en `~/.config/opencode/commands/`) para OpenCode.
2. Reiniciá tu asistente — tanto Claude Code como OpenCode recargan las definiciones de skills (y comandos) al iniciar.
3. Si el problema persiste, ejecutá `asdt-tui` de nuevo y reinstalá los archivos de skill.

### Falta trabajo anterior en la siguiente sesión

Es lo esperable para los archivos de trabajo: el resumen de cada especialista vive en el scratchpad de la sesión, y el scratchpad pertenece a la sesión. Una sesión nueva reconstruye lo que necesita a partir de git, de las notas que recuerda la memoria y de vos, y cada gate (como la aprobación del plan del Developer) vuelve a preguntar.

Lo que sí debería pasar de una sesión a otra es el conocimiento clave: decisiones, conclusiones de debates, convenciones, lecciones, hallazgos de auditorías. Si una sesión nueva no recuerda nada de eso, la memoria no estaba conectada cuando se lo habría recordado.

**Solución:** Asegurate de que Engram (vía MCP) esté corriendo antes de invocar un especialista; cada especialista avisa en una línea cuando no lo está. Las notas solo se recuerdan mientras la memoria está conectada. Para lo que se perdió, volvé a contárselo al especialista, o decí "acordate de que…" para guardarlo.

## Limitaciones conocidas

- **La memoria es lo que cruza sesiones.** Sin un memory provider conectado, ASDT funciona dentro de una sesión, pero ningún conocimiento pasa a la siguiente: tenés que repetir el contexto que habría aportado una consulta a la memoria.
- **Claude Code u OpenCode requerido.** Los especialistas de ASDT son comandos slash invocados dentro de un asistente soportado (Claude Code u OpenCode). No corren en una interfaz de chat estándar ni vía una API de modelo directamente.
- **Solo macOS y Linux.** El script de instalación apunta a macOS y Linux (x86_64 y arm64). Windows vía WSL2 no está probado y no está soportado en esta versión.
- **Las sesiones no comparten archivos de trabajo.** Cada sesión tiene su propio scratchpad, así que una segunda sesión sobre el mismo cambio no ve los resúmenes de la primera: reconstruye lo que necesita a partir de git, la memoria y vos.
