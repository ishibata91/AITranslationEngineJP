import { describe, expect, test } from "vitest"

import type {
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

import { TermTranslationPhasePresenter } from "./term-translation-phase.presenter"

interface TestScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: "refresh" | "start" | "pause" | "resume" | "retry" | null
  hasLoaded: boolean
}

function createState(
  overrides: Partial<TestScreenState> = {}
): TestScreenState {
  return {
    jobId: 7,
    phase: "ready",
    summary: {
      jobId: 7,
      currentPhase: "term_translation",
      phaseState: "ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 10,
        aiTargetCount: 5,
        currentStep: "ready"
      },
      totalTermCount: 10,
      dictionaryHitCount: 5,
      aiTargetCount: 5,
      execution: {
        credentialRef: "cred-main",
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      actionEnablement: {
        canStart: true,
        canPause: false,
        canResume: false,
        canRetry: false,
        canRefresh: true,
        canStartNextPhase: false
      }
    },
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    ...overrides
  }
}

describe("TermTranslationPhasePresenter", () => {
  test("ジョブ未選択時は blocked view と未選択メッセージを返す", () => {
    const presenter = new TermTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(
      createState({ jobId: null, summary: null }),
      false
    )

    expect(viewModel.viewState).toBe("blocked")
    expect(viewModel.statusTitle).toBe("ジョブ未選択")
    expect(viewModel.gatewayStatus).toBe("未接続")
  })

  test("retryable error がある時は recoverable_failed view を返す", () => {
    const presenter = new TermTranslationPhasePresenter()
    const baseState = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...baseState.summary!,
          phaseState: "failed",
          errorSummary: {
            errorKind: "provider_failure",
            reason: "network timeout",
            retryable: true,
            isRedacted: false
          }
        }
      }),
      true
    )

    expect(viewModel.viewState).toBe("recoverable_failed")
    expect(viewModel.retryableLabel).toBe("再試行可能")
  })

  test("完了かつ aiTargetCount が 0 の時は empty_completed view を返す", () => {
    const presenter = new TermTranslationPhasePresenter()
    const baseState = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...baseState.summary!,
          phaseState: "completed",
          aiTargetCount: 0
        }
      }),
      true
    )

    expect(viewModel.viewState).toBe("empty_completed")
    expect(viewModel.statusTitle).toBe("対象語なしで完了")
  })

  test("loading 中は action card を無効化し readiness の blocked reason を優先する", () => {
    const presenter = new TermTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(
      createState({
        phase: "loading",
        nextPhaseReadiness: {
          jobId: 7,
          currentPhase: "term_translation",
          phaseState: "completed",
          canStartNextPhase: false,
          blockedReason: "review pending"
        }
      }),
      true
    )

    expect(viewModel.actionCards.every((card) => card.disabled)).toBe(true)
    const nextPhaseCard = viewModel.actionCards.find(
      (card) => card.id === "next-phase"
    )
    expect(nextPhaseCard?.blockedReason).toBe("review pending")
  })
})
