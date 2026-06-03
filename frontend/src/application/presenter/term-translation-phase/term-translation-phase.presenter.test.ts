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

const defaultProjection = {
  phaseLifecycle: "ready",
  jobLifecycle: "running",
  errorKind: "none",
  aiSettingsConfigured: true,
  aiTargetCount: 5,
  confirmedCount: 0
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
      projection: { ...defaultProjection }
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
      projection: { ...defaultProjection },
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
        projection: {
          ...defaultProjection,
          phaseLifecycle: "completed",
          jobLifecycle: "running", // 終端でない
          confirmedCount: 5, // aiTargetCount 以上
          aiTargetCount: 5
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
        projection: {
          ...defaultProjection,
          phaseLifecycle: "completed",
          jobLifecycle: "running"
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
  test("jobLifecycle=completed のとき次段階不可理由にジョブ終端を示す（terminal_job 相当）", () => {
    // 区別理由: ジョブが終端状態 → projection.jobLifecycle で判定
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
        projection: {
          ...defaultProjection,
          phaseLifecycle: "completed",
          jobLifecycle: "completed" // ジョブが終端状態
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
        projection: {
          ...defaultProjection,
          phaseLifecycle: "paused"
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

  // phaseState を projection.phaseLifecycle に整合させるヘルパー
  // 新実装は projection.phaseLifecycle で判定するため、phaseState と phaseLifecycle を同一値にする
  function buildSummary(
    overrides: Partial<TermTranslationPhaseSummaryResponse>
  ): TermTranslationPhaseSummaryResponse {
    const phaseState = overrides.phaseState ?? "ready"
    const baseProjection = { ...defaultProjection, phaseLifecycle: phaseState }
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
      projection: baseProjection,
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
        projection: {
          ...defaultProjection,
          phaseLifecycle: "ready",
          jobLifecycle: "running" // presenter が導出するため、終端でないこと
        }
      }),
      true
    )

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    // presenter が事実状態から導出した canStart が true になる
    expect(startCard?.disabled).toBe(false)
  })

  // UT-EQV-003: canStart 不可理由 - ジョブが終端状態
  test("jobLifecycle=completed のとき開始不可でジョブ終端理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "ready",
        projection: {
          ...defaultProjection,
          jobLifecycle: "completed" // 終端
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
  // projection.aiSettingsConfigured で判定するため、projection を false に上書きする
  test("実行設定が未構成のとき開始不可で未構成理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "ready",
        execution: {
          credentialRef: "",
          provider: "", // 未構成
          model: "",
          executionMode: ""
        },
        projection: {
          ...defaultProjection,
          phaseLifecycle: "ready",
          aiSettingsConfigured: false // 未構成を projection から導出
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
  test("jobLifecycle=completed のとき一時停止不可でジョブ終端理由を返す", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
        projection: {
          ...defaultProjection,
          phaseLifecycle: "running",
          jobLifecycle: "completed" // 終端
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

// U-PRES-PROVIDER-001/002: provider 一覧連動テスト（単語翻訳）
describe("TermTranslationPhasePresenter - provider 一覧連動（U-PRES-PROVIDER-001/002）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(): TestScreenState {
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
        projection: {
          ...defaultProjection,
          aiSettingsConfigured: false
        }
      },
      nextPhaseReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-PROVIDER-001: provider 一覧が渡ると providerOptions に並ぶ
  test("provider 一覧が渡ると providerOptions に並ぶ", () => {
    const availableProviders = [
      { value: "openai", label: "OpenAI" },
      { value: "anthropic", label: "Anthropic" }
    ]

    const viewModel = presenter.toViewModel(buildState(), true, availableProviders)

    expect(viewModel.providerOptions).toContainEqual({ value: "openai", label: "OpenAI" })
    expect(viewModel.providerOptions).toContainEqual({ value: "anthropic", label: "Anthropic" })
  })

  // U-PRES-PROVIDER-002: 未選択時は先頭が「選んでください」
  test("現在の provider が未設定のとき先頭が選んでください", () => {
    const availableProviders = [
      { value: "openai", label: "OpenAI" },
      { value: "anthropic", label: "Anthropic" }
    ]

    // summary に aiSettings も execution も provider が設定されていない状態
    const viewModel = presenter.toViewModel(buildState(), true, availableProviders)

    expect(viewModel.providerOptions[0]).toEqual({ value: "", label: "選んでください" })
  })

  // U-PRES-PROVIDER-003: provider 一覧が空のときは現在のラベルをフォールバックで返す
  test("provider 一覧が空のとき設定未完了ラベルをフォールバックで返す", () => {
    const viewModel = presenter.toViewModel(buildState(), true, [])

    expect(viewModel.providerOptions).toHaveLength(1)
    expect(viewModel.providerOptions[0]?.label).toBe("設定未完了")
  })
})

// U-PRES-MODEL-001: buildModelOptions の placeholder 重複解消（単語翻訳）
// 根拠: fix-phase-ai-settings-pill-update-after-model-select の secondary 問題修正
// AIModelSelectionCard が固定 placeholder を持つため、buildModelOptions が重複 placeholder を返さないことを証明する
describe("TermTranslationPhasePresenter - buildModelOptions placeholder 重複解消（U-PRES-MODEL-001）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(): TestScreenState {
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
        projection: {
          ...defaultProjection,
          aiSettingsConfigured: false
        }
      },
      nextPhaseReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-MODEL-001: availableModels に要素がある時、modelOptions に placeholder が含まれない
  test("availableModels に要素がある時 modelOptions に空 value の placeholder が含まれない", () => {
    // AIModelSelectionCard が固定 placeholder を描画するため、buildModelOptions は重複 placeholder を返さない。
    const availableModels = [{ value: "fake-model", label: "fake-model" }]

    const viewModel = presenter.toViewModel(buildState(), true, [], availableModels)

    // availableModels を渡した時 placeholder が含まれないことを証明する
    const hasPlaceholder = viewModel.modelOptions.some((opt) => opt.value === "")
    expect(hasPlaceholder).toBe(false)
  })

  // U-PRES-MODEL-001: availableModels に要素がある時、modelOptions は availableModels の内容を返す
  test("availableModels に要素がある時 modelOptions が availableModels をそのまま返す", () => {
    // placeholder を挿入せず availableModels の内容だけを返すことを証明する。
    const availableModels = [
      { value: "model-a", label: "Model A" },
      { value: "model-b", label: "Model B" }
    ]

    const viewModel = presenter.toViewModel(buildState(), true, [], availableModels)

    expect(viewModel.modelOptions).toHaveLength(2)
    expect(viewModel.modelOptions).toContainEqual({ value: "model-a", label: "Model A" })
    expect(viewModel.modelOptions).toContainEqual({ value: "model-b", label: "Model B" })
  })
})

// RAEF-UNIT-001〜009: deriveTermActionEnablement および deriveCanStartNextPhase の境界値証明
// 根拠: design-diff.md H-1〜H-5
describe("deriveTermActionEnablement — H-1 start 境界値（RAEF-UNIT-001〜004）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(projectionOverrides: Partial<TermTranslationPhaseSummaryResponse["projection"]> = {}) {
    const baseProjection = {
      phaseLifecycle: "ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 1,
      confirmedCount: 0,
      ...projectionOverrides
    }
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: {
        jobId: 7,
        currentPhase: "term_translation",
        phaseState: "ready",
        progress: {
          percent: 0,
          processedCount: 0,
          totalCount: 10,
          aiTargetCount: 1,
          currentStep: "ready"
        },
        totalTermCount: 10,
        dictionaryHitCount: 0,
        aiTargetCount: 1,
        projection: baseProjection
      },
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // RAEF-UNIT-001: jobLifecycle が TERMINAL_JOB のとき canStart=false かつ startBlockedReason にジョブ終端文を返す
  test("jobLifecycle が completed のとき開始不可でジョブ終端理由を返す（H-1 terminal 境界）", () => {
    // H-1: jobLifecycle ∈ TERMINAL_JOB → canStart=false、startBlockedReason=「ジョブが終端状態のため開始できません。」
    const vm = presenter.toViewModel(buildState({ jobLifecycle: "completed" }), true)
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のため開始できません。")
  })

  // RAEF-UNIT-002: aiSettingsConfigured=false のとき canStart=false かつ startBlockedReason に未構成文を返す
  test("aiSettingsConfigured=false のとき開始不可で実行設定未構成理由を返す（H-1 aiSettingsConfigured=false 境界）", () => {
    // H-1: aiSettingsConfigured ≠ true → canStart=false、startBlockedReason=「実行設定が未構成のため開始できません。」
    const vm = presenter.toViewModel(buildState({ aiSettingsConfigured: false }), true)
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("実行設定が未構成のため開始できません。")
  })

  // RAEF-UNIT-003: aiTargetCount=0 のとき canStart=false かつ startBlockedReason に 0件文を返す
  test("aiTargetCount=0 のとき開始不可で処理対象 0 件理由を返す（H-1 aiTargetCount≤0 境界値）", () => {
    // H-1: aiTargetCount ≤ 0 → canStart=false、startBlockedReason=「処理対象が 0 件のため開始できません。」
    const vm = presenter.toViewModel(buildState({ aiTargetCount: 0 }), true)
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("処理対象が 0 件のため開始できません。")
  })

  // RAEF-UNIT-004: 全条件満足（¬terminal ∧ ¬running ∧ idleReady ∧ aiSettingsConfigured ∧ aiTargetCount=1）で canStart=true
  test("全条件満足のとき開始可かつ startBlockedReason が空文字になる（H-1 正常パス、aiTargetCount 最小値 1）", () => {
    // H-1: 全有効化条件を満たす最小境界値（aiTargetCount=1）で canStart=true、startBlockedReason=""
    const vm = presenter.toViewModel(
      buildState({ phaseLifecycle: "ready", jobLifecycle: "running", aiSettingsConfigured: true, aiTargetCount: 1 }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })
})

describe("deriveTermActionEnablement — H-2 pause 境界値（RAEF-UNIT-005）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(projectionOverrides: Partial<TermTranslationPhaseSummaryResponse["projection"]> = {}) {
    const baseProjection = {
      phaseLifecycle: "running",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 5,
      confirmedCount: 0,
      ...projectionOverrides
    }
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: {
        jobId: 7,
        currentPhase: "term_translation",
        phaseState: "running",
        progress: { percent: 0, processedCount: 0, totalCount: 10, aiTargetCount: 5, currentStep: "running" },
        totalTermCount: 10,
        dictionaryHitCount: 0,
        aiTargetCount: 5,
        projection: baseProjection
      },
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // RAEF-UNIT-005: RUNNING_PHASE のとき canPause=true、IDLE_READY_PHASE のとき canPause=false かつ pauseBlockedReason を返す
  test("phaseLifecycle が running のとき中断可になる（H-2 running 境界）", () => {
    // H-2: phaseLifecycle ∈ RUNNING_PHASE → canPause=true
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "running" }), true)
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(false)
  })

  test("phaseLifecycle が idle_ready のとき中断不可でフェーズ実行中でない理由を返す（H-2 非 running 境界）", () => {
    // H-2: phaseLifecycle ∉ RUNNING_PHASE → canPause=false、pauseBlockedReason=「フェーズが実行中ではありません。」
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "idle_ready" }), true)
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが実行中ではありません。")
  })
})

describe("deriveTermActionEnablement — H-3 resume 境界値（RAEF-UNIT-006）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(projectionOverrides: Partial<TermTranslationPhaseSummaryResponse["projection"]> = {}) {
    const baseProjection = {
      phaseLifecycle: "idle_ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 5,
      confirmedCount: 0,
      ...projectionOverrides
    }
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: {
        jobId: 7,
        currentPhase: "term_translation",
        phaseState: "idle_ready",
        progress: { percent: 0, processedCount: 0, totalCount: 10, aiTargetCount: 5, currentStep: "idle_ready" },
        totalTermCount: 10,
        dictionaryHitCount: 0,
        aiTargetCount: 5,
        projection: baseProjection
      },
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // RAEF-UNIT-006: paused のとき canResume=true、errorKind=recoverable のとき canResume=true、両条件非該当のとき canResume=false
  test("phaseLifecycle が paused のとき再開可になる（H-3 PAUSED_PHASE 境界）", () => {
    // H-3: phaseLifecycle ∈ PAUSED_PHASE → canResume=true
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "paused" }), true)
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
  })

  test("errorKind=recoverable のとき再開可になる（H-3 errorKind=recoverable 境界）", () => {
    // H-3: errorKind = recoverable → canResume=true（phaseLifecycle が RECOVERABLE_FAILED_PHASE 外でも）
    const vm = presenter.toViewModel(
      buildState({ phaseLifecycle: "idle_ready", errorKind: "recoverable" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
  })

  test("paused でも recoverableFailed でもないとき再開不可でフェーズ再開可能でない理由を返す（H-3 両条件非該当境界）", () => {
    // H-3: phaseLifecycle ∉ PAUSED_PHASE ∧ phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE ∧ errorKind ≠ recoverable → canResume=false
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "running", errorKind: "none" }), true)
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが再開可能な状態ではありません。")
  })
})

describe("deriveTermActionEnablement — H-4 retry 境界値（RAEF-UNIT-007）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(projectionOverrides: Partial<TermTranslationPhaseSummaryResponse["projection"]> = {}) {
    const baseProjection = {
      phaseLifecycle: "idle_ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 5,
      confirmedCount: 0,
      ...projectionOverrides
    }
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: {
        jobId: 7,
        currentPhase: "term_translation",
        phaseState: "idle_ready",
        progress: { percent: 0, processedCount: 0, totalCount: 10, aiTargetCount: 5, currentStep: "idle_ready" },
        totalTermCount: 10,
        dictionaryHitCount: 0,
        aiTargetCount: 5,
        projection: baseProjection
      },
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // RAEF-UNIT-007: RECOVERABLE_FAILED_PHASE のとき canRetry=true、非該当のとき canRetry=false かつ retryBlockedReason を返す
  test("phaseLifecycle が recoverable_failed のとき再試行可になる（H-4 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-4: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE → canRetry=true
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "recoverable_failed" }), true)
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(false)
  })

  test("phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE かつ errorKind ≠ recoverable のとき再試行不可で理由を返す（H-4 非 recoverableFailed 境界）", () => {
    // H-4: phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE ∧ errorKind ≠ recoverable → canRetry=false、retryBlockedReason 返す
    const vm = presenter.toViewModel(buildState({ phaseLifecycle: "paused", errorKind: "none" }), true)
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが再試行可能な状態ではありません。")
  })
})

describe("deriveCanStartNextPhase — H-5 境界値（RAEF-UNIT-008〜009）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildState(projectionOverrides: Partial<TermTranslationPhaseSummaryResponse["projection"]> = {}) {
    const baseProjection = {
      phaseLifecycle: "completed",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 5,
      confirmedCount: 5,
      ...projectionOverrides
    }
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: {
        jobId: 7,
        currentPhase: "term_translation",
        phaseState: "completed",
        progress: { percent: 100, processedCount: 5, totalCount: 10, aiTargetCount: 5, currentStep: "completed" },
        totalTermCount: 10,
        dictionaryHitCount: 0,
        aiTargetCount: 5,
        projection: baseProjection
      },
      nextPhaseReadiness: null as TermTranslationNextPhaseReadinessResponse | null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // RAEF-UNIT-008: COMPLETED_PHASE かつ confirmedCount >= aiTargetCount のとき canStartNextPhase=true
  test("phaseLifecycle が completed かつ confirmedCount >= aiTargetCount のとき次段階開始可（H-5 正常パス）", () => {
    // H-5: phaseLifecycle ∈ COMPLETED_PHASE ∧ confirmedCount ≥ aiTargetCount → canStartNextPhase=true
    const vm = presenter.toViewModel(
      buildState({ phaseLifecycle: "completed", confirmedCount: 5, aiTargetCount: 5 }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "next-phase")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-009: confirmedCount < aiTargetCount のとき canStartNextPhase=false かつ aiTargetSatisfied=false 理由を返す
  test("confirmedCount が aiTargetCount-1 のとき次段階開始不可で確定件数不足理由を返す（H-5 aiTargetSatisfied=false 境界値）", () => {
    // H-5: confirmedCount < aiTargetCount → canStartNextPhase=false、BlockedReason=「確定件数が AI 対象件数に達していないため次段階を開始できません。」
    // confirmedCount=4、aiTargetCount=5 の境界値（aiTargetCount-1）を使う
    const vm = presenter.toViewModel(
      buildState({ phaseLifecycle: "completed", confirmedCount: 4, aiTargetCount: 5 }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "next-phase")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("確定件数が AI 対象件数に達していないため次段階を開始できません。")
  })
})

// U-PRES-002/003/004: execution field 不在時の presenter 判定（単語翻訳）
// 根拠: term-translation-phase-REQ-002「execution field 不在で isExecutionConfigured=false」
describe("TermTranslationPhasePresenter - execution field 不在判定（U-PRES-002/003/004）", () => {
  const presenter = new TermTranslationPhasePresenter()

  function buildSummaryWithoutExecution(): TermTranslationPhaseSummaryResponse {
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
      // execution field を意図的に含めない（不在で「未設定」を表す仕様）
      projection: {
        ...defaultProjection,
        aiSettingsConfigured: false
      }
    }
  }

  function buildStateWithoutExecution() {
    return {
      jobId: 7,
      phase: "ready" as const,
      summary: buildSummaryWithoutExecution(),
      nextPhaseReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-002: execution field が不在の場合、isExecutionConfigured が false を返す
  test("execution field が不在のとき isExecutionConfigured が false を返す", () => {
    // 空文字ではなく field 不在で未設定を判定できることを証明する。
    const viewModel = presenter.toViewModel(buildStateWithoutExecution(), true)

    // field 不在 → isExecutionConfigured=false
    expect(viewModel.isExecutionConfigured).toBe(false)
  })

  // U-PRES-003: execution field が不在の場合、modelLabel が「設定未完了」を返す（空文字 "" を返さない）
  test("execution field が不在のとき modelLabel が設定未完了を返す", () => {
    // 空文字フォールバックを防止し「設定未完了」相当の表示語を返すことを証明する（確定原因 2 の回帰防止）。
    const viewModel = presenter.toViewModel(buildStateWithoutExecution(), true)

    // execution 不在 → modelLabel は空文字 "" ではなく表示語を返す
    expect(viewModel.modelLabel).not.toBe("")
    expect(viewModel.modelLabel).toBe("設定未完了")
  })

  // U-PRES-003: execution field が不在の場合、providerLabel が「設定未完了」を返す
  test("execution field が不在のとき providerLabel が設定未完了を返す", () => {
    // provider も model と同様に空文字フォールバックを防止する。
    const viewModel = presenter.toViewModel(buildStateWithoutExecution(), true)

    expect(viewModel.providerLabel).not.toBe("")
    expect(viewModel.providerLabel).toBe("設定未完了")
  })

  // U-PRES-004: execution field が不在の場合、start card が disabled になる
  test("execution field が不在のとき start card が disabled になる", () => {
    // AI 設定不足での開始ボタン無効化を証明する（確定原因 3 の回帰防止）。
    const viewModel = presenter.toViewModel(buildStateWithoutExecution(), true)

    const startCard = viewModel.actionCards.find((c) => c.id === "start")
    expect(startCard?.disabled).toBe(true)
  })

  // U-PRES-005: execution field が存在する場合、isExecutionConfigured が true を返す
  test("execution field が存在しモデル値が入力済みのとき isExecutionConfigured が true を返す", () => {
    // 設定済み状態と未設定状態を presenter が独立して判定できることを証明する。
    const summaryWithExecution: TermTranslationPhaseSummaryResponse = {
      ...buildSummaryWithoutExecution(),
      execution: {
        credentialRef: "cred-main",
        provider: "openai-compatible",
        model: "gpt-4.1-mini",
        executionMode: "batch"
      }
    }
    const viewModel = presenter.toViewModel(
      {
        jobId: 7,
        phase: "ready" as const,
        summary: summaryWithExecution,
        nextPhaseReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true
      },
      true
    )

    // execution field が存在し値が入力済み → isExecutionConfigured=true
    expect(viewModel.isExecutionConfigured).toBe(true)
    expect(viewModel.modelLabel).toBe("gpt-4.1-mini")
  })
})
