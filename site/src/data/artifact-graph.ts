export type SpecialistId = 'researcher' | 'pm' | 'architect' | 'developer' | 'qa' | 'security' | 'ux-ui'

export interface ArtifactRef {
  key: string
  optional?: boolean
  consumedBy?: SpecialistId[]
  sentinel?: boolean
}

export interface SpecialistArtifacts {
  reads: ArtifactRef[]
  writes: ArtifactRef[]
}

// Every specialist returns exactly ONE summary, which the orchestrator writes to
// the session scratchpad at {scratchpad}/asdt/{change}/{role}.md. A key here is
// that summary's name as workflow.yaml `inputs:` writes it (the role id). Every
// read is optional: a specialist that finds nothing upstream works from the
// request and records the gap.
// Source of truth: skill/asdt-{id}/workflow.yaml
export const artifactGraph: Record<SpecialistId, SpecialistArtifacts> = {
  researcher: {
    reads: [{ key: 'Problem (raw)', sentinel: true }],
    writes: [{ key: 'researcher', consumedBy: ['pm', 'architect'] }],
  },
  pm: {
    reads: [
      { key: 'Request (raw)', sentinel: true },
      { key: 'researcher', optional: true },
    ],
    writes: [{ key: 'pm', consumedBy: ['ux-ui', 'architect', 'developer', 'qa'] }],
  },
  'ux-ui': {
    reads: [{ key: 'pm', optional: true }],
    writes: [{ key: 'ux-ui', consumedBy: ['architect', 'developer', 'qa'] }],
  },
  architect: {
    reads: [
      { key: 'pm', optional: true },
      { key: 'ux-ui', optional: true },
      { key: 'security', optional: true },
      { key: 'researcher', optional: true },
    ],
    writes: [{ key: 'architect', consumedBy: ['developer', 'qa', 'security'] }],
  },
  developer: {
    // `developer` is also read back by the Developer's own `spec` step: the
    // summary a new plan replaces, whose delivery it carries forward.
    reads: [
      { key: 'pm', optional: true },
      { key: 'architect', optional: true },
      { key: 'ux-ui', optional: true },
      { key: 'security', optional: true },
      { key: 'qa', optional: true },
      { key: 'developer', optional: true },
    ],
    writes: [{ key: 'developer', consumedBy: ['qa', 'security', 'developer'] }],
  },
  security: {
    reads: [
      { key: 'developer', optional: true },
      { key: 'architect', optional: true },
    ],
    writes: [{ key: 'security', consumedBy: ['architect', 'developer', 'qa'] }],
  },
  qa: {
    reads: [
      { key: 'pm', optional: true },
      { key: 'developer', optional: true },
      { key: 'architect', optional: true },
      { key: 'ux-ui', optional: true },
      { key: 'security', optional: true },
    ],
    writes: [{ key: 'qa', consumedBy: ['developer'] }],
  },
}

export const PIPELINE_ORDER: SpecialistId[] = ['researcher', 'pm', 'ux-ui', 'architect', 'developer', 'security', 'qa']

export const SPECIALIST_COLOR: Record<SpecialistId, string> = {
  researcher: '--c-res',
  pm: '--c-pm',
  architect: '--c-arch',
  developer: '--c-dev',
  qa: '--c-qa',
  security: '--c-sec',
  'ux-ui': '--c-ux',
}

export const SPECIALIST_LABEL: Record<SpecialistId, string> = {
  researcher: 'Researcher',
  pm: 'PM',
  architect: 'Architect',
  developer: 'Developer',
  qa: 'QA',
  security: 'Security',
  'ux-ui': 'UX/UI',
}
