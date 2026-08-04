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

// Every specialist persists exactly ONE hand-off, at
// {project}/{change}/{role}/handoff. Every read is optional: a specialist that
// finds nothing upstream works from the request and records the gap.
// Source of truth: skill/asdt-{id}/workflow.yaml
export const artifactGraph: Record<SpecialistId, SpecialistArtifacts> = {
  researcher: {
    reads: [{ key: 'Problem (raw)', sentinel: true }],
    writes: [{ key: 'researcher/handoff', consumedBy: ['pm'] }],
  },
  pm: {
    reads: [
      { key: 'Request (raw)', sentinel: true },
      { key: 'researcher/handoff', optional: true },
    ],
    writes: [{ key: 'pm/handoff', consumedBy: ['architect', 'developer', 'qa', 'ux-ui'] }],
  },
  'ux-ui': {
    reads: [{ key: 'pm/handoff', optional: true }],
    writes: [{ key: 'ux-ui/handoff' }],
  },
  architect: {
    reads: [{ key: 'pm/handoff', optional: true }],
    writes: [{ key: 'architect/handoff', consumedBy: ['developer', 'qa', 'security'] }],
  },
  developer: {
    reads: [
      { key: 'pm/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
    ],
    writes: [{ key: 'developer/handoff', consumedBy: ['qa', 'security'] }],
  },
  security: {
    reads: [
      { key: 'developer/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
    ],
    writes: [{ key: 'security/handoff' }],
  },
  qa: {
    reads: [
      { key: 'pm/handoff', optional: true },
      { key: 'developer/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
    ],
    writes: [{ key: 'qa/handoff' }],
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
