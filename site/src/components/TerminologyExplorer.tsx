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
    nameEn: 'Document (artifact)',
    nameEs: 'Documento (artefacto)',
    defEn: 'A written file a step leaves behind — the thing that survives when the chat is gone.',
    defEs: 'Un archivo escrito que un paso deja atrás — lo que sobrevive cuando el chat ya no está.',
    exEn: 'The Architect’s ADR is a document — it still exists after the conversation ends.',
    exEs: 'El ADR del Architect es un documento — sigue existiendo después de que la conversación termina.',
  },
  {
    id: 'brief',
    icon: '✦',
    accentVar: 'var(--dk-c-researcher)',
    nameEn: 'Summary (brief)',
    nameEs: 'Resumen (brief)',
    defEn: 'A short, structured document handed from one step to the next.',
    defEs: 'Un documento corto y estructurado que se pasa de un paso al siguiente.',
    exEn: 'The PM writes a summary; the Architect reads it to start.',
    exEs: 'El PM escribe un resumen; el Architect lo lee para arrancar.',
  },
  {
    id: 'step',
    icon: '→',
    accentVar: 'var(--dk-c-researcher)',
    nameEn: 'Step',
    nameEs: 'Paso (step)',
    defEn: 'One unit of work: reads one document, produces one document.',
    defEs: 'Una unidad de trabajo: lee un documento, produce un documento.',
    exEn: '“feature-brief” is a step inside the UX/UI specialist.',
    exEs: '«feature-brief» es un paso dentro del especialista de UX/UI.',
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
    nameEn: 'Coordinator (orchestrator)',
    nameEs: 'Coordinador (orquestador)',
    defEn: 'Runs a specialist’s steps in order. Coordinates the work, never writes code directly.',
    defEs: 'Ejecuta los pasos de un especialista en orden. Coordina el trabajo, nunca programa directamente.',
    exEn: 'The UX coordinator runs: feature-brief → IA → flows → components → handoff.',
    exEs: 'El coordinador de UX ejecuta: feature-brief → IA → flujos → componentes → handoff.',
  },
  {
    id: 'router',
    icon: '⌥',
    accentVar: 'var(--dk-c-pm)',
    nameEn: 'Router',
    nameEs: 'Router',
    defEn: 'The part that reads your request and routes it — recommending which specialist fits.',
    defEs: 'La parte que lee tu pedido y lo reparte — recomendando qué especialista encaja.',
    exEn: '/asdt add login → Router: “this sounds like Architect + Developer.”',
    exEs: '/asdt agregar login → Router: «esto suena a Architect + Developer».',
  },
  {
    id: 'gate',
    icon: '⊘',
    accentVar: 'var(--dk-c-qa)',
    nameEn: 'Checkpoint',
    nameEs: 'Punto de control',
    defEn: 'A pause for a human to review and approve before work continues.',
    defEs: 'Una pausa para que un humano revise y apruebe antes de continuar.',
    exEn: 'PM hands off to Architect. You open the checkpoint — not the AI.',
    exEs: 'El PM le pasa el trabajo al Architect. Tú activas el punto de control — no la IA.',
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
