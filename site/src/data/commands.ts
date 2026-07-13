import type { SpecialistId } from './artifact-graph'

export type CommandBadge = 'cli' | 'slash'

export interface CommandEntry {
  command: string
  badge: CommandBadge
  /** Only set for the 7 specialist-mapped slash commands — CLI + orchestrator commands (asdt-tui, /asdt, /asdt-init) use --color-accent instead. */
  specialistId?: SpecialistId
  /** Heading id in src/content/docs/{en,es}/commands.md, verified against the built page. Also the key into t.data.commands. */
  anchor: string
}

// Source of truth: src/content/docs/en/commands.md (1 CLI + 9 slash commands, doc order preserved)
export const commands: CommandEntry[] = [
  {
    command: 'asdt-tui',
    badge: 'cli',
    anchor: 'asdt-tui',
  },
  {
    command: '/asdt-init',
    badge: 'slash',
    anchor: 'asdt-init',
  },
  {
    command: '/asdt',
    badge: 'slash',
    anchor: 'asdt',
  },
  {
    command: '/asdt-researcher',
    badge: 'slash',
    specialistId: 'researcher',
    anchor: 'asdt-researcher',
  },
  {
    command: '/asdt-pm',
    badge: 'slash',
    specialistId: 'pm',
    anchor: 'asdt-pm',
  },
  {
    command: '/asdt-architect',
    badge: 'slash',
    specialistId: 'architect',
    anchor: 'asdt-architect',
  },
  {
    command: '/asdt-developer',
    badge: 'slash',
    specialistId: 'developer',
    anchor: 'asdt-developer',
  },
  {
    command: '/asdt-qa',
    badge: 'slash',
    specialistId: 'qa',
    anchor: 'asdt-qa',
  },
  {
    command: '/asdt-security',
    badge: 'slash',
    specialistId: 'security',
    anchor: 'asdt-security',
  },
  {
    command: '/asdt-ux-ui',
    badge: 'slash',
    specialistId: 'ux-ui',
    anchor: 'asdt-ux-ui',
  },
]
