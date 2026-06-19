import { useState } from 'preact/hooks'

interface PipelineStep {
  slug: string
  purposeEn: string
  purposeEs: string
  produces: string
}

interface Specialist {
  id: string
  colorVar: string
  nameEn: string
  nameEs: string
  roleEn: string
  roleEs: string
  steps: PipelineStep[]
  artifactEn: string
  artifactEs: string
}

const PIPELINE_SPECIALISTS: Specialist[] = [
  {
    id: 'researcher',
    colorVar: 'var(--dk-c-researcher)',
    nameEn: 'Researcher',
    nameEs: 'Investigador',
    roleEn: 'Discovery — diverges before PM converges',
    roleEs: 'Descubrimiento — diverge antes de que PM converja',
    artifactEn: 'Discovery brief',
    artifactEs: 'Brief de descubrimiento',
    steps: [
      {
        slug: 'context-recall',
        purposeEn: 'Search organizational memory for prior discovery, related decisions, and known constraints before ideating',
        purposeEs: 'Busca en la memoria organizacional descubrimientos previos, decisiones relacionadas y restricciones conocidas antes de idear',
        produces: 'context (inline)',
      },
      {
        slug: 'divergent-ideation',
        purposeEn: 'Frame the problem and generate divergent candidate directions — deliberately generative, never selective',
        purposeEs: 'Enmarca el problema y genera direcciones candidatas divergentes — deliberadamente generativo, nunca selectivo',
        produces: 'researcher/ideation',
      },
      {
        slug: 'feasibility-scan',
        purposeEn: 'Assess each idea with a green/yellow/red feasibility verdict, supporting evidence, and effort estimate',
        purposeEs: 'Evalúa cada idea con un veredicto de factibilidad verde/amarillo/rojo, evidencia de soporte y estimación de esfuerzo',
        produces: 'researcher/feasibility',
      },
      {
        slug: 'discovery-brief',
        purposeEn: "Converge to ONE recommended direction with rationale; won't-do candidates seed PM's out-of-scope list",
        purposeEs: 'Converge a UNA dirección recomendada con justificación; los candidatos descartados alimentan la lista fuera de alcance del PM',
        produces: 'researcher/discovery-brief',
      },
      {
        slug: 'decision-preservation',
        purposeEn: "Preserve the chosen direction as permanent organizational knowledge via the brief's summary field",
        purposeEs: 'Preserva la dirección elegida como conocimiento organizacional permanente mediante el campo de resumen del brief',
        produces: 'summary (inline)',
      },
    ],
  },
  {
    id: 'pm',
    colorVar: 'var(--dk-c-pm)',
    nameEn: 'PM',
    nameEs: 'PM',
    roleEn: 'Formalizes requirements',
    roleEs: 'Formaliza los requisitos',
    artifactEn: 'Backlog entry',
    artifactEs: 'Entrada de backlog',
    steps: [
      {
        slug: 'feature-intake',
        purposeEn: 'Parse raw request into structured problem statement — extracts problem, goal, stakeholders, flags ambiguities',
        purposeEs: 'Convierte la solicitud cruda en un enunciado de problema estructurado — extrae problema, objetivo, interesados y marca ambigüedades',
        produces: 'pm/feature-intake',
      },
      {
        slug: 'user-stories',
        purposeEn: 'Write user stories with MoSCoW priorities and 1–3 high-level acceptance criteria per story',
        purposeEs: 'Escribe historias de usuario con prioridades MoSCoW y 1–3 criterios de aceptación de alto nivel por historia',
        produces: 'pm/user-stories',
      },
      {
        slug: 'scope-analysis',
        purposeEn: 'Define explicit in/out of scope boundaries, integration points, and scope risk flags',
        purposeEs: 'Define los límites explícitos de dentro/fuera de alcance, los puntos de integración y las alertas de riesgo de alcance',
        produces: 'pm/scope-analysis',
      },
      {
        slug: 'backlog-entry',
        purposeEn: 'Consolidate all PM artifacts into the final backlog entry with executive summary and ordered story list',
        purposeEs: 'Consolida todos los artefactos del PM en la entrada final del backlog con resumen ejecutivo y lista ordenada de historias',
        produces: 'pm/backlog-entry',
      },
    ],
  },
  {
    id: 'architect',
    colorVar: 'var(--dk-c-architect)',
    nameEn: 'Architect',
    nameEs: 'Arquitecto',
    roleEn: 'Technical decisions',
    roleEs: 'Decisiones técnicas',
    artifactEn: 'ADR + system design',
    artifactEs: 'ADR + diseño del sistema',
    steps: [
      {
        slug: 'load-constraints',
        purposeEn: 'Read platform context and classify constraints as HARD (non-negotiable), SOFT (preferences), or OPPORTUNITIES',
        purposeEs: 'Lee el contexto de la plataforma y clasifica las restricciones como DURAS (no negociables), BLANDAS (preferencias) u OPORTUNIDADES',
        produces: 'architect/constraints-analysis',
      },
      {
        slug: 'evaluate-approaches',
        purposeEn: 'Compare 2–3 viable approaches; choose one with explicit rationale; document why alternatives were rejected',
        purposeEs: 'Compara 2–3 enfoques viables; elige uno con justificación explícita; documenta por qué se rechazaron las alternativas',
        produces: 'architect/approaches',
      },
      {
        slug: 'decision-record',
        purposeEn: 'Write the ADR: context, decision, alternatives, and consequences — including negatives. Positives-only = incomplete',
        purposeEs: 'Escribe el ADR: contexto, decisión, alternativas y consecuencias — incluyendo las negativas. Solo positivas = incompleto',
        produces: 'architect/adr',
      },
      {
        slug: 'system-design',
        purposeEn: 'Define data model, API surface (method/inputs/errors), service boundaries, and the happy-path sequence',
        purposeEs: 'Define el modelo de datos, la superficie de API (método/entradas/errores), los límites de servicio y la secuencia del flujo feliz',
        produces: 'architect/system-design',
      },
      {
        slug: 'risk-analysis',
        purposeEn: 'Identify top 3–5 risks (performance, security, reliability, coupling, migration) with concrete mitigations',
        purposeEs: 'Identifica los 3–5 riesgos principales (rendimiento, seguridad, fiabilidad, acoplamiento, migración) con mitigaciones concretas',
        produces: 'architect/risks',
      },
      {
        slug: 'technical-handoff',
        purposeEn: 'Consolidate all architectural work into two final artifacts for Developer and QA',
        purposeEs: 'Consolida todo el trabajo arquitectónico en dos artefactos finales para Developer y QA',
        produces: 'architectural-decision + system-design',
      },
    ],
  },
  {
    id: 'developer',
    colorVar: 'var(--dk-c-developer)',
    nameEn: 'Developer',
    nameEs: 'Desarrollador',
    roleEn: 'Builds the code',
    roleEs: 'Construye el código',
    artifactEn: 'Implementation',
    artifactEs: 'Implementación',
    steps: [
      {
        slug: 'explore',
        purposeEn: 'Read affected files and modules; map naming patterns and constraints before designing anything',
        purposeEs: 'Lee los archivos y módulos afectados; mapea patrones de nombres y restricciones antes de diseñar nada',
        produces: 'developer/dev-exploration',
      },
      {
        slug: 'spec',
        purposeEn: 'Define scope boundary, answer open questions, write acceptance criteria in Given/When/Then format',
        purposeEs: 'Define el límite de alcance, responde preguntas abiertas y escribe criterios de aceptación en formato Given/When/Then',
        produces: 'developer/dev-spec',
      },
      {
        slug: 'design',
        purposeEn: 'Choose technical approach, define data model and API shape, list key implementation constraints',
        purposeEs: 'Elige el enfoque técnico, define el modelo de datos y la forma de la API, lista las restricciones clave de implementación',
        produces: 'developer/dev-design',
      },
      {
        slug: 'implement',
        purposeEn: 'Write code for each task — plan-only mode (snippets) or writing mode (real files to declared targets)',
        purposeEs: 'Escribe código para cada tarea — modo solo-plan (fragmentos) o modo escritura (archivos reales en los destinos declarados)',
        produces: 'developer/dev-implementation',
      },
    ],
  },
  {
    id: 'qa',
    colorVar: 'var(--dk-c-qa)',
    nameEn: 'QA',
    nameEs: 'QA',
    roleEn: 'Quality gate',
    roleEs: 'Control de calidad',
    artifactEn: 'Test plan + quality report',
    artifactEs: 'Plan de pruebas + reporte de calidad',
    steps: [
      {
        slug: 'load-requirements',
        purposeEn: 'Extract and normalize acceptance criteria from upstream artifacts into Given/When/Then format',
        purposeEs: 'Extrae y normaliza los criterios de aceptación de los artefactos anteriores al formato Given/When/Then',
        produces: 'qa/ac-list',
      },
      {
        slug: 'ac-validation',
        purposeEn: 'Review each AC for atomicity, measurability, independence, completeness, unambiguity — rewrite failing ones',
        purposeEs: 'Revisa cada AC por atomicidad, medibilidad, independencia, completitud y claridad — reescribe los que fallen',
        produces: 'qa/ac-gaps',
      },
      {
        slug: 'edge-case-analysis',
        purposeEn: 'Discover edge cases via boundary values, equivalence partitioning, state transitions, concurrent access, permission boundaries',
        purposeEs: 'Descubre casos borde mediante valores límite, particiones de equivalencia, transiciones de estado, acceso concurrente y límites de permisos',
        produces: 'qa/edge-cases',
      },
      {
        slug: 'test-strategy',
        purposeEn: 'Define test pyramid: what each level covers, what it does not, test data strategy, and flakiness tolerance',
        purposeEs: 'Define la pirámide de pruebas: qué cubre cada nivel, qué no, la estrategia de datos de prueba y la tolerancia a inestabilidad',
        produces: 'qa/test-strategy',
      },
      {
        slug: 'test-case-generation',
        purposeEn: 'Write structured test specs (Given/When/Then) for happy path, validated ACs, and critical/high edge cases',
        purposeEs: 'Escribe specs de prueba estructuradas (Given/When/Then) para el flujo feliz, los ACs validados y los casos borde críticos/altos',
        produces: 'qa/test-cases',
      },
      {
        slug: 'quality-report',
        purposeEn: 'Verify AC coverage, compute percentage, render READY / READY WITH CAVEATS / BLOCKED verdict',
        purposeEs: 'Verifica la cobertura de ACs, calcula el porcentaje y emite el veredicto LISTO / LISTO CON RESERVAS / BLOQUEADO',
        produces: 'test-plan',
      },
      {
        slug: 'review',
        purposeEn: 'Go/no-go shipping verdict — holistic release-readiness decision integrating all QA findings before knowledge preservation',
        purposeEs: 'Veredicto final de envío — decisión holística de preparación para release que integra todos los hallazgos de QA antes de preservar el conocimiento',
        produces: 'qa/qa-review',
      },
    ],
  },
  {
    id: 'security',
    colorVar: 'var(--dk-c-security)',
    nameEn: 'Security',
    nameEs: 'Seguridad',
    roleEn: 'Threats — runs at any point',
    roleEs: 'Amenazas — puede correr en cualquier momento',
    artifactEn: 'Threat model + hardening checklist',
    artifactEs: 'Modelo de amenazas + checklist de endurecimiento',
    steps: [
      {
        slug: 'threat-modeling',
        purposeEn: 'Apply STRIDE: Spoofing, Tampering, Repudiation, Information Disclosure, DoS, Elevation of Privilege',
        purposeEs: 'Aplica STRIDE: Suplantación, Manipulación, Repudio, Divulgación de información, DoS y Elevación de privilegios',
        produces: 'security/stride-threats',
      },
      {
        slug: 'attack-surface',
        purposeEn: 'Map entry points, trust boundaries, data flows — verify validation/sanitization/encoding at each step',
        purposeEs: 'Mapea puntos de entrada, límites de confianza y flujos de datos — verifica validación/sanitización/codificación en cada paso',
        produces: 'security/attack-surface',
      },
      {
        slug: 'owasp-analysis',
        purposeEn: 'Check all 10 OWASP Top 10 categories (A01–A10) as APPLICABLE/NOT APPLICABLE and MITIGATED/AT RISK',
        purposeEs: 'Revisa las 10 categorías del OWASP Top 10 (A01–A10) como APLICABLE/NO APLICABLE y MITIGADO/EN RIESGO',
        produces: 'security/owasp-findings',
      },
      {
        slug: 'hardening-checklist',
        purposeEn: 'Deduplicate findings, prioritize by severity, group by effort: quick wins / medium / significant',
        purposeEs: 'Deduplica los hallazgos, prioriza por severidad y agrupa por esfuerzo: victorias rápidas / medio / significativo',
        produces: 'security-findings + hardening-checklist',
      },
    ],
  },
  {
    id: 'ux-ui',
    colorVar: 'var(--dk-c-ux)',
    nameEn: 'UX/UI',
    nameEs: 'UX/UI',
    roleEn: 'The experience',
    roleEs: 'La experiencia',
    artifactEn: 'UX brief + component spec',
    artifactEs: 'Brief de UX + spec de componentes',
    steps: [
      {
        slug: 'feature-brief',
        purposeEn: 'Identify primary actor, define core problem (not solution), establish 3–5 observable success criteria',
        purposeEs: 'Identifica el actor principal, define el problema central (no la solución) y establece 3–5 criterios de éxito observables',
        produces: 'ux-ui/feature-brief',
      },
      {
        slug: 'information-architecture',
        purposeEn: 'Organize content into sections, prioritize immediate vs. progressive disclosure, define navigation path',
        purposeEs: 'Organiza el contenido en secciones, prioriza la divulgación inmediata frente a la progresiva y define el camino de navegación',
        produces: 'ux-ui/ia',
      },
      {
        slug: 'user-flows',
        purposeEn: 'Map happy path, error path, and 2–3 edge case flows as numbered steps from the actor\'s perspective',
        purposeEs: 'Mapea el flujo feliz, el flujo de error y 2–3 flujos de casos borde como pasos numerados desde la perspectiva del actor',
        produces: 'ux-ui/flows',
      },
      {
        slug: 'component-mapping',
        purposeEn: 'Classify each UI state as reuse / extend / new — quality gate: >2:1 reuse ratio required',
        purposeEs: 'Clasifica cada estado de UI como reutilizar / extender / nuevo — compuerta de calidad: se requiere una relación de reúso >2:1',
        produces: 'ux-ui/components',
      },
      {
        slug: 'ux-handoff',
        purposeEn: 'Consolidate all UX work into ux-brief (flows + IA) and component-spec (inventory + props + events)',
        purposeEs: 'Consolida todo el trabajo de UX en ux-brief (flujos + IA) y component-spec (inventario + props + eventos)',
        produces: 'ux-brief + component-spec',
      },
    ],
  },
]

interface Props {
  lang: 'en' | 'es'
}

export function SpecialistExplorer({ lang }: Props) {
  const [activeSpecialist, setActiveSpecialist] = useState<string>('researcher')
  const isEs = lang === 'es'

  const active = PIPELINE_SPECIALISTS.find((s) => s.id === activeSpecialist)!

  function handleNodeClick(id: string) {
    if (id !== activeSpecialist) setActiveSpecialist(id)
    // clicking the active node is a no-op — prevents re-creating blank panel
  }

  function handleKeyDown(e: KeyboardEvent, id: string) {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'Spacebar') {
      e.preventDefault()
      handleNodeClick(id)
    } else if (e.key === 'Escape') {
      // LOAD-BEARING: stop the browser default (exit fullscreen) firing.
      e.stopPropagation()
      e.preventDefault()
      setActiveSpecialist('researcher')
    }
  }

  const framing = isEs
    ? 'Selecciona un especialista para explorar'
    : 'Select a specialist to explore'
  const groupLabel = isEs ? 'Pipeline de especialistas' : 'Specialist pipeline'
  const producesLabel = isEs ? 'Produce:' : 'Produces:'

  return (
    <div class="se" data-island>
      <div class="se-nodes" role="group" aria-label={groupLabel}>
        {PIPELINE_SPECIALISTS.map((s, i) => {
          const isActive = s.id === activeSpecialist
          return (
            <>
              {i > 0 && (
                <span aria-hidden="true" class="se-gate">
                  gate
                </span>
              )}
              <button
                key={s.id}
                type="button"
                class={`se-node${isActive ? ' se-node--active' : ''}`}
                data-specialist={s.id}
                aria-expanded={isActive}
                style={{ '--sc': s.colorVar } as Record<string, string>}
                onClick={() => handleNodeClick(s.id)}
                onKeyDown={(e) => handleKeyDown(e as unknown as KeyboardEvent, s.id)}
              >
                {isEs ? s.nameEs : s.nameEn}
              </button>
            </>
          )
        })}
      </div>

      <p class="se-framing">{framing}</p>

      {active && (
        <div
          class="se-panel se-panel--open"
          role="region"
          aria-label={isEs ? `${active.nameEs} detalles` : `${active.nameEn} details`}
          aria-live="polite"
          style={{ '--sc': active.colorVar } as Record<string, string>}
        >
          <h3 class="se-panel-name" style={{ color: active.colorVar } as Record<string, string>}>
            {isEs ? active.nameEs : active.nameEn}
          </h3>
          <p class="se-panel-role">{isEs ? active.roleEs : active.roleEn}</p>
          <div class="se-steps-table">
            {active.steps.map((step) => (
              <div key={step.slug} class="se-step-row">
                <code class="se-step-id">{step.slug}</code>
                <p class="se-step-purpose">{isEs ? step.purposeEs : step.purposeEn}</p>
                <span class="se-step-produces">→ {step.produces}</span>
              </div>
            ))}
          </div>
          <p class="se-panel-artifact">
            <span class="se-artifact-label">{producesLabel}</span>{' '}
            {isEs ? active.artifactEs : active.artifactEn}
          </p>
        </div>
      )}
    </div>
  )
}
