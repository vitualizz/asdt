import type { UIStrings } from '@i18n/types'

export const en: UIStrings = {
  nav: {
    home: 'ASDT',
    docs: 'Docs',
    github: 'GitHub',
    langPickerLabel: 'Language',
    specialists: 'Specialists',
    howItWorks: 'How it works',
    githubStarLabel: 'Star on GitHub',
  },
  a11y: {
    skipToContent: 'Skip to content',
    themeToggleLabel: 'Toggle theme',
    langPickerLabel: 'Select language',
  },
  hero: {
    eyebrow: 'For Claude Code and OpenCode',
    headline: 'A full',
    headlineGrad: 'software team',
    headlineSuffix: ', in your terminal',
    sub: 'AI specialists that pass work to each other over a shared knowledge base. You direct the route; the team forgets nothing.',
    cta: 'Get started',
    secondaryCta: 'Star on GitHub',
    installLabel: 'Install — one command',
    installCmd: 'curl -fsSL https://raw.githubusercontent.com/vitualizz/asdt/main/install.sh | bash',
    copyLabel: 'Copy',
    copiedLabel: 'Copied!',
    copyErrorLabel: 'Copy failed — select manually',
  },
  specialists: {
    kicker: '01 · The team',
    title: 'Meet the specialists',
    sub: 'Focused roles, each with their craft. Invoke any one directly, or ask /asdt to suggest the route.',
    advisorStrip: 'Start here',
    items: [
      { id: 'researcher', name: 'Researcher', desc: 'Explores fuzzy problems before requirements exist: divergent ideation, feasibility scans, one recommended direction.', command: '/asdt-researcher' },
      { id: 'pm', name: 'Product Manager', desc: 'Turns loose ideas into user stories with clear scope and defined criteria.', command: '/asdt-pm' },
      { id: 'architect', name: 'Architect', desc: 'Makes technical decisions: system design, API contracts, decision records.', command: '/asdt-architect' },
      { id: 'developer', name: 'Developer', desc: 'Turns specs and designs into production code, with an implementation plan.', command: '/asdt-developer' },
      { id: 'qa', name: 'QA Engineer', desc: 'Builds the safety net: test plans, acceptance criteria, quality reports.', command: '/asdt-qa' },
      { id: 'security', name: 'Security', desc: 'Finds the gaps an attacker would see first: threat models and hardening.', command: '/asdt-security' },
      { id: 'ux-ui', name: 'UX/UI Design', desc: 'Shapes the experience: flows, components, and accessibility.', command: '/asdt-ux-ui' },
    ],
    orchestrator: {
      id: 'orchestrator',
      name: 'Route Advisor',
      desc: 'Analyzes what you ask and recommends which specialists to involve and in what order. You confirm the plan and run each command.',
      command: '/asdt',
    },
  },
  terminal: {
    tabs: ['New feature', 'PM, step by step', 'Security · STRIDE', 'Your turn_'],
    tryLabel: 'Try it',
  },
  pipeline: {
    title: 'How it works',
    sub: 'Describe what you need — /asdt tells you which specialists to involve and in what order. Each reads the previous one\'s artifacts from the shared knowledge base.',
    nodes: [
      { id: 'researcher', label: 'Researcher' },
      { id: 'pm', label: 'PM' },
      { id: 'architect', label: 'Architect' },
      { id: 'developer', label: 'Developer' },
      { id: 'qa', label: 'QA' },
      { id: 'security', label: 'Security' },
      { id: 'ux-ui', label: 'UX/UI' },
    ],
    a11yTitle: 'ASDT specialist pipeline',
    a11yDesc: 'Diagram showing seven specialists in sequence: Researcher, PM, Architect, Developer, QA, Security, UX/UI — all connected through a shared knowledge base.',
  },
  recipes: {
    kicker: '02 · How it works',
    title: 'Not a pipeline. A team you compose.',
    sub: '/asdt suggests a route for the task — you confirm it, reorder it, or go straight to a specialist.',
    tabs: ['New feature', 'Hotfix', 'Security audit', 'Design pass'],
    notes: [
      '/asdt add workspace permissions per role',
      '/asdt fix crash on empty pagination',
      '/asdt-security review session handling',
      '/asdt redesign onboarding for mobile',
    ],
    kbNote: 'each step saves to the knowledge base · next specialist reads automatically',
  },
  vs: {
    kicker: '03 · Why ASDT',
    title: 'A chat improvises. A team delivers.',
    sub: 'Code assistants are great. With method and memory, they become a team.',
    chatHead: 'Just a chat',
    asdtHead: 'With ASDT',
    items: [
      {
        chat: 'Context degrades as the conversation grows — what you decided yesterday no longer exists.',
        asdt: 'Every decision lives in a shared knowledge base the whole team consults. ASDT remembers.',
      },
      {
        chat: 'Each prompt reinvents the process from scratch — sometimes brilliant, sometimes chaotic.',
        asdt: 'Each specialist works with defined steps — explore, specify, execute. Every time.',
      },
      {
        chat: 'You get to a demo fast, but without stories, decisions, or tests to back it up.',
        asdt: 'Real artifacts remain — stories, decisions, tests — ready to pick up tomorrow.',
      },
    ],
  },
  ctaBand: {
    title: 'Hire the team with one command',
    sub: 'Works with Claude Code and OpenCode. Open source, MIT license.',
  },
  footer: {
    tagline: 'AI-powered. Human-directed.',
    githubLabel: 'GitHub',
    docsLabel: 'Docs',
    licenseLabel: 'MIT License',
    credit: 'Built in the open by Lee Palacios — vitualizz · © 2026 ASDT',
  },
  docs: {
    fallbackNotice: 'This page is only available in English.',
    fallbackNoticeLink: 'View in English',
    gettingStarted: 'Getting Started',
    specialists: 'Specialists',
    commands: 'Commands',
    userFlows: 'User Flows',
    onThisPage: 'On this page',
    stepsTitle: 'Steps it runs',
    searchPlaceholder: 'Search docs...',
    searchNoResults: 'No results found.',
    searchLabel: 'Search documentation',
    complexity: 'Complexity',
    riskSurface: 'Risk surface',
    notCalled: 'Not called at this level — Developer handles it directly.',
    notEligible: 'Not eligible — falls back to simple.',
    notAutoInvoked: 'Not auto-invoked. Still user-invocable on demand.',
    produces: 'Produces',
    pipelinePosition: 'Pipeline position',
    artifactReads: 'Reads from knowledge base',
    artifactWrites: 'Writes to knowledge base',
    flowStep: 'Step',
    flowOf: 'of',
    flowTitle: 'How it runs',
    troubleshooting: 'Troubleshooting',
    recipes: 'Recipes',
    tutorial: 'Tutorial',
    aiAssistants: 'AI Assistants',
    configuration: 'Configuration',
    howItWorksTitle: 'How It Works',
    specialistModel: 'Specialist Model',
    memoryAndEngram: 'Knowledge Base & Memory',
    contributing: 'Contributing',
    development: 'Development',
    groupGettingStarted: 'Getting Started',
    groupConcepts: 'Concepts',
    subgroupOverview: 'Overview',
    subgroupDetail: 'Detail',
    groupReference: 'Reference',
    sidebarAriaLabel: 'Documentation',
    pageNavAriaLabel: 'Page navigation',
    prev: 'Previous',
    next: 'Next',
  },
  commandCard: {
    slashLabel: 'Slash command',
    cliLabel: 'CLI',
  },
  modelPreset: {
    recommendedLabel: 'Recommended',
    defaultLabel: 'Default',
    customizeCta: {
      title: 'Customize per step',
      desc: 'Set the model for each specialist step yourself.',
    },
  },
  specialistComparison: {
    invokeWhen: 'Invoke when…',
    doNotUseWhen: 'Do NOT use when…',
    exampleCommand: 'Example command',
    fullDetailsLink: 'Full details →',
  },
  recipesBrowser: {
    legend: 'Filter by goal',
    empty: 'No recipes match this filter yet — try All.',
    showingPrefix: 'Showing',
    recipeWord: 'recipe',
    recipeWordPlural: 'recipes',
    forWord: 'for',
  },
  data: {
    specialistSteps: {
      'researcher:context-recall': {
        purpose: 'Search organizational memory for prior discovery, related decisions, and known constraints before ideating',
        produces: 'context (inline)',
      },
      'researcher:divergent-ideation': {
        purpose: 'Frame the problem and generate divergent candidate directions — deliberately generative, never selective',
        produces: 'researcher/ideation',
      },
      'researcher:feasibility-scan': {
        purpose: 'Assess each idea with a green/yellow/red feasibility verdict, supporting evidence, and effort estimate',
        produces: 'researcher/feasibility',
      },
      'researcher:discovery-brief': {
        purpose: "Converge to ONE recommended direction with rationale; won't-do candidates seed PM's out-of-scope list",
        produces: 'researcher/discovery-brief',
      },
      'researcher:decision-preservation': {
        purpose: "Preserve the chosen direction as permanent organizational knowledge via the brief's summary field",
        produces: 'summary (inline)',
      },
      'pm:feature-intake': {
        purpose: 'Parse raw request into structured problem statement — extracts problem, goal, stakeholders, flags ambiguities',
        produces: 'pm/feature-intake',
      },
      'pm:user-stories': {
        purpose: 'Write user stories with MoSCoW priorities and 1–3 high-level acceptance criteria per story',
        produces: 'pm/user-stories',
      },
      'pm:success-metrics': {
        purpose: 'Derive measurable NFR targets — budget, target value, and measurement method — from the user stories',
        produces: 'pm/nfr-targets',
      },
      'pm:scope-analysis': {
        purpose: 'Define explicit in/out of scope boundaries, integration points, and scope risk flags',
        produces: 'pm/scope-analysis',
      },
      'pm:prioritization': {
        purpose: 'Order stories by dependency and risk; move Won\'t stories to deferred with explicit reasons',
        produces: 'pm/prioritization',
      },
      'pm:backlog-entry': {
        purpose: 'Consolidate all PM artifacts into the final backlog entry with executive summary and ordered story list',
        produces: 'pm/backlog-entry',
      },
      'architect:load-constraints': {
        purpose: 'Read platform context and classify constraints as HARD (non-negotiable), SOFT (preferences), or OPPORTUNITIES',
        produces: 'architect/constraints-analysis',
      },
      'architect:evaluate-approaches': {
        purpose: 'Compare 2–3 viable approaches; choose one with explicit rationale; document why alternatives were rejected',
        produces: 'architect/approaches',
      },
      'architect:decision-record': {
        purpose: 'Write the ADR: context, decision, alternatives, and consequences — including negatives. Positives-only = incomplete',
        produces: 'architect/adr',
      },
      'architect:system-design': {
        purpose: 'Define data model, API surface (method/inputs/errors), service boundaries, and the happy-path sequence',
        produces: 'architect/system-design',
      },
      'architect:cost-estimation': {
        purpose: 'Estimate the cost profile of the system design per NFR dimension and judge each against the PM budget',
        produces: 'architect/cost-estimate',
      },
      'architect:risk-analysis': {
        purpose: 'Identify top 3–5 risks (performance, security, reliability, coupling, migration) with concrete mitigations',
        produces: 'architect/risks',
      },
      'architect:technical-handoff': {
        purpose: 'Consolidate all architectural work into two final artifacts for Developer and QA',
        produces: 'architect/architectural-decision + architect/system-design-final',
      },
      'developer:explore': {
        purpose: 'Read affected files and modules; map naming patterns and constraints before designing anything',
        produces: 'developer/dev-exploration',
      },
      'developer:spec': {
        purpose: 'Define scope boundary, answer open questions, write acceptance criteria in Given/When/Then format',
        produces: 'developer/dev-spec',
      },
      'developer:design': {
        purpose: 'Choose technical approach, define data model and API shape, list key implementation constraints',
        produces: 'developer/dev-design',
      },
      'developer:tasks': {
        purpose: 'Break implementation into atomic tasks (<2h each) ordered by dependency, with S/M/L estimates',
        produces: 'developer/dev-tasks',
      },
      'developer:implement': {
        purpose: 'Write code for each task — plan-only mode (snippets) or writing mode (real files to declared targets)',
        produces: 'developer/dev-implementation',
      },
      'developer:test': {
        purpose: 'Write tests covering the implementation tasks and edge cases; emits suggested verification commands — never runs them',
        produces: 'developer/dev-tests',
      },
      'qa:load-requirements': {
        purpose: 'Extract and normalize acceptance criteria from upstream artifacts into Given/When/Then format',
        produces: 'qa/ac-list',
      },
      'qa:ac-validation': {
        purpose: 'Review each AC for atomicity, measurability, independence, completeness, unambiguity — rewrite failing ones',
        produces: 'qa/ac-gaps',
      },
      'qa:edge-case-analysis': {
        purpose: 'Discover edge cases via boundary values, equivalence partitioning, state transitions, concurrent access, permission boundaries',
        produces: 'qa/edge-cases',
      },
      'qa:test-strategy': {
        purpose: 'Define test pyramid: what each level covers, what it does not, test data strategy, and flakiness tolerance',
        produces: 'qa/test-strategy',
      },
      'qa:test-case-generation': {
        purpose: 'Write structured test specs (Given/When/Then) for happy path, validated ACs, and critical/high edge cases',
        produces: 'qa/test-cases',
      },
      'qa:quality-report': {
        purpose: 'Verify AC coverage, compute percentage, render READY / READY WITH CAVEATS / BLOCKED verdict',
        produces: 'qa/test-plan',
      },
      'qa:performance-validation': {
        purpose: 'Validate planned or measured performance against the PM NFR targets and emit a go/no-go gate verdict',
        produces: 'qa/perf-validation',
      },
      'qa:review': {
        purpose: 'Go/no-go shipping verdict — holistic release-readiness decision integrating all QA findings before knowledge preservation',
        produces: 'qa/qa-review',
      },
      'security:threat-modeling': {
        purpose: 'Apply STRIDE: Spoofing, Tampering, Repudiation, Information Disclosure, DoS, Elevation of Privilege',
        produces: 'security/stride-threats',
      },
      'security:attack-surface': {
        purpose: 'Map entry points, trust boundaries, data flows — verify validation/sanitization/encoding at each step',
        produces: 'security/attack-surface',
      },
      'security:owasp-analysis': {
        purpose: 'Check all 10 OWASP Top 10 categories (A01–A10) as APPLICABLE/NOT APPLICABLE and MITIGATED/AT RISK',
        produces: 'security/owasp-findings',
      },
      'security:hardening-checklist': {
        purpose: 'Deduplicate findings, prioritize by severity, group by effort: quick wins / medium / significant',
        produces: 'security/security-findings + security/hardening-checklist',
      },
      'ux-ui:feature-brief': {
        purpose: 'Identify primary actor, define core problem (not solution), establish 3–5 observable success criteria',
        produces: 'ux-ui/feature-brief',
      },
      'ux-ui:design-tokens': {
        purpose: 'Derive the design token set from the feature brief design intent and the platform design system fingerprint',
        produces: 'ux-ui/design-tokens',
      },
      'ux-ui:information-architecture': {
        purpose: 'Organize content into sections, prioritize immediate vs. progressive disclosure, define navigation path',
        produces: 'ux-ui/ia',
      },
      'ux-ui:user-flows': {
        purpose: 'Map happy path, error path, and 2–3 edge case flows as numbered steps from the actor\'s perspective',
        produces: 'ux-ui/flows',
      },
      'ux-ui:content-design': {
        purpose: 'Catalog the text touchpoints surfaced by the flows, elaborating the IA content intent into an illustrative content inventory',
        produces: 'ux-ui/content-inventory',
      },
      'ux-ui:component-mapping': {
        purpose: 'Classify each UI state as reuse / extend / new — quality gate: >2:1 reuse ratio required',
        produces: 'ux-ui/components',
      },
      'ux-ui:design-critique': {
        purpose: 'Single-pass annotation of the component inventory against the design tokens for consistency and accessibility; derives needs_review',
        produces: 'ux-ui/design-critique',
      },
      'ux-ui:ux-handoff': {
        purpose: 'Consolidate all UX work into ux-brief (flows + IA) and component-spec (inventory + props + events)',
        produces: 'ux-ui/ux-brief + ux-ui/component-spec',
      },
    },
    pipelineFlows: {
      'full-feature': {
        title: 'Getting a pipeline suggestion',
        steps: {
          orchestrate: {
            description: 'Ask the orchestrator — ASDT analyzes the request and recommends which specialists to involve and in what order.',
          },
          pm: {
            description: 'PM defines scope, writes user stories with acceptance criteria, saves pm/backlog-entry to the knowledge base.',
          },
          architect: {
            description: 'Architect reads the backlog entry, designs the token flow and API contracts, saves architectural-decision + system-design-final.',
          },
          developer: {
            description: 'Developer reads the ADR and system design, implements the magic link handler, saves dev-implementation.',
          },
          security: {
            description: 'Security reviews the auth mechanism, runs STRIDE and OWASP analysis, saves security-findings + hardening-checklist.',
          },
        },
      },
      'mid-pipeline': {
        title: 'Picking up mid-pipeline',
        steps: {
          developer: {
            description: 'Developer reads prior artifacts from the knowledge base automatically — even from a previous session. No manual context passing.',
          },
          qa: {
            description: 'QA loads dev-implementation and runs its full workflow: AC validation, edge-case analysis, and test case generation.',
          },
        },
      },
    },
    recipes: {
      'ship-new-feature': {
        title: 'Ship a new user-facing feature',
        note: 'Start here when a feature request is in vague language and needs the full PM → Architect → Developer sequence.',
        kbNote: 'Review the suggested specialist sequence, then run each in order.',
      },
      'new-rest-endpoint': {
        title: 'Build a new REST API endpoint',
        note: 'For backend-only changes where the API contract is the primary design decision.',
      },
      'new-screen-ux-first': {
        title: 'New screen with UX design first',
        note: 'For user-visible features where UI design should precede implementation.',
      },
      'security-review-before-shipping': {
        title: 'Feature with security review before shipping',
        note: 'For anything touching auth, payments, PII (personally identifiable information), or external integrations.',
      },
      'explore-before-planning': {
        title: 'Explore before planning (fuzzy problem)',
        note: 'When the problem is unclear and you need discovery before requirements.',
        kbNote: 'Researcher produces a discovery brief with a recommended direction — hand it to PM next.',
      },
      'lock-scope-user-stories': {
        title: 'Lock scope and write user stories',
        note: 'When you have a clear feature idea and just need structured requirements.',
      },
      'document-architecture-decision': {
        title: 'Document an architecture decision',
        note: 'When the technical approach needs a formal ADR and system design.',
      },
      'write-code-from-settled-spec': {
        title: 'Write production code from a settled spec',
        note: 'When scope and architecture are already locked in the knowledge base.',
      },
      'security-audit-existing-feature': {
        title: 'Security audit on an existing feature',
        note: 'Run at any point — no prior specialist run required.',
      },
      'design-new-ui-component': {
        title: 'Design a new UI component',
        note: 'When you need a component spec before the developer starts coding.',
      },
      'validate-test-coverage': {
        title: 'Validate test coverage before shipping',
        note: 'Run after Developer — QA reads the implementation artifact automatically.',
      },
      'pickup-developer-existing-adr': {
        title: 'Pick up at Developer after an existing ADR',
        note: 'When PM and Architect ran in a previous session. Developer loads prior artifacts automatically.',
        kbNote: 'No need to re-run PM or Architect — artifacts are in the knowledge base.',
      },
      'add-security-review-inflight': {
        title: 'Add a security review to an in-flight pipeline',
        note: 'Run Security at any point without restarting the pipeline.',
        kbNote: 'Security reads system-design-final and dev-implementation from the knowledge base.',
      },
      'qa-completed-feature-no-pipeline': {
        title: 'QA a completed feature without a full pipeline',
        note: 'When a feature was built without ASDT — run QA against the existing code.',
      },
    },
    recipeCategories: {
      all: { label: 'All' },
      'from-scratch': { label: 'Start from scratch' },
      'add-to-existing': { label: 'Add to existing code' },
      'review-harden': { label: 'Review & harden' },
      'understand-document': { label: 'Understand & document' },
    },
    commands: {
      'asdt-tui': {
        title: 'Interactive terminal UI',
        oneLiner: 'The only CLI tool — checks your setup and installs or updates the ASDT specialist skills.',
      },
      'asdt-init': {
        title: 'Initialize ASDT',
        oneLiner: 'Detects your stack and creates .asdt/config.yaml and .asdt/knowledge/platform.yaml.',
      },
      asdt: {
        title: 'Pipeline routing suggestion',
        oneLiner: 'Analyzes the request and recommends which specialists to involve and in what order.',
      },
      'asdt-researcher': {
        title: 'Researcher',
        oneLiner: 'Runs the Researcher specialist only — discovery for fuzzy problems before requirements exist.',
      },
      'asdt-pm': {
        title: 'Product Manager',
        oneLiner: 'Runs the Product Manager specialist only.',
      },
      'asdt-architect': {
        title: 'Architect',
        oneLiner: 'Runs the Architect specialist only.',
      },
      'asdt-developer': {
        title: 'Developer',
        oneLiner: 'Runs the Developer specialist only.',
      },
      'asdt-qa': {
        title: 'QA',
        oneLiner: 'Runs the QA specialist only.',
      },
      'asdt-security': {
        title: 'Security',
        oneLiner: 'Runs the Security specialist only.',
      },
      'asdt-ux-ui': {
        title: 'UX/UI',
        oneLiner: 'Runs the UX/UI specialist only.',
      },
    },
    specialistComparison: {
      pm: {
        teaser: 'Locks scope and turns vague requests into structured user stories.',
        invokeWhen: "The request is vague or user-facing, scope isn't locked, or user stories don't exist yet",
        produces: 'pm/backlog-entry — feature summary, ordered user stories with AC, in/out scope, risk flags',
        doNotUseWhen: 'You already have a clear backlog entry — re-running PM regenerates stories from scratch',
      },
      architect: {
        teaser: 'Decides how the pieces fit together before anyone writes code.',
        invokeWhen: 'The solution touches service boundaries, data models, or API contracts, or has two viable technical approaches worth documenting',
        produces: 'architectural-decision (ADR) + system-design-final — data model, API surface, service boundaries',
        doNotUseWhen: 'You need implementation code, test plans, or UX specs — Architect produces decisions, not code',
      },
      developer: {
        teaser: 'Turns a settled design into an ordered implementation plan or real code.',
        invokeWhen: 'The shape of the solution is settled and you need an ordered implementation plan or production code written to the repo',
        produces: 'developer/dev-implementation — ordered file manifest and code plan',
        doNotUseWhen: "You haven't locked scope or architecture yet — Developer will implement against ambiguous requirements",
      },
      qa: {
        teaser: 'Turns acceptance criteria into a systematic test plan with a go/no-go verdict.',
        invokeWhen: "Code is ready for review, AC exists but hasn't been validated, or you need systematic edge case and boundary coverage",
        produces: 'test-plan — AC coverage %, uncovered gaps, full Given/When/Then test case list, quality verdict',
        doNotUseWhen: 'You want executable test code written — QA produces test specifications, not runnable code',
      },
      security: {
        teaser: 'Hunts for auth, data, and integration risks before they ship.',
        invokeWhen: 'The feature touches auth, sessions, PII, external integrations, webhooks, or new public API endpoints',
        produces: 'security-findings (severity-rated, CWE-referenced) + hardening-checklist — must-fix vs can-defer',
        doNotUseWhen: 'You want implementation code or architectural decisions — Security produces findings and checklists only',
      },
      'ux-ui': {
        teaser: 'Maps flows and components before implementation starts.',
        invokeWhen: 'A new screen or feature-level UI needs design before implementation begins, or user flows need mapping',
        produces: 'ux-brief (flows, IA, success criteria) + component-spec — inventory of reused/extended/new components',
        doNotUseWhen: 'The screen has already been built — a UX spec delivered after implementation is too late to shape it',
      },
      researcher: {
        teaser: 'Explores a fuzzy problem before anyone commits to a direction.',
        invokeWhen: 'The problem is fuzzy or open-ended — you need discovery and framing before requirements can be written',
        produces: 'researcher/discovery-brief — problem framing, divergent idea set, feasibility scan, single recommended direction',
        doNotUseWhen: 'You already have a defined problem statement — Researcher explores; it does not produce user stories or ADRs',
      },
    },
    tutorialStages: {
      install: { label: 'Install' },
      initialize: { label: 'Initialize' },
      recommendation: { label: 'Recommendation' },
      pm: { label: 'PM' },
      architect: { label: 'Architect' },
      developer: { label: 'Developer' },
    },
    artifactAnatomy: {
      title: { desc: 'The human-readable name of the artifact.' },
      topicKey: { desc: 'The machine key used to retrieve it automatically, no fuzzy matching involved.' },
      type: { desc: 'The artifact\'s category — architecture, decision, bugfix, etc.' },
      project: { desc: 'The project it belongs to, so results never mix across projects.' },
    },
    modelPresets: {
      chameleon: {
        label: 'Chameleon',
        desc: 'Keeps the model your assistant already has defined (strips the model: field so each assistant uses its own default).',
      },
      sprinter: {
        label: 'Sprinter',
        desc: 'Fastest and cheapest across the board.',
      },
      craftsman: {
        label: 'Craftsman',
        desc: 'The recommended balance of speed and capability (the shipped defaults, verbatim).',
      },
      strategist: {
        label: 'Strategist',
        desc: 'More capability for analysis and decisions.',
      },
      mastermind: {
        label: 'Mastermind',
        desc: 'Maximum capability where it matters most.',
      },
    },
    artifactSentinels: {
      'Problem (raw)': { label: 'Problem (raw)' },
      'Request (raw)': { label: 'Request (raw)' },
    },
  },
}
