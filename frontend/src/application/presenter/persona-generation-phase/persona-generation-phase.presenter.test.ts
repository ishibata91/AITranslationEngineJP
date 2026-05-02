import { describe, expect, test } from "vitest"
import { PersonaGenerationPhasePresenter } from "./persona-generation-phase.presenter"

const presenter = new PersonaGenerationPhasePresenter()

describe("PersonaGenerationPhasePresenter", () => {
  test("未接続状態を view model に反映する", () => {
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready",
        summary: null,
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: false
      },
      false
    )

    expect(vm.gatewayStatus).toContain("未接続")
    expect(vm.viewState).toBe("loading")
  })

  test("snapshot missing を優先表示する", () => {
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready",
        summary: {
          jobId: 10,
          currentPhase: "persona_generation",
          phaseState: "completed",
          progress: {
            percent: 100,
            processedCount: 1,
            totalCount: 1,
            targetCount: 1,
            currentStep: "completed"
          },
          targetSummary: {
            targetCount: 1,
            commonPersonaHitCount: 0,
            commonPersonaMissCount: 1,
            skippedCount: 0,
            skippedReasons: [],
            targetSnapshotDigest: "sha256:1"
          },
          execution: {
            credentialRef: "cred",
            provider: "fake",
            model: "m",
            executionMode: "single_request",
            promptDigest: "sha256:1",
            inputCount: 1,
            outputCount: 0,
            evidenceRefs: []
          },
          resultSummary: {
            generatedCount: 0,
            failedCount: 1,
            personaCount: 0,
            missingCount: 1,
            snapshotId: "",
            snapshotDigest: "sha256:1",
            snapshotReferenceStatus: "missing",
            bodyReadiness: false
          },
          errorSummary: {
            errorKind: "snapshot_missing",
            reason: "redacted",
            retryable: false,
            isRedacted: true
          },
          actionEnablement: {
            canStart: false,
            canPause: false,
            canResume: false,
            canRetry: true,
            canCancel: false,
            canStartBodyPhase: false
          }
        },
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true
      },
      true
    )

    expect(vm.viewState).toBe("snapshot_missing")
    expect(vm.errorKindLabel).toContain("snapshot_missing")
  })
})
