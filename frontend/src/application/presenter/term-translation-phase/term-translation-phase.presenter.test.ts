import { describe, expect, test } from "vitest"

import type {
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

import { TermTranslationPhasePresenter } from "./term-translation-phase.presenter"

// 等価性テスト用: backend 旧導出の期待値を仕様（detail-spec term-translation-phase-REQ-008）から固定する
// 成立条件: ジョブが終端でない・phaseState が completed・confirmedCount >= aiTargetCount
// 区別理由: terminal_job, term_phase_incomplete

interface TestScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: "start" | "pause" | "resume" | "retry" | null
  hasLoaded: boolean
  initialFetchDone: boolean
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
        canRetry: false
      }
    },
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    initialFetchDone: true,
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

  test("pending phase state は開始待ちとして表示する", () => {
    const presenter = new TermTranslationPhasePresenter()
    const baseState = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...baseState.summary!,
          phaseState: "pending",
          phaseRunId: 44,
          progress: {
            ...baseState.summary!.progress,
            currentStep: "pending"
          }
        }
      }),
      true
    )

    expect(viewModel.viewState).toBe("idle_ready")
    expect(viewModel.phaseStateLabel).toBe("開始待ち")
    expect(viewModel.progressDetail).toContain("開始待ち")
    expect(viewModel.progressDetail).not.toContain("pending")
  })

  test("進行詳細は AI 翻訳対象語件数を分母として表示する", () => {
    const presenter = new TermTranslationPhasePresenter()
    const baseState = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...baseState.summary!,
          progress: {
            ...baseState.summary!.progress,
            processedCount: 2,
            totalCount: 10,
            aiTargetCount: 5
          },
          totalTermCount: 10,
          aiTargetCount: 5
        }
      }),
      true
    )

    expect(viewModel.progressDetail).toBe(
      "2 / 5 件 / AI 翻訳対象語 5 件 / 開始可能"
    )
    expect(viewModel.totalTermCountLabel).toBe("10 件")
    expect(viewModel.aiTargetCountLabel).toBe("5 件")
  })

  test("current phase key は画面表示名へ変換する", () => {
    const presenter = new TermTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(createState(), true)

    expect(viewModel.currentPhaseLabel).toBe("単語翻訳")
  })

  test("loading 中は action card を無効化し next-phase blocked reason を返す", () => {
    const presenter = new TermTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(
      createState({
        phase: "loading",
        summary: {
          ...createState().summary!,
          phaseState: "pending"
        }
      }),
      true
    )

    expect(viewModel.actionCards.every((card) => card.disabled)).toBe(true)
    const nextPhaseCard = viewModel.actionCards.find(
      (card) => card.id === "next-phase"
    )
    expect(nextPhaseCard?.blockedReason).toBe(
      "単語翻訳段階が未完了のため次段階を開始できません。"
    )
  })
})

// UT-EQV-001/002: 次段階開始可否の等価性テスト（term-translation-phase-REQ-008 成立条件・区別理由）
// 期待値は detail-spec の成立条件（backend 旧導出撤去後の仕様根拠）から固定する
describe("TermTranslationPhasePresenter - 次段階開始可否の等価性（UT-EQV-001/002）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildSummary(
    overrides: Partial<TermTranslationPhaseSummaryResponse>
  ): TermTranslationPhaseSummaryResponse {
    return {
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
        credentialRef: "cred",
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      actionEnablement: {
        canStart: true,
        canPause: false,
        canResume: false,
        canRetry: false
      },
      ...overrides
    }
  }

  function buildState(
    summaryOverrides: Partial<TermTranslationPhaseSummaryResponse> = {}
  ) {
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: buildSummary(summaryOverrides),
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // UT-EQV-001: 成立条件を全て満たす場合は canStartNextPhase=true（next-phase カードが活性）
  test("ジョブ終端でない・completed・confirmedCount>=aiTargetCount のとき次段階を開始できる", () => {
    // 成立条件: ジョブが終端状態でない・phaseState が completed・confirmedCount >= aiTargetCount
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "completed",
        aiTargetCount: 5,
        resultSummary: {
          confirmedCount: 5,
          jobDictionaryAppliedCount: 0,
          replacementTargetCount: 0,
          unmatchedCount: 0,
          providerSkipped: false
        },
        actionEnablement: {
          canStart: false,
          canPause: false,
          canResume: false,
          canRetry: false
          // startBlockedReason なし = 終端ジョブでない
        }
      }),
      true
    )

    const nextPhaseCard = viewModel.actionCards.find((c) => c.id === "next-phase")
    // 仕様通り: 成立条件を満たすので次段階を開始できる
    expect(nextPhaseCard?.disabled).toBe(false)
    expect(nextPhaseCard?.blockedReason).toBe("")
    expect(viewModel.nextPhaseStatusLabel).toBe("開始可能")
  })

  // UT-EQV-001: 成立条件の一つが欠ける場合（confirmedCount < aiTargetCount）は開始不可
  test("confirmedCount が aiTargetCount 未満のとき次段階を開始できない（term_phase_incomplete）", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "completed",
        aiTargetCount: 5,
        resultSummary: {
          confirmedCount: 3, // 5 未満
          jobDictionaryAppliedCount: 0,
          replacementTargetCount: 0,
          unmatchedCount: 0,
          providerSkipped: false
        },
        actionEnablement: {
          canStart: false,
          canPause: false,
          canResume: false,
          canRetry: false
        }
      }),
      true
    )

    const nextPhaseCard = viewModel.actionCards.find((c) => c.id === "next-phase")
    expect(nextPhaseCard?.disabled).toBe(true)
    expect(viewModel.nextPhaseStatusLabel).toBe("開始不可")
  })

  // UT-EQV-001: 成立条件の一つが欠ける場合（phaseState が completed でない）は開始不可
  test("phaseState が running のとき次段階を開始できない（term_phase_incomplete）", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running"
      }),
      true
    )

    const nextPhaseCard = viewModel.actionCards.find((c) => c.id === "next-phase")
    expect(nextPhaseCard?.disabled).toBe(true)
    expect(viewModel.nextPhaseStatusLabel).toBe("開始不可")
  })

  // UT-EQV-002: 区別理由 - ジョブが終端状態のとき「ジョブが終端状態のため」の理由を返す
  test("startBlockedReason=terminal_job のとき次段階不可理由にジョブ終端を示す（terminal_job 相当）", () => {
    // 区別理由: ジョブが終端状態 → startBlockedReason で判定
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "completed",
        aiTargetCount: 5,
        resultSummary: {
          confirmedCount: 5,
          jobDictionaryAppliedCount: 0,
          replacementTargetCount: 0,
          unmatchedCount: 0,
          providerSkipped: false
        },
        actionEnablement: {
          canStart: false,
          startBlockedReason: "terminal_job", // ジョブが終端状態の識別子
          canPause: false,
          canResume: false,
          canRetry: false
        }
      }),
      true
    )

    const nextPhaseCard = viewModel.actionCards.find((c) => c.id === "next-phase")
    expect(nextPhaseCard?.disabled).toBe(true)
    // ジョブ終端の区別理由を示す
    expect(nextPhaseCard?.blockedReason).toContain("終端")
    expect(viewModel.nextPhaseBlockedReason).toContain("終端")
  })

  // UT-EQV-002: 区別理由 - phaseState 未完了のとき「単語翻訳段階が未完了」の理由を返す
  test("phaseState が paused のとき次段階不可理由に段階未完了を示す（term_phase_incomplete 相当）", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "paused",
        actionEnablement: {
          canStart: false,
          canPause: false,
          canResume: true,
          canRetry: false
        }
      }),
      true
    )

    const nextPhaseCard = viewModel.actionCards.find((c) => c.id === "next-phase")
    expect(nextPhaseCard?.disabled).toBe(true)
    expect(nextPhaseCard?.blockedReason).toContain("未完了")
    expect(viewModel.nextPhaseBlockedReason).toContain("未完了")
  })
})

// UT-EQV-003〜006: 操作可否の等価性テスト（term-translation-phase-REQ-008 各操作の区別理由）
// 期待値は detail-spec の区別理由から固定する
describe("TermTranslationPhasePresenter - 操作可否の等価性（UT-EQV-003〜006）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildSummary(
    overrides: Partial<TermTranslationPhaseSummaryResponse>
  ): TermTranslationPhaseSummaryResponse {
    return {
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
        credentialRef: "cred",
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      },
      actionEnablement: {
        canStart: true,
        canPause: false,
        canResume: false,
        canRetry: false
      },
      ...overrides
    }
  }

  function buildState(
    summaryOverrides: Partial<TermTranslationPhaseSummaryResponse> = {}
  ) {
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: buildSummary(summaryOverrides),
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // UT-EQV-003: canStart - 終端でない・実行中フェーズ存在しない・Ready 状態・実行設定構成済みのとき開始可
  test("Ready 状態かつ実行設定構成済みのとき開始できる", () => {
    // 成立条件: 終端でない・activePhaseExists なし・isIdleReady=true・configured=true
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "ready",
        execution: {
          credentialRef: "cred",
          provider: "openai-compatible",
          model: "gpt-4.1-mini",
          executionMode: "batch"
        },
        actionEnablement: {
          canStart: false, // presenter が導出するため gateway 値は参考にしない
          canPause: false,
          canResume: false,
          canRetry: false
        }
      }),
      true
    )

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    // presenter が事実状態から導出した canStart が true になる
    expect(startCard?.disabled).toBe(false)
  })

  // UT-EQV-003: canStart 不可理由 - ジョブが終端状態
  test("terminal_job のとき開始不可でジョブ終端理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "ready",
        actionEnablement: {
          canStart: false,
          startBlockedReason: "terminal_job",
          canPause: false,
          canResume: false,
          canRetry: false
        }
      }),
      true
    )

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    expect(startCard?.disabled).toBe(true)
    expect(startCard?.blockedReason).toContain("終端")
  })

  // UT-EQV-003: canStart 不可理由 - 実行中フェーズが存在する（active_phase_exists 相当）
  test("phaseState=running のとき開始不可で実行中理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running"
      }),
      true
    )

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    expect(startCard?.disabled).toBe(true)
    expect(startCard?.blockedReason).toContain("実行中")
  })

  // UT-EQV-003: canStart 不可理由 - 実行設定が未構成
  test("実行設定が未構成のとき開始不可で未構成理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "ready",
        execution: {
          credentialRef: "",
          provider: "", // 未構成
          model: "",
          executionMode: ""
        }
      }),
      true
    )

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    expect(startCard?.disabled).toBe(true)
    expect(startCard?.blockedReason).toContain("未構成")
  })

  // UT-EQV-004: canPause - 実行中のとき一時停止可
  test("phaseState=running のとき一時停止できる", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "running" }),
      true
    )

    const pauseCard = viewModel.actionCards.find((c) => c.id === "pause")
    expect(pauseCard?.disabled).toBe(false)
  })

  // UT-EQV-004: canPause 不可理由 - 終端ジョブ
  test("terminal_job のとき一時停止不可でジョブ終端理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
        actionEnablement: {
          canStart: false,
          startBlockedReason: "terminal_job",
          canPause: false,
          canResume: false,
          canRetry: false
        }
      }),
      true
    )

    const pauseCard = viewModel.actionCards.find((c) => c.id === "pause")
    expect(pauseCard?.disabled).toBe(true)
    expect(pauseCard?.blockedReason).toContain("終端")
  })

  // UT-EQV-004: canPause 不可理由 - フェーズが実行中でない（phase_not_running 相当）
  test("phaseState=paused のとき一時停止不可でフェーズ実行中でない理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "paused" }),
      true
    )

    const pauseCard = viewModel.actionCards.find((c) => c.id === "pause")
    expect(pauseCard?.disabled).toBe(true)
    expect(pauseCard?.blockedReason).toContain("実行中ではありません")
  })

  // UT-EQV-005: canResume - paused のとき再開可
  test("phaseState=paused のとき再開できる", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "paused" }),
      true
    )

    const resumeCard = viewModel.actionCards.find((c) => c.id === "resume")
    expect(resumeCard?.disabled).toBe(false)
  })

  // UT-EQV-005: canResume - recoverable_fail のとき再開可
  test("phaseState=recoverable_failed のとき再開できる", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "recoverable_failed" }),
      true
    )

    const resumeCard = viewModel.actionCards.find((c) => c.id === "resume")
    expect(resumeCard?.disabled).toBe(false)
  })

  // UT-EQV-005: canResume 不可理由 - フェーズが再開可能でない（phase_not_resumable 相当）
  test("phaseState=running のとき再開不可でフェーズ再開可能でない理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "running" }),
      true
    )

    const resumeCard = viewModel.actionCards.find((c) => c.id === "resume")
    expect(resumeCard?.disabled).toBe(true)
    expect(resumeCard?.blockedReason).toContain("再開可能な状態ではありません")
  })

  // UT-EQV-006: canRetry - recoverable_fail のとき再試行可
  test("phaseState=recoverable_failed のとき再試行できる", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "recoverable_failed" }),
      true
    )

    const retryCard = viewModel.actionCards.find((c) => c.id === "retry")
    expect(retryCard?.disabled).toBe(false)
  })

  // UT-EQV-006: canRetry 不可理由 - フェーズが再試行可能でない（phase_not_retryable 相当）
  test("phaseState=paused のとき再試行不可でフェーズ再試行可能でない理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({ phaseState: "paused" }),
      true
    )

    const retryCard = viewModel.actionCards.find((c) => c.id === "retry")
    expect(retryCard?.disabled).toBe(true)
    expect(retryCard?.blockedReason).toContain("再試行可能な状態ではありません")
  })
})
