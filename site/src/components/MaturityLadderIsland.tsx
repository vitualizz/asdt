import { useState, useEffect, useRef } from 'preact/hooks'

const RW = 180
const RH = 68

interface Props {
  lang: 'en' | 'es'
}

// Per-rung definition text — swaps with revealedStep. Replaces the MDX formula
// paragraph (step 2) and the removed bottomLabel SVG text.
const DEFINITIONS = [
  {
    en: 'Prompt, paste, pray. No spec, no trace, no second opinion. Feels fast until something breaks.',
    es: 'Escribe, pega, reza. Sin spec, sin trazabilidad, sin segunda opinión. Parece rápido hasta que algo falla.',
  },
  {
    en: 'Spec-driven development. Write the spec first; let the code follow the document. Disciplined but still a single lens.',
    es: 'Desarrollo guiado por especificación. Primero el spec, luego el código. Disciplinado, pero sigue siendo una sola perspectiva.',
  },
  {
    en: 'SDD on steroids. Seven specialists, each a fresh context, with a human checkpoint between every handoff. Artifacts survive every reset.',
    es: 'SDD con esteroides. Siete especialistas, cada uno con contexto fresco, con un punto de control humano entre cada entrega. Los artefactos sobreviven cada reset.',
  },
]

// Step-reveal of 3 maturity rungs: Vibe Coding (step 0) → SDD (step 1) → ASDT
// (step 2). Mirrors the static MaturityLadderDiagram.astro SVG body, wrapping
// each rung in a <g class="rung"> so island state can drive per-rung opacity.
export function MaturityLadderIsland({ lang }: Props) {
  const [revealedStep, setRevealedStep] = useState(0)
  const isEs = lang === 'es'

  const opacityFor = (rung: number) => (rung <= revealedStep ? 1 : 0.15)
  const connectorOpacity = revealedStep === 0 ? 0.1 : revealedStep === 1 ? 0.6 : 1

  const rootRef = useRef<HTMLDivElement>(null)
  const motion = useRef(
    typeof window !== 'undefined'
      ? !window.matchMedia('(prefers-reduced-motion: reduce)').matches
      : true
  )

  // Expose the fragment API on the DOM node so the DeckController can drive
  // step-reveal via global next/prev. Re-assigned on each state change so the
  // closures stay fresh (deps include revealedStep).
  useEffect(() => {
    const el = rootRef.current
    if (!el) return
    ;(el as any).fragmentAdvance = () => {
      if (revealedStep >= 2) return false // already at last; navigate
      setRevealedStep((s) => s + 1)
      return true // consumed a step; stay
    }
    ;(el as any).fragmentBack = () => {
      if (revealedStep <= 0) return false // already at first; navigate
      setRevealedStep((s) => s - 1)
      return true // consumed a step; stay
    }
    ;(el as any).fragmentReset = () => setRevealedStep(0)
    ;(el as any).fragmentShowAll = () => setRevealedStep(2)
  }, [revealedStep])

  return (
    <div class="ml-island" data-fragment-island ref={rootRef}>
      <svg
        class="ml-diagram"
        viewBox="0 0 540 300"
        role="img"
        aria-labelledby="ml-title ml-desc"
        xmlns="http://www.w3.org/2000/svg"
      >
        <title id="ml-title">ASDT maturity ladder: Vibe Coding to SDD to ASDT</title>
        <desc id="ml-desc">
          Three ascending rungs showing the progression from unstructured vibe
          coding through spec-driven development to ASDT with specialist teams
          and human gates.
        </desc>
        <defs>
          <filter id="ml-glow" x="-40%" y="-40%" width="180%" height="180%">
            <feDropShadow dx="0" dy="0" stdDeviation="6" flood-color="#B7CC85" flood-opacity="0.25" />
          </filter>
        </defs>

        {/* connectors */}
        <line
          x1={30 + RW} y1="212" x2="250" y2={122 + RH}
          stroke="#313342" stroke-width="2" stroke-dasharray="4 3"
          style={{ opacity: connectorOpacity, transition: motion.current ? 'opacity 300ms ease' : 'none' }}
        />
        <line
          x1={160 + RW} y1="122" x2="380" y2={32 + RH}
          stroke="#313342" stroke-width="2" stroke-dasharray="4 3"
          style={{ opacity: revealedStep >= 2 ? 1 : 0.1, transition: motion.current ? 'opacity 300ms ease' : 'none' }}
        />

        {/* Rung 1: Vibe Coding */}
        <g class="rung" style={{ opacity: opacityFor(0), transition: motion.current ? 'opacity 300ms ease' : 'none' }}>
          <rect x="30" y="212" width={RW} height={RH} rx="10" fill="#232A36" stroke="#313342" stroke-width="1" />
          <text x={30 + RW / 2} y="246" text-anchor="middle" fill="#F3F6F9" font-family="var(--dk-font-display, sans-serif)" font-size="14" font-weight="700">Vibe Coding</text>
          <text x={30 + RW / 2} y="266" text-anchor="middle" fill="#A1AABB" font-family="var(--dk-font-body, serif)" font-size="11">{isEs ? 'escribe, pega, reza' : 'prompt, paste, pray'}</text>
        </g>

        {/* Rung 2: SDD */}
        <g class="rung" style={{ opacity: opacityFor(1), transition: motion.current ? 'opacity 300ms ease' : 'none' }}>
          <rect x="160" y="122" width={RW} height={RH} rx="10" fill="#232A36" stroke="#7FB4CA" stroke-width="1.5" />
          <text x={160 + RW / 2} y="156" text-anchor="middle" fill="#7FB4CA" font-family="var(--dk-font-display, sans-serif)" font-size="14" font-weight="700">SDD</text>
          <text x={160 + RW / 2} y="176" text-anchor="middle" fill="#A1AABB" font-family="var(--dk-font-body, serif)" font-size="11">{isEs ? 'guiado por especificación' : 'spec-driven, disciplined'}</text>
        </g>

        {/* Rung 3: ASDT (glowing) */}
        <g class="rung" style={{ opacity: opacityFor(2), transition: motion.current ? 'opacity 300ms ease' : 'none' }}>
          <rect x="290" y="32" width={RW} height={RH} rx="10" fill="#232A36" stroke="#B7CC85" stroke-width="2" filter="url(#ml-glow)" />
          <rect x={290 + RW / 2 - 55} y="6" width="110" height="20" rx="4" fill="#B7CC85" />
          <text x={290 + RW / 2} y="20" text-anchor="middle" fill="#1c212c" font-family="var(--dk-font-mono, monospace)" font-size="10" font-weight="700">{isEs ? 'ESTÁS AQUÍ' : 'YOU ARE HERE'}</text>
          <text x={290 + RW / 2} y="62" text-anchor="middle" fill="#B7CC85" font-family="var(--dk-font-display, sans-serif)" font-size="16" font-weight="700">ASDT</text>
          <text x={290 + RW / 2} y="82" text-anchor="middle" fill="#B7CC85" opacity="0.85" font-family="var(--dk-font-body, serif)" font-size="11">{isEs ? 'SDD con esteroides' : 'SDD on steroids'}</text>
        </g>
      </svg>

      <div class="ml-definition" aria-live="polite">
        <p class="ml-def-text">{isEs ? DEFINITIONS[revealedStep].es : DEFINITIONS[revealedStep].en}</p>
      </div>
    </div>
  )
}
