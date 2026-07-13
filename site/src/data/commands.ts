import type { SpecialistId } from './artifact-graph'

export type CommandBadge = 'cli' | 'slash'

export interface CommandEntry {
  command: string
  badge: CommandBadge
  /** Only set for the 7 specialist-mapped slash commands — CLI + orchestrator commands (asdt-tui, /asdt, /asdt-init) use --color-accent instead. */
  specialistId?: SpecialistId
  titleEn: string
  titleEs: string
  oneLinerEn: string
  oneLinerEs: string
  /** Heading id in src/content/docs/{en,es}/commands.md, verified against the built page. */
  anchor: string
}

// Source of truth: src/content/docs/en/commands.md (1 CLI + 9 slash commands, doc order preserved)
export const commands: CommandEntry[] = [
  {
    command: 'asdt-tui',
    badge: 'cli',
    titleEn: 'Interactive terminal UI',
    titleEs: 'Interfaz interactiva de terminal',
    oneLinerEn: 'The only CLI tool — checks your setup and installs or updates the ASDT specialist skills.',
    oneLinerEs: 'La única herramienta CLI — revisa tu setup e instala o actualiza las skills de los especialistas ASDT.',
    anchor: 'asdt-tui',
  },
  {
    command: '/asdt-init',
    badge: 'slash',
    titleEn: 'Initialize ASDT',
    titleEs: 'Inicializar ASDT',
    oneLinerEn: 'Detects your stack and creates .asdt/config.yaml and .asdt/knowledge/platform.yaml.',
    oneLinerEs: 'Detecta tu stack y crea .asdt/config.yaml y .asdt/knowledge/platform.yaml.',
    anchor: 'asdt-init',
  },
  {
    command: '/asdt',
    badge: 'slash',
    titleEn: 'Pipeline routing suggestion',
    titleEs: 'Sugerencia de ruteo del pipeline',
    oneLinerEn: 'Analyzes the request and recommends which specialists to involve and in what order.',
    oneLinerEs: 'Analiza el pedido y recomienda qué especialistas involucrar y en qué orden.',
    anchor: 'asdt',
  },
  {
    command: '/asdt-researcher',
    badge: 'slash',
    specialistId: 'researcher',
    titleEn: 'Researcher',
    titleEs: 'Researcher',
    oneLinerEn: 'Runs the Researcher specialist only — discovery for fuzzy problems before requirements exist.',
    oneLinerEs: 'Corre solo al especialista Researcher — descubrimiento para problemas difusos antes de que existan requisitos.',
    anchor: 'asdt-researcher',
  },
  {
    command: '/asdt-pm',
    badge: 'slash',
    specialistId: 'pm',
    titleEn: 'Product Manager',
    titleEs: 'Product Manager',
    oneLinerEn: 'Runs the Product Manager specialist only.',
    oneLinerEs: 'Corre solo al especialista Product Manager.',
    anchor: 'asdt-pm',
  },
  {
    command: '/asdt-architect',
    badge: 'slash',
    specialistId: 'architect',
    titleEn: 'Architect',
    titleEs: 'Architect',
    oneLinerEn: 'Runs the Architect specialist only.',
    oneLinerEs: 'Corre solo al especialista Architect.',
    anchor: 'asdt-architect',
  },
  {
    command: '/asdt-developer',
    badge: 'slash',
    specialistId: 'developer',
    titleEn: 'Developer',
    titleEs: 'Developer',
    oneLinerEn: 'Runs the Developer specialist only.',
    oneLinerEs: 'Corre solo al especialista Developer.',
    anchor: 'asdt-developer',
  },
  {
    command: '/asdt-qa',
    badge: 'slash',
    specialistId: 'qa',
    titleEn: 'QA',
    titleEs: 'QA',
    oneLinerEn: 'Runs the QA specialist only.',
    oneLinerEs: 'Corre solo al especialista QA.',
    anchor: 'asdt-qa',
  },
  {
    command: '/asdt-security',
    badge: 'slash',
    specialistId: 'security',
    titleEn: 'Security',
    titleEs: 'Security',
    oneLinerEn: 'Runs the Security specialist only.',
    oneLinerEs: 'Corre solo al especialista Security.',
    anchor: 'asdt-security',
  },
  {
    command: '/asdt-ux-ui',
    badge: 'slash',
    specialistId: 'ux-ui',
    titleEn: 'UX/UI',
    titleEs: 'UX/UI',
    oneLinerEn: 'Runs the UX/UI specialist only.',
    oneLinerEs: 'Corre solo al especialista UX/UI.',
    anchor: 'asdt-ux-ui',
  },
]
