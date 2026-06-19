import type { SpecialistId } from './artifact-graph'
import type { FlowStep, PipelineFlow } from './pipeline-flow-steps'

// The dogfood flow: the six specialists that actually built THIS deck.
// Rendered by FlowWalker on slide 10. Commands are the real /asdt-* slash
// commands Lee ran; descriptions name the one artifact each step produced
// for this deck. specialistId drives the FlowWalker tint via SPECIALIST_COLOR.
const _ux: SpecialistId = 'ux-ui'

// Local extension of FlowStep: the dogfood flow adds a short terminal-style
// "produces" output line (bilingual) so the FlowWalker right card fills out.
// FlowWalker reads these optionally, so the shared FlowStep type stays untouched.
export interface DogfoodStep extends FlowStep {
  producesEn: string
  producesEs: string
}

export const dogfoodFlow: PipelineFlow & { steps: DogfoodStep[] } = {
  id: 'dogfood',
  titleEn: 'How this deck was built',
  titleEs: 'Cómo se construyó este deck',
  steps: [
    {
      command: '/asdt-researcher',
      specialistId: 'researcher',
      descriptionEn: 'Brainstormed the angles and recommended the Three Laws as the deck’s spine — a discovery brief.',
      descriptionEs: 'Hizo lluvia de ideas de los enfoques y recomendó las Tres Leyes como columna del deck — un resumen de descubrimiento (brief).',
      producesEn: 'researcher/discovery-brief',
      producesEs: 'researcher/discovery-brief',
    },
    {
      command: '/asdt-pm',
      specialistId: 'pm',
      descriptionEn: 'Turned the discovery brief into a 12-slide outline with acceptance criteria — the backlog entry.',
      descriptionEs: 'Convirtió el resumen de descubrimiento en un esquema de 12 slides con criterios de aceptación — la entrada de backlog.',
      producesEn: 'pm/backlog-entry',
      producesEs: 'pm/backlog-entry',
    },
    {
      command: '/asdt-architect',
      specialistId: 'architect',
      descriptionEn: 'Chose the format at the human checkpoint — Astro + MDX over a slide framework — and recorded the ADR.',
      descriptionEs: 'Eligió el formato en el punto de control humano — Astro + MDX en vez de un framework de slides — y registró el ADR.',
      producesEn: 'architect/adr + system-design',
      producesEs: 'architect/adr + system-design',
    },
    {
      command: '/asdt-developer',
      specialistId: 'developer',
      descriptionEn: 'Built the deck-app shell, the DeckController, and the interactive islands — the implementation.',
      descriptionEs: 'Construyó el shell deck-app, el DeckController y las islas interactivas — la implementación.',
      producesEn: 'developer/dev-implementation',
      producesEs: 'developer/dev-implementation',
    },
    {
      command: '/asdt-qa',
      specialistId: 'qa',
      descriptionEn: 'Reviewed keyboard navigation, the no-JS fallback, and bilingual parity — the quality report.',
      descriptionEs: 'Revisó la navegación por teclado, el fallback sin JS y la paridad bilingüe — el reporte de calidad.',
      producesEn: 'qa/qa-review + test-plan',
      producesEs: 'qa/qa-review + test-plan',
    },
    {
      command: '/asdt-ux-ui',
      specialistId: _ux,
      descriptionEn: 'Designed the rail, transitions, and island interactions — the UX brief and component spec.',
      descriptionEs: 'Diseñó el rail, las transiciones y las interacciones de las islas — el resumen de UX (brief) y la spec de componentes.',
      producesEn: 'ux-brief + component-spec',
      producesEs: 'ux-brief + component-spec',
    },
  ],
}
