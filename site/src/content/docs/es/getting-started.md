---
title: Comenzando
description: Instalá y ejecutá tu primer pipeline ASDT en minutos.
order: 1
locale: es
---

# Comenzando

## Requisitos

Antes de usar ASDT necesitás:

- **Claude Code** u **OpenCode** — instalado y autenticado
- **Un memory provider** — requerido para persistencia entre sesiones (por defecto: [Engram](https://github.com/Gentleman-Programming/engram))
- Una terminal (bash o zsh)

> **¿Compilando desde el código fuente?** Se requiere Go 1.22+. El instalador de una línea descarga un binario precompilado — no necesitás compilador.

## Instalación

```bash
curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash
```

Descarga el último binario precompilado (`asdt-tui`) para tu plataforma e instala en `~/.local/bin/`. No necesitás Go ni compilador.

> Si `~/.local/bin` no está en tu `PATH`, el instalador te muestra la línea exacta para agregar a tu shell profile.

## Instalá los especialistas

Ejecutá el instalador interactivo:

```bash
asdt-tui
```

Verifica tu setup, te deja elegir a qué asistente(s) apuntar (Claude Code, OpenCode, o ambos), y copia las skills de ASDT — cada especialista como una skill invocable de forma independiente.

## Inicializá tu proyecto

Abrí tu asistente de IA en la carpeta de tu proyecto y ejecutá:

```
/asdt-init
```

Detecta tu stack, hace un par de preguntas, y escribe `.asdt/config.yaml` con valores por defecto razonables.

## Tu primer pipeline

```
/asdt Agregar autenticación de usuario con email y contraseña
```

ASDT analiza la petición y recomienda una secuencia de especialistas — por ejemplo: `/asdt-pm` → `/asdt-architect` → `/asdt-developer`. Confirmá el plan y ejecutá cada comando. Cada especialista guarda su output en la base de conocimiento para que el siguiente retome donde dejó el anterior.

## Ejecutar especialistas individuales

```
/asdt-pm Agregar modo oscuro a la página de configuración
/asdt-architect Diseñar la estrategia de caché
/asdt-developer Implementar el componente de perfil de usuario
```
