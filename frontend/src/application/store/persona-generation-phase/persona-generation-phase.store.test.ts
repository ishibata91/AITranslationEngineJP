import { describe, expect, test, vi } from "vitest"

import type {
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseSummaryResponse
} from "@application/gateway-contract/persona-generation-phase"

import { PersonaGenerationPhaseStore } from "./persona-generation-phase.store"

function createSummary(): PersonaGenerationPhaseSummaryResponse {
  return {
    jobId: 10,
    currentPhase: "persona_generation",
    phaseState: "ready",
    phaseRunId: 20,
    progress: {
      percent: 40,
      processedCount: 4,
      totalCount: 10,
      targetCount: 6,
      currentStep: "step"
    },
    targetSummary: {
      targetCount: 6,
      commonPersonaHitCount: 1,
      commonPersonaMissCount: 5,
      skippedCount: 1,
      skippedReasons: ["zero dialogue"],
      targetSnapshotId: "snapshot",
      targetSnapshotDigest: "digest"
    },
    execution: {
      credentialRef: "credential",
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "single_request",
      promptDigest: "prompt",
      inputCount: 6,
      outputCount: 4,
      evidenceRefs: ["evidence"]
    },
    resultSummary: {
      generatedCount: 4,
      failedCount: 0,
      personaCount: 4,
      missingCount: 0,
      snapshotId: "persona-snapshot",
      snapshotDigest: "persona-digest",
      snapshotReferenceStatus: "ready",
      bodyReadiness: true
    },
    errorSummary: {
      errorKind: "provider_failure",
      reason: "reason",
      retryable: true,
      isRedacted: true
    },
    actionEnablement: {
      canStart: false,
      canPause: true,
      canResume: false,
      canRetry: true,
      canCancel: true
    }
  }
}

function createBodyReadiness(): PersonaGenerationBodyReadinessResponse {
  return {
    jobId: 10,
    currentPhase: "persona_generation",
    phaseState: "ready",
    inputSummary: {
      personaCount: 4,
      missingCount: 0,
      snapshotId: "persona-snapshot",
      snapshotDigest: "persona-digest",
      evidenceRefs: ["evidence"]
    }
  }
}

describe("PersonaGenerationPhaseStore", () => {
  test("subscribe は初期 state を即時通知する", () => {
    const store = new PersonaGenerationPhaseStore()
    const listener = vi.fn()

    store.subscribe(listener)

    expect(listener).toHaveBeenCalledTimes(1)
    expect(listener.mock.calls[0]?.[0]).toMatchObject({
      jobId: null,
      phase: "idle",
      summary: null,
      bodyReadiness: null,
      pendingAction: null,
      hasLoaded: false
    })
  })

  test("unsubscribe 後は update しても通知されない", () => {
    const store = new PersonaGenerationPhaseStore()
    const listener = vi.fn()

    const unsubscribe = store.subscribe(listener)
    unsubscribe()
    store.update((draft) => {
      draft.jobId = 31
    })

    expect(listener).toHaveBeenCalledTimes(1)
  })

  test("snapshot は summary と body readiness の defensive copy を返す", () => {
    const store = new PersonaGenerationPhaseStore()

    store.update((draft) => {
      draft.summary = createSummary()
      draft.bodyReadiness = createBodyReadiness()
    })

    const snapshot = store.snapshot()
    snapshot.summary!.progress.percent = 0
    snapshot.summary!.targetSummary.skippedReasons[0] = "changed"
    if (snapshot.summary?.execution) {
      snapshot.summary.execution.evidenceRefs[0] = "changed"
    }
    snapshot.summary!.resultSummary!.generatedCount = 0
    snapshot.summary!.errorSummary!.reason = "changed"
    snapshot.summary!.actionEnablement.canPause = false
    snapshot.bodyReadiness!.inputSummary.evidenceRefs[0] = "changed"

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.summary?.progress.percent).toBe(40)
    expect(nextSnapshot.summary?.targetSummary.skippedReasons[0]).toBe(
      "zero dialogue"
    )
    expect(nextSnapshot.summary?.execution?.evidenceRefs[0]).toBe("evidence")
    expect(nextSnapshot.summary?.resultSummary?.generatedCount).toBe(4)
    expect(nextSnapshot.summary?.errorSummary?.reason).toBe("reason")
    expect(nextSnapshot.summary?.actionEnablement.canPause).toBe(true)
    expect(nextSnapshot.bodyReadiness?.inputSummary.evidenceRefs[0]).toBe(
      "evidence"
    )
  })
})
