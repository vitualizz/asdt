import { SPECIALIST_COLOR, type SpecialistId } from './artifact-graph'

export interface TutorialStage {
  id: string
  /** Heading id in src/content/docs/en/tutorial.mdx, verified against the built page. */
  anchorIdEn: string
  /** Heading id in src/content/docs/es/tutorial.mdx (Spanish slugs differ), verified against the built page. */
  anchorIdEs: string
  specialistId?: SpecialistId
  /** Reuses SPECIALIST_COLOR for the last 3 (specialist) stages; neutral --color-accent for the first 3. */
  color: string
}

const NEUTRAL_STAGE_COLOR = '--color-accent'

// Source of truth: src/content/docs/en/tutorial.md — the tutorial's OWN 6 step headings
// (NOT the full 7-specialist PIPELINE_ORDER — the tutorial never exercises QA, Security,
// UX/UI, or Researcher). Only the last 3 stages map to a SpecialistId/color.
export const tutorialStages: TutorialStage[] = [
  {
    id: 'install',
    anchorIdEn: 'step-1--install-asdt',
    anchorIdEs: 'paso-1--instalar-asdt',
    color: NEUTRAL_STAGE_COLOR,
  },
  {
    id: 'initialize',
    anchorIdEn: 'step-2--initialize-asdt-in-your-project',
    anchorIdEs: 'paso-2--inicializar-asdt-en-tu-proyecto',
    color: NEUTRAL_STAGE_COLOR,
  },
  {
    id: 'recommendation',
    anchorIdEn: 'step-3--ask-asdt-for-a-pipeline-recommendation',
    anchorIdEs: 'paso-3--pedirle-a-asdt-una-recomendación-de-pipeline',
    color: NEUTRAL_STAGE_COLOR,
  },
  {
    id: 'pm',
    anchorIdEn: 'step-4--run-the-product-manager',
    anchorIdEs: 'paso-4--correr-el-product-manager',
    specialistId: 'pm',
    color: SPECIALIST_COLOR.pm,
  },
  {
    id: 'architect',
    anchorIdEn: 'step-5--run-the-architect',
    anchorIdEs: 'paso-5--correr-el-arquitecto',
    specialistId: 'architect',
    color: SPECIALIST_COLOR.architect,
  },
  {
    id: 'developer',
    anchorIdEn: 'step-6--run-the-developer',
    anchorIdEs: 'paso-6--correr-el-developer',
    specialistId: 'developer',
    color: SPECIALIST_COLOR.developer,
  },
]
