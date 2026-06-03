import { describe, expect, test, vi } from "vitest"

import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

import { BodyTranslationPhaseStore } from "./body-translation-phase.store"

function createSummary(): BodyTranslationPhaseSummaryResponse {
  return {
    jobId: 10,
    currentPhase: "body_translation",
    phaseState: "ready",
    phaseRunId: 20,
    progress: {
      percent: 40,
      processedCount: 4,
      totalCount: 10,
      targetCount: 6,
      translatedCount: 3,
      skippedCount: 1,
      currentStep: "step"
    },
    inputSummary: {
      targetCount: 6,
      skippedReasons: ["exact dictionary"],
      inputSnapshotRef: "input",
      dictionaryDigest: "dictionary",
      personaDigest: "persona",
      metadataDigest: "metadata",
      promptDigest: "prompt"
    },
    requestSummary: {
      providerTargetCount: 5,
      exactDictionaryExclusionCount: 1,
      partialDictionaryConstraintCount: 2
    },
    execution: {
      credentialRef: "credential",
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMode: "single_request",
      requestUnitCount: 5,
      outputCount: 4
    },
    fieldResults: [
      {
        identity: { formId: "FE01A812", editorId: "TEST_FIELD" },
        translatedText: "translated"
      }
    ],
    resultSummary: {
      translatedCount: 3,
      failedCount: 0,
      skippedCount: 1,
      protectionFailedCount: 0,
      outputReadyCount: 3,
      outputCount: 3,
      fieldResults: [
        {
          identity: { formId: "FE01A813", editorId: "TEST_RESULT" },
          translatedText: "result"
        }
      ]
    },
    errorSummary: {
      errorKind: "provider_failure",
      reason: "reason",
      retryable: true,
      isRedacted: true
    },
    projection: {
      phaseLifecycle: "running",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      targetCount: 10,
      previousPhaseLifecycle: "completed",
      personaBodyReadiness: {
        bodyReadiness: true,
        snapshotReferenceStatus: "available"
      }
    },
    outputReadiness: {
      completedFieldCount: 3,
      statusConsistent: true,
      outputCount: 3
    }
  }
}

function createOutputReadiness(): BodyTranslationOutputReadinessResponse {
  return {
    jobId: 10,
    currentPhase: "body_translation",
    phaseState: "ready",
    ready: true,
    completedFieldCount: 3,
    statusConsistent: true,
    outputCount: 3
  }
}

describe("BodyTranslationPhaseStore", () => {
  test("subscribe は初期 state を即時通知する", () => {
    const store = new BodyTranslationPhaseStore()
    const listener = vi.fn()

    store.subscribe(listener)

    expect(listener).toHaveBeenCalledTimes(1)
    expect(listener.mock.calls[0]?.[0]).toMatchObject({
      jobId: null,
      phase: "idle",
      summary: null,
      outputReadiness: null,
      pendingAction: null,
      hasLoaded: false
    })
  })

  test("unsubscribe 後は update しても通知されない", () => {
    const store = new BodyTranslationPhaseStore()
    const listener = vi.fn()

    const unsubscribe = store.subscribe(listener)
    unsubscribe()
    store.update((draft) => {
      draft.jobId = 31
    })

    expect(listener).toHaveBeenCalledTimes(1)
  })

  test("snapshot は summary と output readiness の defensive copy を返す", () => {
    const store = new BodyTranslationPhaseStore()

    store.update((draft) => {
      draft.summary = createSummary()
      draft.outputReadiness = createOutputReadiness()
    })

    const snapshot = store.snapshot()
    snapshot.summary!.progress.percent = 0
    snapshot.summary!.inputSummary.skippedReasons![0] = "changed"
    snapshot.summary!.requestSummary!.providerTargetCount = 0
    if (snapshot.summary?.execution) {
      snapshot.summary.execution.provider = "changed"
    }
    snapshot.summary!.fieldResults![0].identity!.formId = "changed"
    snapshot.summary!.resultSummary!.fieldResults![0].identity!.formId =
      "changed"
    snapshot.summary!.errorSummary!.reason = "changed"
    snapshot.summary!.projection.phaseLifecycle = "changed"
    snapshot.summary!.outputReadiness.outputCount = 0
    snapshot.outputReadiness!.ready = false

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.summary?.progress.percent).toBe(40)
    expect(nextSnapshot.summary?.inputSummary.skippedReasons?.[0]).toBe(
      "exact dictionary"
    )
    expect(nextSnapshot.summary?.requestSummary?.providerTargetCount).toBe(5)
    expect(nextSnapshot.summary?.execution?.provider).toBe("gemini")
    expect(nextSnapshot.summary?.fieldResults?.[0]?.identity?.formId).toBe(
      "FE01A812"
    )
    expect(
      nextSnapshot.summary?.resultSummary?.fieldResults?.[0]?.identity?.formId
    ).toBe("FE01A813")
    expect(nextSnapshot.summary?.errorSummary?.reason).toBe("reason")
    expect(nextSnapshot.summary?.projection.phaseLifecycle).toBe("running")
    expect(nextSnapshot.summary?.outputReadiness.outputCount).toBe(3)
    expect(nextSnapshot.outputReadiness?.ready).toBe(true)
  })
})
