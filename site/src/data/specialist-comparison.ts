import type { SpecialistId } from './artifact-graph'

export interface SpecialistComparisonRow {
  specialistId: SpecialistId
  exampleCommand: string
}

// Source of truth: src/content/docs/en/specialist-comparison.md (table order preserved)
export const specialistComparison: SpecialistComparisonRow[] = [
  {
    specialistId: 'pm',
    exampleCommand: '/asdt-pm "Add dark mode toggle to user settings"',
  },
  {
    specialistId: 'architect',
    exampleCommand: '/asdt-architect "Document the decision to use PostgreSQL row-level security for multi-tenancy"',
  },
  {
    specialistId: 'developer',
    exampleCommand: "/asdt-developer \"Implement the multi-tenancy RLS policy from the Architect's ADR\"",
  },
  {
    specialistId: 'qa',
    exampleCommand: '/asdt-qa "Review the checkout flow for edge cases"',
  },
  {
    specialistId: 'security',
    exampleCommand: '/asdt-security "Review the session management in the auth module"',
  },
  {
    specialistId: 'ux-ui',
    exampleCommand: '/asdt-ux-ui "Design a data table component with sorting, filtering, and pagination"',
  },
  {
    specialistId: 'researcher',
    exampleCommand: '/asdt-researcher "Explore ways to reduce onboarding drop-off"',
  },
]
