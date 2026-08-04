import { useState, useEffect, useRef } from 'preact/hooks'
import { SPECIALIST_COLOR } from '@data/artifact-graph'
import type { SpecialistId } from '@data/artifact-graph'

const ID_TO_COLOR: Record<string, string> = {
  asdt: 'var(--color-accent)',
  ...Object.fromEntries(
    Object.entries(SPECIALIST_COLOR).map(([id, cssVar]) => [id, `var(${cssVar})`])
  ),
}

interface FlowWalkerStep {
  command: string
  specialistId: SpecialistId | 'asdt'
  description: string
}

interface FragmentIsland extends HTMLDivElement {
  fragmentAdvance: () => boolean
  fragmentBack: () => boolean
  fragmentReset: () => void
  fragmentShowAll: () => void
}

interface Props {
  title: string
  steps: FlowWalkerStep[]
  labelStep: string
  labelOf: string
}

export function FlowWalker({ title, steps, labelStep, labelOf }: Props) {
  const [step, setStep] = useState(0)
  const current = steps[step]
  const total = steps.length

  const colorVar = ID_TO_COLOR[current.specialistId] ?? 'var(--color-accent)'

  const rootRef = useRef<HTMLDivElement>(null)

  // Expose the fragment API on the DOM node so the DeckController can drive
  // step-stepping via global next/prev. Re-assigned on each state change so the
  // closures stay fresh (deps include step, total).
  useEffect(() => {
    const el = rootRef.current as FragmentIsland | null
    if (!el) return
    el.fragmentAdvance = () => {
      if (step >= total - 1) return false // already at last; navigate
      setStep((s) => s + 1)
      return true // consumed a step; stay
    }
    el.fragmentBack = () => {
      if (step <= 0) return false // already at first; navigate
      setStep((s) => s - 1)
      return true // consumed a step; stay
    }
    el.fragmentReset = () => setStep(0)
    el.fragmentShowAll = () => setStep(total - 1)
  }, [step, total])

  return (
    <div class="fw" data-fragment-island ref={rootRef} aria-label={title} style={{ '--sc': colorVar } as Record<string, string>}>
      <div class="fw-header">
        <span class="fw-counter">{labelStep} {step + 1} {labelOf} {total}</span>
      </div>
      <code class="fw-command">{current.command}</code>
      <p class="fw-desc">{current.description}</p>
      <div class="fw-dots" role="presentation">
        {steps.map((_, i) => (
          <span
            key={i}
            class={`fw-dot${i === step ? ' fw-dot--active' : ''}`}
            aria-hidden="true"
          />
        ))}
      </div>
    </div>
  )
}
