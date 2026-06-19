import type { SpecialistId } from './artifact-graph'
import type { PipelineFlow } from './pipeline-flow-steps'

// The dogfood flow: the six specialists that actually built THIS deck.
// Rendered by FlowWalker on slide 10. Commands are the real /asdt-* slash
// commands Lee ran; descriptions name the one artifact each step produced
// for this deck. specialistId drives the FlowWalker tint via SPECIALIST_COLOR.
const _ux: SpecialistId = 'ux-ui'

export const dogfoodFlow: PipelineFlow = {
  id: 'dogfood',
  titleEn: 'How this deck was built',
  titleEs: 'Cómo se construyó este deck',
  steps: [
    {
      command: '/asdt-researcher',
      specialistId: 'researcher',
      descriptionEn: 'Brainstormed the angles and recommended the Three Laws as the deck’s spine — a discovery brief.',
      descriptionEs: 'Hizo lluvia de ideas de los enfoques y recomendó las Tres Leyes como columna del deck — un discovery brief.',
    },
    {
      command: '/asdt-pm',
      specialistId: 'pm',
      descriptionEn: 'Turned the discovery brief into a 12-slide outline with acceptance criteria — the backlog entry.',
      descriptionEs: 'Convirtió el discovery brief en un esquema de 12 slides con criterios de aceptación — la entrada de backlog.',
    },
    {
      command: '/asdt-architect',
      specialistId: 'architect',
      descriptionEn: 'Chose the format at the human gate — Astro + MDX over a slide framework — and recorded the ADR.',
      descriptionEs: 'Eligió el formato en el gate humano — Astro + MDX en vez de un framework de slides — y registró el ADR.',
    },
    {
      command: '/asdt-developer',
      specialistId: 'developer',
      descriptionEn: 'Built the deck-app shell, the DeckController, and the interactive islands — the implementation.',
      descriptionEs: 'Construyó el shell deck-app, el DeckController y las islas interactivas — la implementación.',
    },
    {
      command: '/asdt-qa',
      specialistId: 'qa',
      descriptionEn: 'Reviewed keyboard navigation, the no-JS fallback, and bilingual parity — the quality report.',
      descriptionEs: 'Revisó la navegación por teclado, el fallback sin JS y la paridad bilingüe — el reporte de calidad.',
    },
    {
      command: '/asdt-ux-ui',
      specialistId: _ux,
      descriptionEn: 'Designed the rail, transitions, and island interactions — the UX brief and component spec.',
      descriptionEs: 'Diseñó el rail, las transiciones y las interacciones de las islas — el UX brief y la spec de componentes.',
    },
  ],
}
