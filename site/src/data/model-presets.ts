export interface ModelPreset {
  id: string
  labelEn: string
  labelEs: string
  descEn: string
  descEs: string
  /** Pre-selected in the installer wizard (Chameleon). Distinct from `recommended`. */
  default?: boolean
  /** Gets the single accent "Recommended" badge on ModelPresetTierCard — no per-tier color scale (Craftsman only). */
  recommended?: boolean
}

// Source of truth: src/content/docs/en/getting-started.md ("Model presets" section)
export const modelPresets: ModelPreset[] = [
  {
    id: 'chameleon',
    labelEn: 'Chameleon',
    labelEs: 'Chameleon',
    descEn: 'Keeps the model your assistant already has defined (strips the model: field so each assistant uses its own default).',
    descEs: 'Mantiene el modelo que tu asistente ya tiene definido (quita el campo model: para que cada asistente use su propio valor por defecto).',
    default: true,
  },
  {
    id: 'sprinter',
    labelEn: 'Sprinter',
    labelEs: 'Sprinter',
    descEn: 'Fastest and cheapest across the board.',
    descEs: 'El más rápido y económico en todo.',
  },
  {
    id: 'craftsman',
    labelEn: 'Craftsman',
    labelEs: 'Craftsman',
    descEn: 'The recommended balance of speed and capability (the shipped defaults, verbatim).',
    descEs: 'El balance recomendado entre velocidad y capacidad (los valores por defecto que se envían, tal cual).',
    recommended: true,
  },
  {
    id: 'strategist',
    labelEn: 'Strategist',
    labelEs: 'Strategist',
    descEn: 'More capability for analysis and decisions.',
    descEs: 'Más capacidad para análisis y decisiones.',
  },
  {
    id: 'mastermind',
    labelEn: 'Mastermind',
    labelEs: 'Mastermind',
    descEn: 'Maximum capability where it matters most.',
    descEs: 'Máxima capacidad donde más importa.',
  },
]

// Separate CTA below the tier-card grid, NOT a 6th card.
export const CUSTOMIZE_CTA = {
  titleEn: 'Customize per step',
  titleEs: 'Personalizar por paso',
  descEn: 'Set the model for each specialist step yourself.',
  descEs: 'Configurá el modelo de cada paso de especialista vos mismo.',
}
