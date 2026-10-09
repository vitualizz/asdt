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
    writes: [{ key: 'researcher/handoff', consumedBy: ['pm', 'architect'] }],
  },
  pm: {
    reads: [
      { key: 'Request (raw)', sentinel: true },
      { key: 'researcher/handoff', optional: true },
    ],
    writes: [{ key: 'pm/handoff', consumedBy: ['ux-ui', 'architect', 'developer', 'qa'] }],
  },
  'ux-ui': {
    reads: [{ key: 'pm/handoff', optional: true }],
    writes: [{ key: 'ux-ui/handoff', consumedBy: ['architect', 'developer', 'qa'] }],
  },
  architect: {
    reads: [
      { key: 'pm/handoff', optional: true },
      { key: 'ux-ui/handoff', optional: true },
      { key: 'security/handoff', optional: true },
      { key: 'researcher/handoff', optional: true },
    ],
    writes: [{ key: 'architect/handoff', consumedBy: ['developer', 'qa', 'security'] }],
  },
  developer: {
    // developer/handoff is also read back by the Developer itself when a later
    // run resumes the persisted plan ("implement the plan we approved").
    reads: [
      { key: 'pm/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
      { key: 'ux-ui/handoff', optional: true },
      { key: 'security/handoff', optional: true },
      { key: 'developer/handoff', optional: true },
    ],
    writes: [{ key: 'developer/handoff', consumedBy: ['qa', 'security', 'developer'] }],
  },
  security: {
    reads: [
      { key: 'developer/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
    ],
    writes: [{ key: 'security/handoff', consumedBy: ['architect', 'developer', 'qa'] }],
  },
  qa: {
    reads: [
      { key: 'pm/handoff', optional: true },
      { key: 'developer/handoff', optional: true },
      { key: 'architect/handoff', optional: true },
      { key: 'ux-ui/handoff', optional: true },
      { key: 'security/handoff', optional: true },
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
