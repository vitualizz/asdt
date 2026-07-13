export type ComplexityTier = 'trivial' | 'simple' | 'moderate' | 'complex'
export type RiskTier = 'none' | 'moderate' | 'high'

export interface Step {
  id: string
}

export interface SpecialistConfig {
  color: string
  tierType: 'complexity' | 'risk-surface'
  tiers: Partial<Record<ComplexityTier | RiskTier, string[]>>
  special?: Partial<Record<ComplexityTier | RiskTier, 'not-called' | 'not-eligible' | 'not-auto-invoked'>>
  steps: Record<string, Step>
}

// Source of truth: skill/asdt-{id}/SKILL.md and steps/*.md files
export const specialistSteps: Record<string, SpecialistConfig> = {
  researcher: {
    color: '--c-res',
    tierType: 'complexity',
    tiers: {
      trivial:  ['context-recall', 'divergent-ideation', 'decision-preservation'],
      simple:   ['context-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
      moderate: ['context-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
      complex:  ['context-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
    },
    steps: {
      'context-recall': {
        id: 'context-recall',
      },
      'divergent-ideation': {
        id: 'divergent-ideation',
      },
      'feasibility-scan': {
        id: 'feasibility-scan',
      },
      'discovery-brief': {
        id: 'discovery-brief',
      },
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  pm: {
    color: '--c-pm',
    tierType: 'complexity',
    tiers: {
      trivial:  ['feature-intake'],
      simple:   ['feature-intake', 'user-stories', 'backlog-entry'],
      moderate: ['feature-intake', 'user-stories', 'scope-analysis', 'backlog-entry'],
      complex:  ['feature-intake', 'user-stories', 'scope-analysis', 'prioritization', 'backlog-entry'],
    },
    steps: {
      'feature-intake': {
        id: 'feature-intake',
      },
      'user-stories': {
        id: 'user-stories',
      },
      'scope-analysis': {
        id: 'scope-analysis',
      },
      'prioritization': {
        id: 'prioritization',
      },
      'backlog-entry': {
        id: 'backlog-entry',
      },
    },
  },

  architect: {
    color: '--c-arch',
    tierType: 'complexity',
    tiers: {
      trivial:  ['load-constraints'],
      simple:   [],
      moderate: ['load-constraints', 'evaluate-approaches', 'decision-record'],
      complex:  ['load-constraints', 'evaluate-approaches', 'decision-record', 'system-design', 'risk-analysis', 'technical-handoff'],
    },
    special: {
      simple: 'not-called',
    },
    steps: {
      'load-constraints': {
        id: 'load-constraints',
      },
      'evaluate-approaches': {
        id: 'evaluate-approaches',
      },
      'decision-record': {
        id: 'decision-record',
      },
      'system-design': {
        id: 'system-design',
      },
      'risk-analysis': {
        id: 'risk-analysis',
      },
      'technical-handoff': {
        id: 'technical-handoff',
      },
    },
  },

  developer: {
    color: '--c-dev',
    tierType: 'complexity',
    tiers: {
      trivial:  ['explore'],
      simple:   ['explore', 'spec', 'implement'],
      moderate: ['explore', 'spec', 'design', 'implement'],
      complex:  ['explore', 'spec', 'design', 'tasks', 'implement'],
    },
    steps: {
      'explore': {
        id: 'explore',
      },
      'spec': {
        id: 'spec',
      },
      'design': {
        id: 'design',
      },
      'tasks': {
        id: 'tasks',
      },
      'implement': {
        id: 'implement',
      },
    },
  },

  qa: {
    color: '--c-qa',
    tierType: 'complexity',
    tiers: {
      trivial:  [],
      simple:   ['load-requirements', 'ac-validation', 'test-case-generation', 'quality-report'],
      moderate: ['load-requirements', 'ac-validation', 'edge-case-analysis', 'test-strategy', 'test-case-generation', 'quality-report', 'review'],
      complex:  ['load-requirements', 'ac-validation', 'edge-case-analysis', 'test-strategy', 'test-case-generation', 'quality-report', 'review'],
    },
    special: {
      trivial: 'not-eligible',
    },
    steps: {
      'load-requirements': {
        id: 'load-requirements',
      },
      'ac-validation': {
        id: 'ac-validation',
      },
      'edge-case-analysis': {
        id: 'edge-case-analysis',
      },
      'test-strategy': {
        id: 'test-strategy',
      },
      'test-case-generation': {
        id: 'test-case-generation',
      },
      'quality-report': {
        id: 'quality-report',
      },
      'review': {
        id: 'review',
      },
    },
  },

  security: {
    color: '--c-sec',
    tierType: 'risk-surface',
    tiers: {
      none:     [],
      moderate: ['threat-modeling', 'hardening-checklist'],
      high:     ['threat-modeling', 'attack-surface', 'owasp-analysis', 'hardening-checklist'],
    },
    special: {
      none: 'not-auto-invoked',
    },
    steps: {
      'threat-modeling': {
        id: 'threat-modeling',
      },
      'attack-surface': {
        id: 'attack-surface',
      },
      'owasp-analysis': {
        id: 'owasp-analysis',
      },
      'hardening-checklist': {
        id: 'hardening-checklist',
      },
    },
  },

  'ux-ui': {
    color: '--c-ux',
    tierType: 'complexity',
    tiers: {
      trivial:  ['feature-brief'],
      simple:   ['feature-brief', 'information-architecture', 'user-flows', 'component-mapping', 'ux-handoff'],
      moderate: ['feature-brief', 'information-architecture', 'user-flows', 'component-mapping', 'ux-handoff'],
      complex:  ['feature-brief', 'information-architecture', 'user-flows', 'component-mapping', 'responsive-strategy', 'ux-handoff'],
    },
    steps: {
      'feature-brief': {
        id: 'feature-brief',
      },
      'information-architecture': {
        id: 'information-architecture',
      },
      'user-flows': {
        id: 'user-flows',
      },
      'component-mapping': {
        id: 'component-mapping',
      },
      'responsive-strategy': {
        id: 'responsive-strategy',
      },
      'ux-handoff': {
        id: 'ux-handoff',
      },
    },
  },
}
