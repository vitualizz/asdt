---
title: Changelog
description: Changelog for ASDT releases. v0.4.2 added the Researcher specialist, interactive docs search, step-flow visualization, a full landing page, preset tiers, and live install progress.
order: 13
locale: en
---

# Changelog

## v0.4.2

- **Added** Researcher specialist (`/asdt-researcher`) — explores fuzzy problems and open-ended opportunities before requirements exist, producing a structured discovery brief with a single recommended direction.
- **Added** Interactive search to the docs site — find any page by keyword from anywhere in the docs using Fuse.js-powered search. Press `⌘K` to open.
- **Added** Step-flow visualization to specialist pages — each specialist page now shows exactly which steps it runs and at what complexity level, with an interactive tier switcher.
- **Added** Full landing page — live pipeline visualizer, specialist cards, Claude vs ASDT comparison, and a recipes tab.
- **Improved** Installer: stale skill files are now pruned on install, keeping `.asdt/skills/` clean and in sync with the installed version.
- **Added** Preset tiers for installer workflow configuration — choose a complexity profile during `asdt init`.
- **Added** DashboardStrings i18n and live install progress UI to the TUI dashboard — progress is displayed in the terminal with English and Spanish support.

---

To check your installed version, run `asdt --version` or inspect `.asdt/config.yaml`.
