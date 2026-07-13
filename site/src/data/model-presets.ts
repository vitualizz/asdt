export interface ModelPreset {
  id: string
  /** Pre-selected in the installer wizard (Chameleon). Distinct from `recommended`. */
  default?: boolean
  /** Gets the single accent "Recommended" badge on ModelPresetTierCard — no per-tier color scale (Craftsman only). */
  recommended?: boolean
}

// Source of truth: src/content/docs/en/getting-started.md ("Model presets" section)
export const modelPresets: ModelPreset[] = [
  {
    id: 'chameleon',
    default: true,
  },
  {
    id: 'sprinter',
  },
  {
    id: 'craftsman',
    recommended: true,
  },
  {
    id: 'strategist',
  },
  {
    id: 'mastermind',
  },
]
