import { describe, expect, test, vi } from "vitest"

import { TermTranslationPhaseStore } from "./term-translation-phase.store"

describe("TermTranslationPhaseStore", () => {
  test("subscribe は初期 state を即時通知する", () => {
    const store = new TermTranslationPhaseStore()
    const listener = vi.fn()

    store.subscribe(listener)

    expect(listener).toHaveBeenCalledTimes(1)
    expect(listener.mock.calls[0]?.[0]).toMatchObject({
      jobId: null,
      phase: "idle",
      summary: null,
      nextPhaseReadiness: null,
      pendingAction: null,
      hasLoaded: false
    })
  })

  test("unsubscribe 後は update しても通知されない", () => {
    const store = new TermTranslationPhaseStore()
    const listener = vi.fn()

    const unsubscribe = store.subscribe(listener)
    unsubscribe()
    store.update((draft) => {
      draft.jobId = 31
    })

    expect(listener).toHaveBeenCalledTimes(1)
  })

  test("snapshot は nested object の defensive copy を返す", () => {
    const store = new TermTranslationPhaseStore()

    store.update((draft) => {
      draft.summary = {
        jobId: 10,
        currentPhase: "term_translation",
        phaseState: "ready",
        progress: {
          percent: 40,
          processedCount: 4,
          totalCount: 10,
          aiTargetCount: 6,
          currentStep: "step"
        },
        totalTermCount: 10,
        dictionaryHitCount: 4,
        aiTargetCount: 6,
        execution: {
          credentialRef: "cred",
          provider: "openai-compatible",
          model: "gpt-4.1-mini",
          executionMode: "batch"
        },
        actionEnablement: {
          canStart: false,
          canPause: true,
          canResume: false,
          canRetry: false,
          canStartNextPhase: false
        }
      }
    })

    const snapshot = store.snapshot()
    snapshot.summary!.progress.percent = 0

    expect(store.snapshot().summary?.progress.percent).toBe(40)
  })
})
