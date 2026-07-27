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
      trivial:  ['knowledge-recall', 'divergent-ideation', 'decision-preservation'],
      simple:   ['knowledge-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
      moderate: ['knowledge-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
      complex:  ['knowledge-recall', 'divergent-ideation', 'feasibility-scan', 'discovery-brief', 'decision-preservation'],
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
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
      trivial:  ['knowledge-recall', 'feature-intake', 'decision-preservation'],
      simple:   ['knowledge-recall', 'feature-intake', 'user-stories', 'backlog-entry', 'decision-preservation'],
      moderate: ['knowledge-recall', 'feature-intake', 'user-stories', 'success-metrics', 'scope-analysis', 'backlog-entry', 'decision-preservation'],
      complex:  ['knowledge-recall', 'feature-intake', 'user-stories', 'success-metrics', 'scope-analysis', 'prioritization', 'backlog-entry', 'decision-preservation'],
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
      'feature-intake': {
        id: 'feature-intake',
      },
      'user-stories': {
        id: 'user-stories',
      },
      'success-metrics': {
        id: 'success-metrics',
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
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  architect: {
    color: '--c-arch',
    tierType: 'complexity',
    tiers: {
      trivial:  ['knowledge-recall', 'platform-analysis', 'load-constraints', 'decision-preservation'],
      simple:   [],
      moderate: ['knowledge-recall', 'platform-analysis', 'load-constraints', 'evaluate-approaches', 'decision-record', 'technical-handoff', 'decision-preservation'],
      complex:  ['knowledge-recall', 'platform-analysis', 'load-constraints', 'evaluate-approaches', 'decision-record', 'system-design', 'cost-estimation', 'risk-analysis', 'technical-handoff', 'decision-preservation'],
    },
    special: {
      simple: 'not-called',
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
      'platform-analysis': {
        id: 'platform-analysis',
      },
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
      'cost-estimation': {
        id: 'cost-estimation',
      },
      'risk-analysis': {
        id: 'risk-analysis',
      },
      'technical-handoff': {
        id: 'technical-handoff',
      },
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  developer: {
    color: '--c-dev',
    tierType: 'complexity',
    tiers: {
      trivial:  ['knowledge-recall', 'explore', 'decision-preservation'],
      simple:   ['knowledge-recall', 'explore', 'spec', 'implement', 'decision-preservation'],
      moderate: ['knowledge-recall', 'explore', 'spec', 'design', 'implement', 'test', 'decision-preservation'],
      complex:  ['knowledge-recall', 'explore', 'spec', 'design', 'tasks', 'implement', 'test', 'decision-preservation'],
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
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
      'test': {
        id: 'test',
      },
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  qa: {
    color: '--c-qa',
    tierType: 'complexity',
    tiers: {
      trivial:  [],
      simple:   ['knowledge-recall', 'load-requirements', 'ac-validation', 'test-case-generation', 'quality-report', 'performance-validation', 'review', 'decision-preservation'],
      moderate: ['knowledge-recall', 'load-requirements', 'ac-validation', 'edge-case-analysis', 'test-strategy', 'test-case-generation', 'quality-report', 'performance-validation', 'review', 'decision-preservation'],
      complex:  ['knowledge-recall', 'load-requirements', 'ac-validation', 'edge-case-analysis', 'test-strategy', 'test-case-generation', 'quality-report', 'performance-validation', 'review', 'decision-preservation'],
    },
    special: {
      trivial: 'not-eligible',
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
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
      'performance-validation': {
        id: 'performance-validation',
      },
      'review': {
        id: 'review',
      },
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  security: {
    color: '--c-sec',
    tierType: 'risk-surface',
    tiers: {
      none:     [],
      moderate: ['knowledge-recall', 'platform-analysis', 'threat-modeling', 'hardening-checklist', 'decision-preservation'],
      high:     ['knowledge-recall', 'platform-analysis', 'threat-modeling', 'attack-surface', 'owasp-analysis', 'hardening-checklist', 'decision-preservation'],
    },
    special: {
      none: 'not-auto-invoked',
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
      'platform-analysis': {
        id: 'platform-analysis',
      },
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
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },

  'ux-ui': {
    color: '--c-ux',
    tierType: 'complexity',
    tiers: {
      trivial:  ['knowledge-recall', 'platform-analysis', 'feature-brief', 'decision-preservation'],
      simple:   ['knowledge-recall', 'platform-analysis', 'feature-brief', 'design-tokens', 'information-architecture', 'user-flows', 'component-mapping', 'ux-handoff', 'decision-preservation'],
      moderate: ['knowledge-recall', 'platform-analysis', 'feature-brief', 'design-tokens', 'information-architecture', 'user-flows', 'content-design', 'component-mapping', 'ux-handoff', 'decision-preservation'],
      complex:  ['knowledge-recall', 'platform-analysis', 'feature-brief', 'design-tokens', 'information-architecture', 'user-flows', 'content-design', 'component-mapping', 'design-critique', 'ux-handoff', 'decision-preservation'],
    },
    steps: {
      'knowledge-recall': {
        id: 'knowledge-recall',
      },
      'platform-analysis': {
        id: 'platform-analysis',
      },
      'feature-brief': {
        id: 'feature-brief',
      },
      'design-tokens': {
        id: 'design-tokens',
      },
      'information-architecture': {
        id: 'information-architecture',
      },
      'user-flows': {
        id: 'user-flows',
      },
      'content-design': {
        id: 'content-design',
      },
      'component-mapping': {
        id: 'component-mapping',
      },
      'design-critique': {
        id: 'design-critique',
      },
      'ux-handoff': {
        id: 'ux-handoff',
      },
      'decision-preservation': {
        id: 'decision-preservation',
      },
    },
  },
}
