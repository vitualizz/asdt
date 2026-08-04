export interface SpecialistItem {
  id: string
  name: string
  desc: string
  command: string
}

export interface PipelineNode {
  id: string
  label: string
}

export interface VsItem {
  chat: string
  asdt: string
}

export interface UIStrings {
  nav: {
    home: string
    docs: string
    github: string
    langPickerLabel: string
    specialists: string
    howItWorks: string
    githubStarLabel: string
  }
  a11y: {
    skipToContent: string
    themeToggleLabel: string
    langPickerLabel: string
  }
  hero: {
    eyebrow: string
    headline: string
    headlineGrad: string
    headlineSuffix: string
    sub: string
    cta: string
    secondaryCta: string
    installLabel: string
    installCmd: string
    copyLabel: string
    copiedLabel: string
    copyErrorLabel: string
  }
  specialists: {
    kicker: string
    title: string
    sub: string
    advisorStrip: string
    items: SpecialistItem[]
    orchestrator: SpecialistItem
  }
  terminal: {
    tabs: [string, string, string, string]
    tryLabel: string
  }
  pipeline: {
    title: string
    sub: string
    nodes: PipelineNode[]
    a11yTitle: string
    a11yDesc: string
  }
  recipes: {
    kicker: string
    title: string
    sub: string
    tabs: [string, string, string, string]
    notes: [string, string, string, string]
    kbNote: string
  }
  vs: {
    kicker: string
    title: string
    sub: string
    chatHead: string
    asdtHead: string
    items: VsItem[]
  }
  ctaBand: {
    title: string
    sub: string
  }
  footer: {
    tagline: string
    githubLabel: string
    docsLabel: string
    licenseLabel: string
    credit: string
  }
  docs: {
    fallbackNotice: string
    fallbackNoticeLink: string
    gettingStarted: string
    specialists: string
    commands: string
    userFlows: string
    onThisPage: string
    stepsTitle: string
    searchPlaceholder: string
    searchNoResults: string
    searchLabel: string
    complexity: string
    riskSurface: string
    notCalled: string
    notEligible: string
    notAutoInvoked: string
    produces: string
    pipelinePosition: string
    artifactReads: string
    artifactWrites: string
    flowStep: string
    flowOf: string
    flowTitle: string
    troubleshooting: string
    recipes: string
    tutorial: string
    aiAssistants: string
    configuration: string
    howItWorksTitle: string
    specialistModel: string
    memoryAndEngram: string
    contributing: string
    development: string
    groupGettingStarted: string
    groupConcepts: string
    subgroupOverview: string
    subgroupDetail: string
    groupReference: string
    sidebarAriaLabel: string
    pageNavAriaLabel: string
    prev: string
    next: string
  }
  commandCard: {
    slashLabel: string
    cliLabel: string
  }
  modelPreset: {
    recommendedLabel: string
    defaultLabel: string
    customizeCta: {
      title: string
      desc: string
    }
  }
  specialistComparison: {
    invokeWhen: string
    doNotUseWhen: string
    exampleCommand: string
    fullDetailsLink: string
  }
  recipesBrowser: {
    legend: string
    empty: string
    showingPrefix: string
    recipeWord: string
    recipeWordPlural: string
    forWord: string
  }
  data: {
    chainsLabel: string
    chains: Record<string, string>
    specialistSteps: Record<string, { purpose: string; produces: string }>
    pipelineFlows: Record<string, {
      title: string
      steps: Record<string, { description: string }>
    }>
    recipes: Record<string, { title: string; note: string; kbNote?: string }>
    recipeCategories: Record<string, { label: string }>
    commands: Record<string, { title: string; oneLiner: string }>
    specialistComparison: Record<string, {
      teaser: string
      invokeWhen: string
      produces: string
      doNotUseWhen: string
    }>
    tutorialStages: Record<string, { label: string }>
    artifactAnatomy: Record<string, { desc: string }>
    modelPresets: Record<string, { label: string; desc: string }>
    artifactSentinels: Record<string, { label: string }>
  }
}
