export interface Step {
  id: string
  execution: 'inline' | 'subagent'
}

export interface Chain {
  /** i18n key suffix under data.chains — what the request has to look like for this chain to run */
  when?: string
  steps: string[]
}

export interface SpecialistConfig {
  color: string
  /** Every specialist runs one fixed chain. Only the Developer picks between chains,
   *  and it picks from what the request asks for — not from a level anyone passes in. */
  chains: Chain[]
  steps: Record<string, Step>
}

// Source of truth: skill/asdt-{id}/workflow.yaml
export const specialistSteps: Record<string, SpecialistConfig> = {
  researcher: {
    color: '--c-res',
    chains: [{ steps: ['knowledge-recall', 'discovery'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      discovery: { id: 'discovery', execution: 'subagent' },
    },
  },

  pm: {
    color: '--c-pm',
    chains: [{ steps: ['knowledge-recall', 'backlog'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      backlog: { id: 'backlog', execution: 'subagent' },
    },
  },

  'ux-ui': {
    color: '--c-ux',
    chains: [{ steps: ['knowledge-recall', 'platform-analysis', 'ux-spec'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      'platform-analysis': { id: 'platform-analysis', execution: 'inline' },
      'ux-spec': { id: 'ux-spec', execution: 'subagent' },
    },
  },

  architect: {
    color: '--c-arch',
    chains: [{ steps: ['knowledge-recall', 'platform-analysis', 'design'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      'platform-analysis': { id: 'platform-analysis', execution: 'inline' },
      design: { id: 'design', execution: 'subagent' },
    },
  },

  developer: {
    color: '--c-dev',
    chains: [
      { when: 'question', steps: ['knowledge-recall', 'explore'] },
      { when: 'plan', steps: ['knowledge-recall', 'explore', 'spec'] },
      { when: 'build', steps: ['knowledge-recall', 'explore', 'spec', 'implement'] },
    ],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      explore: { id: 'explore', execution: 'subagent' },
      spec: { id: 'spec', execution: 'subagent' },
      implement: { id: 'implement', execution: 'subagent' },
    },
  },

  security: {
    color: '--c-sec',
    chains: [{ steps: ['knowledge-recall', 'platform-analysis', 'assess', 'harden'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      'platform-analysis': { id: 'platform-analysis', execution: 'inline' },
      assess: { id: 'assess', execution: 'subagent' },
      harden: { id: 'harden', execution: 'subagent' },
    },
  },

  qa: {
    color: '--c-qa',
    chains: [{ steps: ['knowledge-recall', 'test-plan'] }],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      'test-plan': { id: 'test-plan', execution: 'subagent' },
    },
  },
}
