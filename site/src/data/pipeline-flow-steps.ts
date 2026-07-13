import type { SpecialistId } from './artifact-graph'

export interface FlowStep {
  id: string
  command: string
  specialistId: SpecialistId | 'asdt'
}

export interface PipelineFlow {
  id: string
  steps: FlowStep[]
}

export const pipelineFlows: PipelineFlow[] = [
  {
    id: 'full-feature',
    steps: [
      {
        id: 'orchestrate',
        command: '/asdt Add passwordless login with magic links',
        specialistId: 'asdt',
      },
      {
        id: 'pm',
        command: '/asdt-pm',
        specialistId: 'pm',
      },
      {
        id: 'architect',
        command: '/asdt-architect',
        specialistId: 'architect',
      },
      {
        id: 'developer',
        command: '/asdt-developer',
        specialistId: 'developer',
      },
      {
        id: 'security',
        command: '/asdt-security',
        specialistId: 'security',
      },
    ],
  },
  {
    id: 'mid-pipeline',
    steps: [
      {
        id: 'developer',
        command: "/asdt-developer Implement based on the Architect's ADR",
        specialistId: 'developer',
      },
      {
        id: 'qa',
        command: '/asdt-qa',
        specialistId: 'qa',
      },
    ],
  },
]
