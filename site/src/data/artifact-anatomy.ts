export interface ArtifactAnatomyField {
  id: string
  label: string
  value: string
}

// The anatomy of one remembered note — the only thing ASDT keeps in memory.
// Source of truth: skill/asdt-core/protocol.md §2 ("The note").
export const artifactAnatomy: ArtifactAnatomyField[] = [
  {
    id: 'subject',
    label: 'Subject',
    value: 'session-storage',
  },
  {
    id: 'kind',
    label: 'Kind',
    value: 'decision',
  },
  {
    id: 'source',
    label: 'Source',
    value: 'architect on add-auth, 2026-05-04',
  },
  {
    id: 'project',
    label: 'Project',
    value: 'asdt-1a2b3c4d',
  },
]
