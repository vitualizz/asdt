import type { SpecialistId } from './artifact-graph'

export interface SpecialistComparisonRow {
  specialistId: SpecialistId
  teaserEn: string
  teaserEs: string
  invokeWhenEn: string
  invokeWhenEs: string
  producesEn: string
  producesEs: string
  doNotUseWhenEn: string
  doNotUseWhenEs: string
  exampleCommand: string
}

// Source of truth: src/content/docs/en/specialist-comparison.md (table order preserved)
export const specialistComparison: SpecialistComparisonRow[] = [
  {
    specialistId: 'pm',
    teaserEn: 'Locks scope and turns vague requests into structured user stories.',
    teaserEs: 'Fija el alcance y convierte pedidos vagos en historias de usuario estructuradas.',
    invokeWhenEn: "The request is vague or user-facing, scope isn't locked, or user stories don't exist yet",
    invokeWhenEs: 'El pedido es vago o de cara al usuario, el alcance no está definido, o todavía no existen historias de usuario',
    producesEn: 'pm/backlog-entry — feature summary, ordered user stories with AC, in/out scope, risk flags',
    producesEs: 'pm/backlog-entry — resumen de la feature, historias de usuario ordenadas con AC, alcance dentro/fuera, riesgos',
    doNotUseWhenEn: 'You already have a clear backlog entry — re-running PM regenerates stories from scratch',
    doNotUseWhenEs: 'Ya tenés un backlog entry claro — volver a correr PM regenera las historias desde cero',
    exampleCommand: '/asdt-pm "Add dark mode toggle to user settings"',
  },
  {
    specialistId: 'architect',
    teaserEn: 'Decides how the pieces fit together before anyone writes code.',
    teaserEs: 'Decide cómo encajan las piezas antes de escribir código.',
    invokeWhenEn: 'The solution touches service boundaries, data models, or API contracts, or has two viable technical approaches worth documenting',
    invokeWhenEs: 'La solución toca límites de servicios, modelos de datos o contratos de API, o tiene dos enfoques técnicos viables que vale la pena documentar',
    producesEn: 'architectural-decision (ADR) + system-design — data model, API surface, service boundaries',
    producesEs: 'architectural-decision (ADR) + system-design — modelo de datos, superficie de API, límites de servicios',
    doNotUseWhenEn: 'You need implementation code, test plans, or UX specs — Architect produces decisions, not code',
    doNotUseWhenEs: 'Necesitás código de implementación, planes de pruebas o specs de UX — Architect produce decisiones, no código',
    exampleCommand: '/asdt-architect "Document the decision to use PostgreSQL row-level security for multi-tenancy"',
  },
  {
    specialistId: 'developer',
    teaserEn: 'Turns a settled design into an ordered implementation plan or real code.',
    teaserEs: 'Convierte un diseño ya definido en un plan de implementación ordenado o código real.',
    invokeWhenEn: 'The shape of the solution is settled and you need an ordered implementation plan or production code written to the repo',
    invokeWhenEs: 'La forma de la solución está definida y necesitás un plan de implementación ordenado o código de producción escrito en el repo',
    producesEn: 'developer/dev-implementation — ordered file manifest and code plan',
    producesEs: 'developer/dev-implementation — manifiesto de archivos ordenado y plan de código',
    doNotUseWhenEn: "You haven't locked scope or architecture yet — Developer will implement against ambiguous requirements",
    doNotUseWhenEs: 'Todavía no fijaste el alcance ni la arquitectura — Developer va a implementar contra requisitos ambiguos',
    exampleCommand: "/asdt-developer \"Implement the multi-tenancy RLS policy from the Architect's ADR\"",
  },
  {
    specialistId: 'qa',
    teaserEn: 'Turns acceptance criteria into a systematic test plan with a go/no-go verdict.',
    teaserEs: 'Convierte los criterios de aceptación en un plan de pruebas sistemático con un veredicto de go/no-go.',
    invokeWhenEn: "Code is ready for review, AC exists but hasn't been validated, or you need systematic edge case and boundary coverage",
    invokeWhenEs: 'El código está listo para revisión, existen AC pero no fueron validados, o necesitás cobertura sistemática de edge cases y límites',
    producesEn: 'test-plan — AC coverage %, uncovered gaps, full Given/When/Then test case list, quality verdict',
    producesEs: 'test-plan — % de cobertura de AC, gaps sin cubrir, lista completa de test cases Given/When/Then, veredicto de calidad',
    doNotUseWhenEn: 'You want executable test code written — QA produces test specifications, not runnable code',
    doNotUseWhenEs: 'Querés código de test ejecutable — QA produce especificaciones de test, no código que corra',
    exampleCommand: '/asdt-qa "Review the checkout flow for edge cases"',
  },
  {
    specialistId: 'security',
    teaserEn: 'Hunts for auth, data, and integration risks before they ship.',
    teaserEs: 'Busca riesgos de auth, datos e integraciones antes de que se publiquen.',
    invokeWhenEn: 'The feature touches auth, sessions, PII, external integrations, webhooks, or new public API endpoints',
    invokeWhenEs: 'La feature toca auth, sesiones, PII, integraciones externas, webhooks, o nuevos endpoints públicos de API',
    producesEn: 'security-findings (severity-rated, CWE-referenced) + hardening-checklist — must-fix vs can-defer',
    producesEs: 'security-findings (con severidad y CWE) + hardening-checklist — qué arreglar sí o sí vs qué se puede postergar',
    doNotUseWhenEn: 'You want implementation code or architectural decisions — Security produces findings and checklists only',
    doNotUseWhenEs: 'Querés código de implementación o decisiones arquitectónicas — Security solo produce findings y checklists',
    exampleCommand: '/asdt-security "Review the session management in the auth module"',
  },
  {
    specialistId: 'ux-ui',
    teaserEn: 'Maps flows and components before implementation starts.',
    teaserEs: 'Mapea flujos y componentes antes de que arranque la implementación.',
    invokeWhenEn: 'A new screen or feature-level UI needs design before implementation begins, or user flows need mapping',
    invokeWhenEs: 'Una pantalla nueva o una UI a nivel de feature necesita diseño antes de implementar, o hay que mapear flujos de usuario',
    producesEn: 'ux-brief (flows, IA, success criteria) + component-spec — inventory of reused/extended/new components',
    producesEs: 'ux-brief (flujos, IA, criterios de éxito) + component-spec — inventario de componentes reusados/extendidos/nuevos',
    doNotUseWhenEn: 'The screen has already been built — a UX spec delivered after implementation is too late to shape it',
    doNotUseWhenEs: 'La pantalla ya está construida — una spec de UX entregada después de implementar llega demasiado tarde para darle forma',
    exampleCommand: '/asdt-ux-ui "Design a data table component with sorting, filtering, and pagination"',
  },
  {
    specialistId: 'researcher',
    teaserEn: 'Explores a fuzzy problem before anyone commits to a direction.',
    teaserEs: 'Explora un problema difuso antes de comprometerse con una dirección.',
    invokeWhenEn: 'The problem is fuzzy or open-ended — you need discovery and framing before requirements can be written',
    invokeWhenEs: 'El problema es difuso o abierto — necesitás descubrimiento y enmarcado antes de poder escribir requisitos',
    producesEn: 'researcher/discovery-brief — problem framing, divergent idea set, feasibility scan, single recommended direction',
    producesEs: 'researcher/discovery-brief — enmarcado del problema, conjunto divergente de ideas, análisis de factibilidad, una dirección recomendada',
    doNotUseWhenEn: 'You already have a defined problem statement — Researcher explores; it does not produce user stories or ADRs',
    doNotUseWhenEs: 'Ya tenés un problema bien definido — Researcher explora; no produce historias de usuario ni ADRs',
    exampleCommand: '/asdt-researcher "Explore ways to reduce onboarding drop-off"',
  },
]
