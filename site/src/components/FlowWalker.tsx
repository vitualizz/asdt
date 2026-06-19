import { useState, useEffect, useRef } from 'preact/hooks'
import { SPECIALIST_COLOR } from '@data/artifact-graph'
import type { PipelineFlow } from '@data/pipeline-flow-steps'

const ID_TO_COLOR: Record<string, string> = {
  asdt: 'var(--color-accent)',
  ...Object.fromEntries(
    Object.entries(SPECIALIST_COLOR).map(([id, cssVar]) => [id, `var(${cssVar})`])
  ),
}

interface Props {
  flow: PipelineFlow
  lang: 'en' | 'es'
  labelStep: string
  labelOf: string
}

export function FlowWalker({ flow, lang, labelStep, labelOf }: Props) {
  const [step, setStep] = useState(0)
  const current = flow.steps[step]
  const total = flow.steps.length
  const isEs = lang === 'es'

  const colorVar = ID_TO_COLOR[current.specialistId] ?? 'var(--color-accent)'

  const rootRef = useRef<HTMLDivElement>(null)

  // Expose the fragment API on the DOM node so the DeckController can drive
  // step-stepping via global next/prev. Re-assigned on each state change so the
  // closures stay fresh (deps include step, total).
  useEffect(() => {
    const el = rootRef.current
    if (!el) return
    ;(el as any).fragmentAdvance = () => {
      if (step >= total - 1) return false
      setStep((s) => s + 1)
      return step + 1 < total - 1
    }
    ;(el as any).fragmentBack = () => {
      if (step <= 0) return false
      setStep((s) => s - 1)
      return step - 1 > 0
    }
    ;(el as any).fragmentReset = () => setStep(0)
    ;(el as any).fragmentShowAll = () => setStep(total - 1)
  }, [step, total])

  return (
    <div class="fw" data-fragment-island ref={rootRef} style={{ '--sc': colorVar } as Record<string, string>}>
      <div class="fw-header">
        <span class="fw-counter">{labelStep} {step + 1} {labelOf} {total}</span>
      </div>
      <code class="fw-command">{current.command}</code>
      <p class="fw-desc">{isEs ? current.descriptionEs : current.descriptionEn}</p>
      <div class="fw-dots" role="presentation">
        {flow.steps.map((_, i) => (
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
