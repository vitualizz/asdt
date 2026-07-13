import type { SpecialistId } from './artifact-graph'

// Goal-based taxonomy (locked decision, docs-ux-overhaul): resolves the "8 by-goal
// categories, PROPOSAL pending sign-off" open item from ux-ui/component-spec —
// replaced with 4 goal buckets + 'all', no 5th bucket needed (all 14 recipes fit).
export type RecipeCategory = 'all' | 'from-scratch' | 'add-to-existing' | 'review-harden' | 'understand-document'

export interface RecipeCategoryDef {
  id: RecipeCategory
  labelEn: string
  labelEs: string
}

export const RECIPE_CATEGORIES: RecipeCategoryDef[] = [
  { id: 'all', labelEn: 'All', labelEs: 'Todas' },
  { id: 'from-scratch', labelEn: 'Start from scratch', labelEs: 'Empezar de cero' },
  { id: 'add-to-existing', labelEn: 'Add to existing code', labelEs: 'Añadir a código existente' },
  { id: 'review-harden', labelEn: 'Review & harden', labelEs: 'Revisar y endurecer' },
  { id: 'understand-document', labelEn: 'Understand & document', labelEs: 'Entender y documentar' },
]

export interface RecipeChip {
  specialistId: SpecialistId
}

export interface Recipe {
  id: string
  category: Exclude<RecipeCategory, 'all'>
  titleEn: string
  titleEs: string
  chips: RecipeChip[]
  /** Raw command sequence, in run order (routing-only `/asdt ...` calls included when present in the source doc). */
  commands: string[]
  noteEn: string
  noteEs: string
  kbNoteEn?: string
  kbNoteEs?: string
}

// Source of truth: src/content/docs/en/recipes.md (all 14 recipes, doc order preserved)
export const recipes: Recipe[] = [
  {
    id: 'ship-new-feature',
    category: 'from-scratch',
    titleEn: 'Ship a new user-facing feature',
    titleEs: 'Lanzar una feature nueva de cara al usuario',
    chips: [{ specialistId: 'pm' }, { specialistId: 'architect' }, { specialistId: 'developer' }],
    commands: [
      `/asdt "Add a contact form with email notification"`,
      `/asdt-pm "Add a contact form with email notification"`,
      `/asdt-architect`,
      `/asdt-developer`,
    ],
    noteEn: 'Start here when a feature request is in vague language and needs the full PM → Architect → Developer sequence.',
    noteEs: 'Empezá acá cuando el pedido de una feature está en lenguaje vago y necesita la secuencia completa PM → Architect → Developer.',
    kbNoteEn: 'Review the suggested specialist sequence, then run each in order.',
    kbNoteEs: 'Revisá la secuencia de especialistas sugerida y después corré cada uno en orden.',
  },
  {
    id: 'new-rest-endpoint',
    category: 'from-scratch',
    titleEn: 'Build a new REST API endpoint',
    titleEs: 'Construir un nuevo endpoint REST',
    chips: [{ specialistId: 'pm' }, { specialistId: 'architect' }, { specialistId: 'developer' }],
    commands: [
      `/asdt "Add a POST /api/v1/subscriptions endpoint"`,
      `/asdt-pm "Add a POST /api/v1/subscriptions endpoint"`,
      `/asdt-architect`,
      `/asdt-developer`,
    ],
    noteEn: 'For backend-only changes where the API contract is the primary design decision.',
    noteEs: 'Para cambios solo de backend donde el contrato de la API es la decisión de diseño principal.',
  },
  {
    id: 'new-screen-ux-first',
    category: 'from-scratch',
    titleEn: 'New screen with UX design first',
    titleEs: 'Pantalla nueva con diseño UX primero',
    chips: [
      { specialistId: 'pm' },
      { specialistId: 'ux-ui' },
      { specialistId: 'architect' },
      { specialistId: 'developer' },
      { specialistId: 'qa' },
    ],
    commands: [
      `/asdt "Add an onboarding wizard for new users"`,
      `/asdt-pm "Add an onboarding wizard for new users"`,
      `/asdt-ux-ui`,
      `/asdt-architect`,
      `/asdt-developer`,
      `/asdt-qa`,
    ],
    noteEn: 'For user-visible features where UI design should precede implementation.',
    noteEs: 'Para features de cara al usuario donde el diseño de UI debería preceder a la implementación.',
  },
  {
    id: 'security-review-before-shipping',
    category: 'review-harden',
    titleEn: 'Feature with security review before shipping',
    titleEs: 'Feature con revisión de seguridad antes de publicar',
    chips: [
      { specialistId: 'pm' },
      { specialistId: 'architect' },
      { specialistId: 'security' },
      { specialistId: 'developer' },
      { specialistId: 'qa' },
    ],
    commands: [
      `/asdt "Add OAuth login with GitHub"`,
      `/asdt-pm "Add OAuth login with GitHub"`,
      `/asdt-architect`,
      `/asdt-security`,
      `/asdt-developer`,
      `/asdt-qa`,
    ],
    noteEn: 'For anything touching auth, payments, PII (personally identifiable information), or external integrations.',
    noteEs: 'Para cualquier cosa que toque auth, pagos, PII (información de identificación personal) o integraciones externas.',
  },
  {
    id: 'explore-before-planning',
    category: 'understand-document',
    titleEn: 'Explore before planning (fuzzy problem)',
    titleEs: 'Explorar antes de planificar (problema difuso)',
    chips: [{ specialistId: 'researcher' }, { specialistId: 'pm' }],
    commands: [
      `/asdt-researcher "We're losing users at the signup step — what could we do?"`,
      `/asdt-pm "Based on the discovery brief: add progressive disclosure to the signup flow"`,
    ],
    noteEn: 'When the problem is unclear and you need discovery before requirements.',
    noteEs: 'Cuando el problema no está claro y necesitás descubrimiento antes de definir requisitos.',
    kbNoteEn: 'Researcher produces a discovery brief with a recommended direction — hand it to PM next.',
    kbNoteEs: 'Researcher produce un discovery brief con una dirección recomendada — pasáselo al PM después.',
  },
  {
    id: 'lock-scope-user-stories',
    category: 'from-scratch',
    titleEn: 'Lock scope and write user stories',
    titleEs: 'Fijar el alcance y escribir historias de usuario',
    chips: [{ specialistId: 'pm' }],
    commands: [`/asdt-pm "Add dark mode toggle to user settings"`],
    noteEn: 'When you have a clear feature idea and just need structured requirements.',
    noteEs: 'Cuando tenés una idea de feature clara y solo necesitás requisitos estructurados.',
  },
  {
    id: 'document-architecture-decision',
    category: 'understand-document',
    titleEn: 'Document an architecture decision',
    titleEs: 'Documentar una decisión de arquitectura',
    chips: [{ specialistId: 'architect' }],
    commands: [`/asdt-architect "Document the decision to use PostgreSQL row-level security for multi-tenancy"`],
    noteEn: 'When the technical approach needs a formal ADR and system design.',
    noteEs: 'Cuando el enfoque técnico necesita un ADR formal y un diseño de sistema.',
  },
  {
    id: 'write-code-from-settled-spec',
    category: 'add-to-existing',
    titleEn: 'Write production code from a settled spec',
    titleEs: 'Escribir código de producción a partir de una spec ya definida',
    chips: [{ specialistId: 'developer' }],
    commands: [`/asdt-developer "Implement the multi-tenancy RLS policy from the Architect's ADR"`],
    noteEn: 'When scope and architecture are already locked in the knowledge base.',
    noteEs: 'Cuando el alcance y la arquitectura ya están definidos en la base de conocimiento.',
  },
  {
    id: 'security-audit-existing-feature',
    category: 'review-harden',
    titleEn: 'Security audit on an existing feature',
    titleEs: 'Auditoría de seguridad sobre una feature existente',
    chips: [{ specialistId: 'security' }],
    commands: [`/asdt-security "Review the session management in the auth module"`],
    noteEn: 'Run at any point — no prior specialist run required.',
    noteEs: 'Corré en cualquier momento — no requiere que haya corrido un especialista antes.',
  },
  {
    id: 'design-new-ui-component',
    category: 'from-scratch',
    titleEn: 'Design a new UI component',
    titleEs: 'Diseñar un nuevo componente de UI',
    chips: [{ specialistId: 'ux-ui' }],
    commands: [`/asdt-ux-ui "Design a data table component with sorting, filtering, and pagination"`],
    noteEn: 'When you need a component spec before the developer starts coding.',
    noteEs: 'Cuando necesitás una spec de componente antes de que el developer empiece a codear.',
  },
  {
    id: 'validate-test-coverage',
    category: 'review-harden',
    titleEn: 'Validate test coverage before shipping',
    titleEs: 'Validar la cobertura de tests antes de publicar',
    chips: [{ specialistId: 'qa' }],
    commands: [`/asdt-qa`],
    noteEn: 'Run after Developer — QA reads the implementation artifact automatically.',
    noteEs: 'Corré después de Developer — QA lee la implementación automáticamente.',
  },
  {
    id: 'pickup-developer-existing-adr',
    category: 'add-to-existing',
    titleEn: 'Pick up at Developer after an existing ADR',
    titleEs: 'Retomar en Developer después de un ADR existente',
    chips: [{ specialistId: 'developer' }],
    commands: [`/asdt-developer`],
    noteEn: 'When PM and Architect ran in a previous session. Developer loads prior artifacts automatically.',
    noteEs: 'Cuando PM y Architect corrieron en una sesión anterior. Developer carga los artefactos previos automáticamente.',
    kbNoteEn: 'No need to re-run PM or Architect — artifacts are in the knowledge base.',
    kbNoteEs: 'No hace falta volver a correr PM o Architect — los artefactos están en la base de conocimiento.',
  },
  {
    id: 'add-security-review-inflight',
    category: 'review-harden',
    titleEn: 'Add a security review to an in-flight pipeline',
    titleEs: 'Agregar una revisión de seguridad a un pipeline en curso',
    chips: [{ specialistId: 'security' }],
    commands: [`/asdt-security`],
    noteEn: 'Run Security at any point without restarting the pipeline.',
    noteEs: 'Corré Security en cualquier momento sin reiniciar el pipeline.',
    kbNoteEn: 'Security reads system-design and dev-implementation from the knowledge base.',
    kbNoteEs: 'Security lee system-design y dev-implementation de la base de conocimiento.',
  },
  {
    id: 'qa-completed-feature-no-pipeline',
    category: 'review-harden',
    titleEn: 'QA a completed feature without a full pipeline',
    titleEs: 'Hacer QA de una feature terminada sin pipeline completo',
    chips: [{ specialistId: 'pm' }, { specialistId: 'qa' }],
    commands: [`/asdt-pm "Add acceptance criteria for the existing checkout flow"`, `/asdt-qa`],
    noteEn: 'When a feature was built without ASDT — run QA against the existing code.',
    noteEs: 'Cuando se construyó una feature sin ASDT — corré QA contra el código existente.',
  },
]
