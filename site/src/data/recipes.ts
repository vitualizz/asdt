import type { SpecialistId } from './artifact-graph'

// Goal-based taxonomy (locked decision, docs-ux-overhaul): resolves the "8 by-goal
// categories, PROPOSAL pending sign-off" open item from ux-ui/component-spec —
// replaced with 4 goal buckets + 'all', no 5th bucket needed (all 14 recipes fit).
export type RecipeCategory = 'all' | 'from-scratch' | 'add-to-existing' | 'review-harden' | 'understand-document'

export interface RecipeCategoryDef {
  id: RecipeCategory
}

export const RECIPE_CATEGORIES: RecipeCategoryDef[] = [
  { id: 'all' },
  { id: 'from-scratch' },
  { id: 'add-to-existing' },
  { id: 'review-harden' },
  { id: 'understand-document' },
]

export interface RecipeChip {
  specialistId: SpecialistId
}

export interface Recipe {
  id: string
  category: Exclude<RecipeCategory, 'all'>
  chips: RecipeChip[]
  /** Raw command sequence, in run order (routing-only `/asdt ...` calls included when present in the source doc). */
  commands: string[]
}

// Source of truth: src/content/docs/en/recipes.md (all 14 recipes, doc order preserved)
export const recipes: Recipe[] = [
  {
    id: 'ship-new-feature',
    category: 'from-scratch',
    chips: [{ specialistId: 'pm' }, { specialistId: 'architect' }, { specialistId: 'developer' }],
    commands: [
      `/asdt "Add a contact form with email notification"`,
      `/asdt-pm "Add a contact form with email notification"`,
      `/asdt-architect`,
      `/asdt-developer`,
    ],
  },
  {
    id: 'new-rest-endpoint',
    category: 'from-scratch',
    chips: [{ specialistId: 'pm' }, { specialistId: 'architect' }, { specialistId: 'developer' }],
    commands: [
      `/asdt "Add a POST /api/v1/subscriptions endpoint"`,
      `/asdt-pm "Add a POST /api/v1/subscriptions endpoint"`,
      `/asdt-architect`,
      `/asdt-developer`,
    ],
  },
  {
    id: 'new-screen-ux-first',
    category: 'from-scratch',
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
  },
  {
    id: 'security-review-before-shipping',
    category: 'review-harden',
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
  },
  {
    id: 'explore-before-planning',
    category: 'understand-document',
    chips: [{ specialistId: 'researcher' }, { specialistId: 'pm' }],
    commands: [
      `/asdt-researcher "We're losing users at the signup step — what could we do?"`,
      `/asdt-pm "Based on the discovery brief: add progressive disclosure to the signup flow"`,
    ],
  },
  {
    id: 'lock-scope-user-stories',
    category: 'from-scratch',
    chips: [{ specialistId: 'pm' }],
    commands: [`/asdt-pm "Add dark mode toggle to user settings"`],
  },
  {
    id: 'document-architecture-decision',
    category: 'understand-document',
    chips: [{ specialistId: 'architect' }],
    commands: [`/asdt-architect "Document the decision to use PostgreSQL row-level security for multi-tenancy"`],
  },
  {
    id: 'write-code-from-settled-spec',
    category: 'add-to-existing',
    chips: [{ specialistId: 'developer' }],
    commands: [`/asdt-developer "Implement the multi-tenancy RLS policy from the Architect's ADR"`],
  },
  {
    id: 'security-audit-existing-feature',
    category: 'review-harden',
    chips: [{ specialistId: 'security' }],
    commands: [`/asdt-security "Review the session management in the auth module"`],
  },
  {
    id: 'design-new-ui-component',
    category: 'from-scratch',
    chips: [{ specialistId: 'ux-ui' }],
    commands: [`/asdt-ux-ui "Design a data table component with sorting, filtering, and pagination"`],
  },
  {
    id: 'validate-test-coverage',
    category: 'review-harden',
    chips: [{ specialistId: 'qa' }],
    commands: [`/asdt-qa`],
  },
  {
    id: 'pickup-developer-existing-adr',
    category: 'add-to-existing',
    chips: [{ specialistId: 'developer' }],
    commands: [`/asdt-developer`],
  },
  {
    id: 'add-security-review-inflight',
    category: 'review-harden',
    chips: [{ specialistId: 'security' }],
    commands: [`/asdt-security`],
  },
  {
    id: 'qa-completed-feature-no-pipeline',
    category: 'review-harden',
    chips: [{ specialistId: 'pm' }, { specialistId: 'qa' }],
    commands: [`/asdt-pm "Add acceptance criteria for the existing checkout flow"`, `/asdt-qa`],
  },
]
