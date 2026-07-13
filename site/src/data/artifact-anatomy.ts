export interface ArtifactAnatomyField {
  id: string
  label: string
  value: string
}

export const artifactAnatomy: ArtifactAnatomyField[] = [
  {
    id: 'title',
    label: 'Title',
    value: 'add-auth/developer/dev-spec',
  },
  {
    id: 'topicKey',
    label: 'Topic Key',
    value: 'asdt/add-auth/developer/dev-spec',
  },
  {
    id: 'type',
    label: 'Type',
    value: 'decision',
  },
  {
    id: 'project',
    label: 'Project',
    value: 'asdt',
  },
]
