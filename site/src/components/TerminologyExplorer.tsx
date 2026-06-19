import { useState } from 'preact/hooks'

interface Term {
  id: string
  icon: string
  accentVar: string
  nameEn: string
  nameEs: string
  defEn: string
  defEs: string
  exEn: string
  exEs: string
}

const TERMS: Term[] = [
  {
    id: 'artifact',
    icon: '▣',
    accentVar: 'var(--dk-c-researcher)',
    nameEn: 'Artifact',
    nameEs: 'Artefacto',
    defEn: 'A written output a step leaves behind — the thing that survives a context reset.',
    defEs: 'Un resultado escrito que un step deja atrás — lo que sobrevive a un reset de contexto.',
    exEn: 'The Architect’s ADR is an artifact — it exists after the conversation ends.',
    exEs: 'El ADR del Architect es un artefacto — existe después de que la conversación termina.',
  },
  {
    id: 'brief',
    icon: '✦',
    accentVar: 'var(--dk-c-researcher)',
    nameEn: 'Brief',
    nameEs: 'Brief',
    defEn: 'A short, structured document handed from one step to the next.',
    defEs: 'Un documento corto y estructurado que se pasa de un step al siguiente.',
    exEn: 'The PM writes a brief; the Architect reads it to start.',
    exEs: 'El PM escribe un brief; el Architect lo lee para arrancar.',
  },
  {
    id: 'step',
    icon: '→',
    accentVar: 'var(--dk-c-researcher)',
    nameEn: 'Step',
    nameEs: 'Step',
    defEn: 'One unit of work: reads one artifact, produces one artifact.',
    defEs: 'Una unidad de trabajo: lee un artefacto, produce un artefacto.',
    exEn: '“feature-brief” is a step inside the UX/UI specialist.',
    exEs: '«feature-brief» es un step dentro del especialista de UX/UI.',
  },
  {
    id: 'specialist',
    icon: '◎',
    accentVar: 'var(--dk-c-pm)',
    nameEn: 'Specialist',
    nameEs: 'Especialista',
    defEn: 'One role, one lens — Researcher, Architect, Developer, and more.',
    defEs: 'Un rol, una mirada — Researcher, Architect, Developer y más.',
    exEn: 'Seven specialists = seven different ways of looking at the same feature.',
    exEs: 'Siete especialistas = siete formas distintas de mirar la misma funcionalidad.',
  },
  {
    id: 'orchestrator',
    icon: '⟳',
    accentVar: 'var(--dk-c-pm)',
    nameEn: 'Orchestrator',
    nameEs: 'Orquestador',
    defEn: 'Runs a specialist’s steps in order. Coordinates, never codes directly.',
    defEs: 'Corre los steps de un especialista en orden. Coordina, nunca programa directamente.',
    exEn: 'The UX orchestrator runs: feature-brief → IA → flows → components → handoff.',
    exEs: 'El orquestador de UX corre: feature-brief → IA → flujos → componentes → handoff.',
  },
  {
    id: 'router',
    icon: '⌥',
    accentVar: 'var(--dk-c-pm)',
    nameEn: 'Router',
    nameEs: 'Router',
    defEn: 'Reads your request and recommends which specialist fits.',
    defEs: 'Lee tu pedido y recomienda qué especialista encaja.',
    exEn: '/asdt add login → Router: “this sounds like Architect + Developer.”',
    exEs: '/asdt agregar login → Router: «esto suena a Architect + Developer».',
  },
  {
    id: 'gate',
    icon: '⊘',
    accentVar: 'var(--dk-c-qa)',
    nameEn: 'Gate',
    nameEs: 'Gate',
    defEn: 'A pause for a human to review and approve before work continues.',
    defEs: 'Una pausa para que un humano revise y apruebe antes de continuar.',
    exEn: 'PM hands off to Architect. You open the gate — not the AI.',
    exEs: 'El PM le pasa el trabajo al Architect. Tú abres el gate — no la IA.',
  },
]

interface Props {
  lang: 'en' | 'es'
}

export function TerminologyExplorer({ lang }: Props) {
  const [openIndex, setOpenIndex] = useState<number | null>(0)
  const isEs = lang === 'es'

  function toggle(i: number) {
    setOpenIndex((prev) => (prev === i ? null : i))
  }

  function handleKeyDown(e: KeyboardEvent, i: number) {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'Spacebar') {
      e.stopPropagation()
      e.preventDefault()
      toggle(i)
    } else if (e.key === 'Escape' && openIndex !== null) {
      // LOAD-BEARING: stop the browser default (exit fullscreen) firing.
      e.stopPropagation()
      e.preventDefault()
      setOpenIndex(null)
    }
  }

  const defLabel = isEs ? 'Definición' : 'Definition'
  const exLabel = isEs ? 'Ejemplo' : 'Example'

  return (
    <div
      class="te"
      data-island
      role="list"
      aria-label={isEs ? 'Conceptos clave de ASDT' : 'Key ASDT terminology'}
    >
      {TERMS.map((t, i) => {
        const isOpen = openIndex === i
        const bodyId = `te-body-${t.id}`
        return (
          <div
            key={t.id}
            class={`te-card${isOpen ? ' te-card--open' : ''}`}
            role="listitem"
            style={{ '--card-accent': t.accentVar } as Record<string, string>}
          >
            <button
              type="button"
              class="te-card-trigger"
              aria-expanded={isOpen}
              aria-controls={bodyId}
              onClick={() => toggle(i)}
              onKeyDown={(e) => handleKeyDown(e as unknown as KeyboardEvent, i)}
            >
              <span class="te-card-head">
                <span class="te-icon" aria-hidden="true">{t.icon}</span>
                <span class="te-term">{isEs ? t.nameEs : t.nameEn}</span>
                <span class="te-chevron" aria-hidden="true">{isOpen ? '−' : '+'}</span>
              </span>
            </button>
            <div id={bodyId} class="te-card-body" aria-hidden={!isOpen}>
              <div class="te-divider" />
              <p class="te-def">
                <span class="te-def-label">{defLabel}</span> {isEs ? t.defEs : t.defEn}
              </p>
              <blockquote class="te-example">
                <span class="te-def-label">{exLabel}</span> {isEs ? t.exEs : t.exEn}
              </blockquote>
            </div>
          </div>
        )
      })}
    </div>
  )
}
