---
title: Solución de problemas
description: 'Solucioná errores comunes de ASDT: command not found, problemas de PATH, Engram sin conectar, fallo de autenticación, especialista que no aparece, especialista equivocado.'
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

ASDT guarda los artefactos de los especialistas en Engram. Si Engram no está configurado, los especialistas no pueden guardar ni cargar artefactos entre pasos.

**Solución:**
1. Confirmá que Engram esté conectado a la configuración MCP (Model Context Protocol — el protocolo conector que permite que tu asistente hable con Engram; ver [Base de conocimiento y memoria](/asdt/docs/memory-and-engram) para entender cómo encaja todo) de tu asistente. Claude Code y OpenCode configuran los servidores MCP de forma distinta, así que seguí el setup de MCP del asistente que uses (y la [guía de setup de Engram](https://github.com/Gentleman-Programming/engram)). Ambos se conectan al mismo servidor de Engram; solo cambia la ubicación de la configuración.
2. Reiniciá tu asistente después de editar la configuración MCP — tanto Claude Code como OpenCode inician los servidores MCP configurados al arrancar.
3. Ejecutá `/asdt-init` dentro de tu asistente — escribe la configuración del memory provider en `.asdt/config.yaml`.
4. Revisá `.asdt/config.yaml` y confirmá que `memory.provider: engram` esté presente.

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

Si ASDT reporta que el modelo configurado no está disponible, tu `.asdt/config.yaml` puede referenciar un ID de modelo desactualizado.

**Solución:** Abrí `.asdt/config.yaml` y actualizá el campo `model` a un modelo que tu asistente pueda servir. Los IDs de modelo disponibles dependen de tu asistente y su provider configurado — para Claude Code, los modelos soportados están listados en la [documentación de Claude](https://docs.anthropic.com/en/docs/about-claude/models); para OpenCode, usá un ID de modelo expuesto por tu provider configurado. También podés elegir el preset **Chameleon** durante la instalación para quitar el campo `model:` por completo y dejar que cada asistente use su propio default.

## Problemas de especialistas

### Especialista equivocado — cómo recuperarse

Detené la corrida actual. No se pierde nada — los artefactos previos permanecen en la base de conocimiento. Invocá el especialista correcto directamente.

**Ejemplo:** Si corriste `/asdt-developer` antes de producir un ADR, ejecutá `/asdt-architect` para crear el registro de decisión. Luego volvé a correr `/asdt-developer` — lee el ADR automáticamente.

Ver [Comparación de especialistas](/asdt/docs/specialist-comparison) para elegir el especialista correcto para tu escenario.

### El comando del especialista no aparece en tu asistente

Si al escribir `/asdt-pm` (o cualquier especialista) no aparece el autocompletado, los archivos de skill no están instalados.

**Solución:**
1. Ejecutá `asdt-tui` y (re)instalá las skills para tu asistente. Esto instala en `~/.claude/skills` para Claude Code, o en `~/.config/opencode/skills` (más los wrappers de comandos en `~/.config/opencode/commands/`) para OpenCode.
2. Reiniciá tu asistente — tanto Claude Code como OpenCode recargan las definiciones de skills (y comandos) al iniciar.
3. Si el problema persiste, ejecutá `asdt-tui` de nuevo y reinstalá los archivos de skill.

### Los artefactos no cargan en la siguiente sesión

Cada especialista lee los artefactos previos desde Engram. Si una sesión nueva no encuentra los artefactos de la anterior, Engram no estaba corriendo durante la sesión previa cuando se guardaron los artefactos.

**Solución:** Asegurate de que Engram (vía MCP) esté corriendo antes de invocar cualquier especialista. Revisá el estado del servidor MCP en tu asistente. Los artefactos solo se persisten si Engram está activo en el momento en que el especialista los guarda.

## Limitaciones conocidas

- **Memory provider requerido.** ASDT requiere una instancia de Engram corriendo (vía MCP) para persistir artefactos entre corridas de especialistas. No hay almacenamiento de respaldo — si Engram no está conectado, los artefactos no se guardan y el siguiente especialista del pipeline no encontrará sus inputs.
- **Claude Code u OpenCode requerido.** Los especialistas de ASDT son comandos slash invocados dentro de un asistente soportado (Claude Code u OpenCode). No corren en una interfaz de chat estándar ni vía una API de modelo directamente.
- **Solo macOS y Linux.** El script de instalación apunta a macOS y Linux (x86_64 y arm64). Windows vía WSL2 no está probado y no está soportado en esta versión.
- **Una sesión de pipeline activa a la vez.** Correr dos pipelines de especialistas simultáneamente en el mismo directorio de proyecto puede causar colisiones de claves de artefactos en Engram.
