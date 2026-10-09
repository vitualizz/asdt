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
   *  and it picks from what the request asks for — not from a level anyone passes in.
   *  The last chain must be the fullest one: components render its steps as the full list. */
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
    // `build` stays last: StepFlow and AllSpecialistsOverview read the last chain
    // as the full step list.
    chains: [
      { when: 'question', steps: ['knowledge-recall', 'platform-analysis', 'explore'] },
      { when: 'plan', steps: ['knowledge-recall', 'platform-analysis', 'explore', 'spec'] },
      // A resume of a plan built but never verified runs `verify` alone.
      { when: 'resume', steps: ['knowledge-recall', 'platform-analysis', 'approve', 'implement', 'verify'] },
      { when: 'build', steps: ['knowledge-recall', 'platform-analysis', 'explore', 'spec', 'approve', 'implement', 'verify'] },
    ],
    steps: {
      'knowledge-recall': { id: 'knowledge-recall', execution: 'inline' },
      'platform-analysis': { id: 'platform-analysis', execution: 'inline' },
      explore: { id: 'explore', execution: 'subagent' },
      spec: { id: 'spec', execution: 'subagent' },
      approve: { id: 'approve', execution: 'inline' },
      implement: { id: 'implement', execution: 'subagent' },
      verify: { id: 'verify', execution: 'inline' },
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
