import { describe, expect, test } from "vitest"

import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

import { BodyTranslationPhasePresenter } from "./body-translation-phase.presenter"

interface TestScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: BodyTranslationOutputReadinessResponse | null
  errorMessage: string
  pendingAction:
    | "refresh"
    | "start"
    | "pause"
    | "resume"
    | "retry"
    | "cancel"
    | "check-output-readiness"
    | null
  hasLoaded: boolean
}

function createState(
  overrides: Partial<TestScreenState> = {}
): TestScreenState {
  return {
    jobId: 9,
    phase: "ready",
    summary: {
      jobId: 9,
      currentPhase: "body_translation",
      phaseState: "running",
      phaseRunId: 19,
      progress: {
        percent: 50,
        processedCount: 5,
        totalCount: 10,
        targetCount: 10,
        translatedCount: 3,
        skippedCount: 1,
        currentStep: "provider_request"
      },
      inputSummary: {
        targetCount: 10,
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "cred",
        provider: "fake",
        model: "model",
        executionMode: "batch",
        requestUnitCount: 10,
        outputCount: 3
      },
      actionEnablement: {
        canStart: false,
        canPause: true,
        canResume: false,
        canRetry: false,
        canCancel: true,
        canCheckOutputReadiness: true
      },
      outputReadiness: {
        ready: false,
        blockedReason: "phase_running",
        completedFieldCount: 3,
        statusConsistent: true
      }
    },
    outputReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    ...overrides
  }
}

describe("BodyTranslationPhasePresenter", () => {
  test("初回 loading は loading view を返す", () => {
    const presenter = new BodyTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(
      createState({ phase: "loading", summary: null, hasLoaded: false }),
      true
    )

    expect(viewModel.viewState).toBe("loading")
    expect(viewModel.statusTitle).toBe("summary を取得中")
    expect(viewModel.isLoading).toBe(true)
  })

  test("retryable error summary は recoverable_failed view を返す", () => {
    const presenter = new BodyTranslationPhasePresenter()
    const base = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...base.summary!,
          phaseState: "failed",
          errorSummary: {
            errorKind: "provider_failure",
            reason: "network timeout",
            retryable: true,
            isRedacted: true
          }
        }
      }),
      true
    )

    expect(viewModel.viewState).toBe("recoverable_failed")
    expect(viewModel.retryableLabel).toBe("再試行可能")
  })

  test("pending phase state は開始待ちとして表示する", () => {
    const presenter = new BodyTranslationPhasePresenter()
    const base = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...base.summary!,
          phaseState: "pending",
          progress: {
            ...base.summary!.progress,
            currentStep: "pending"
          },
          actionEnablement: {
            ...base.summary!.actionEnablement,
            canStart: true
          }
        }
      }),
      true
    )

    expect(viewModel.viewState).toBe("ready")
    expect(viewModel.phaseStateLabel).toBe("開始待ち")
    expect(viewModel.progressDetail).toContain("開始待ち")
    expect(viewModel.progressDetail).not.toContain("pending")
  })

  test("output readiness の上書き値を表示に使う", () => {
    const presenter = new BodyTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(
      createState({
        outputReadiness: {
          jobId: 9,
          currentPhase: "body_translation",
          phaseState: "running",
          ready: true,
          blockedReason: "",
          completedFieldCount: 10,
          statusConsistent: true,
          outputCount: 10
        }
      }),
      true
    )

    expect(viewModel.outputReadinessLabel).toBe("ready")
    expect(viewModel.outputReadinessCompletedFieldCountLabel).toBe("10 件")
  })

  test("result summary 不在時は progress から failed count を補完する", () => {
    const presenter = new BodyTranslationPhasePresenter()
    const base = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...base.summary!,
          resultSummary: undefined,
          progress: {
            ...base.summary!.progress,
            processedCount: 7,
            translatedCount: 4,
            skippedCount: 2
          }
        }
      }),
      true
    )

    expect(viewModel.failedCountLabel).toBe("1 件")
  })
})
