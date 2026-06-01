import { describe, expect, test } from "vitest"

import type {
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

import { BodyTranslationPhasePresenter } from "./body-translation-phase.presenter"

interface TestScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: null
  errorMessage: string
  pendingAction:
    | "start"
    | "pause"
    | "resume"
    | "retry"
    | "cancel"
    | "check-output-readiness"
    | null
  hasLoaded: boolean
  initialFetchDone: boolean
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
      requestSummary: {
        providerTargetCount: 8,
        exactDictionaryExclusionCount: 2,
        partialDictionaryConstraintCount: 3
      },
      actionEnablement: {
        canStart: false,
        canPause: true,
        canResume: false,
        canRetry: false,
        canCancel: true
      },
      outputReadiness: {
        completedFieldCount: 3,
        statusConsistent: true,
        outputCount: 0
      }
    },
    outputReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    initialFetchDone: true,
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

  test("進行詳細は AI 送信対象件数を分母として表示する", () => {
    const presenter = new BodyTranslationPhasePresenter()
    const base = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...base.summary!,
          progress: {
            ...base.summary!.progress,
            processedCount: 4,
            totalCount: 10,
            targetCount: 10,
            translatedCount: 3,
            skippedCount: 1
          },
          requestSummary: {
            providerTargetCount: 8,
            exactDictionaryExclusionCount: 2,
            partialDictionaryConstraintCount: 3
          }
        }
      }),
      true
    )

    expect(viewModel.progressDetail).toBe(
      "4 / 8 件 / AI 送信対象 8 件 / 成功 3 件 / スキップ 1 件 / AI 処理中"
    )
    expect(viewModel.targetCountLabel).toBe("10 件")
    expect(viewModel.providerTargetCountLabel).toBe("8 件")
  })

  test("current phase key は画面表示名へ変換する", () => {
    const presenter = new BodyTranslationPhasePresenter()

    const viewModel = presenter.toViewModel(createState(), true)

    expect(viewModel.currentPhaseLabel).toBe("本文翻訳")
  })

  test("段階完了かつ statusConsistent のとき output readiness は ready になる", () => {
    const presenter = new BodyTranslationPhasePresenter()
    const base = createState()

    const viewModel = presenter.toViewModel(
      createState({
        summary: {
          ...base.summary!,
          phaseState: "completed",
          outputReadiness: {
            completedFieldCount: 10,
            statusConsistent: true,
            outputCount: 10
          }
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

// UT-EQV-009/010: body の成果物出力確認可否と操作可否の等価性テスト
// 期待値は detail-spec body-translation-phase-REQ-007 成立条件・区別理由から固定する
// 段階要約の事実（completedFieldCount/statusConsistent/outputCount）から application 層が導出する
describe("BodyTranslationPhasePresenter - 成果物出力確認可否・操作可否の等価性（UT-EQV-009/010）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  function buildBaseSummary(
    overrides: Partial<BodyTranslationPhaseSummaryResponse> = {}
  ): BodyTranslationPhaseSummaryResponse {
    return {
      jobId: 9,
      currentPhase: "body_translation",
      phaseState: "completed",
      phaseRunId: 19,
      progress: {
        percent: 100,
        processedCount: 10,
        totalCount: 10,
        targetCount: 10,
        translatedCount: 10,
        skippedCount: 0,
        currentStep: "completed"
      },
      inputSummary: {
        targetCount: 10,
        dictionaryDigest: "sha256:dict",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:meta",
        promptDigest: "sha256:prompt"
      },
      execution: {
        credentialRef: "cred",
        provider: "fake",
        model: "model",
        executionMode: "batch",
        requestUnitCount: 10,
        outputCount: 10
      },
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: false,
        canCancel: false
      },
      outputReadiness: {
        completedFieldCount: 10,
        statusConsistent: true,
        outputCount: 10
      },
      ...overrides
    }
  }

  function buildState(
    summaryOverrides: Partial<BodyTranslationPhaseSummaryResponse> = {}
  ) {
    return {
      jobId: 9,
      phase: "ready" as const,
      summary: buildBaseSummary(summaryOverrides),
      outputReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // UT-EQV-009: 成立条件 - 段階完了・statusConsistent=true のとき成果物出力確認可（ready=true）
  test("phaseState=completed かつ statusConsistent=true のとき成果物出力確認が ready になる", () => {
    // 成立条件: phaseState が completed・statusConsistent=true・completedFieldCount >= 0
    const viewModel = presenter.toViewModel(buildState(), true)

    expect(viewModel.outputReadinessLabel).toBe("ready")
    expect(viewModel.outputReadinessBlockedReason).toBe("")
    // 事実状態が正しく表示される
    expect(viewModel.outputReadinessCompletedFieldCountLabel).toBe("10 件")
  })

  // UT-EQV-009: 段階が完了していないとき成果物出力確認不可（ready=false）
  test("phaseState=running のとき成果物出力確認が not ready になる", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
        outputReadiness: {
          completedFieldCount: 5,
          statusConsistent: true,
          outputCount: 5
        }
      }),
      true
    )

    // 段階未完了なので ready でない
    expect(viewModel.outputReadinessLabel).toBe("not ready")
    expect(viewModel.outputReadinessBlockedReason).toContain("完了していません")
  })

  // UT-EQV-009: statusConsistent=false のとき成果物出力確認不可
  test("phaseState=completed でも statusConsistent=false のとき成果物出力確認が not ready になる", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "completed",
        outputReadiness: {
          completedFieldCount: 10,
          statusConsistent: false, // 整合性なし
          outputCount: 8
        }
      }),
      true
    )

    expect(viewModel.outputReadinessLabel).toBe("not ready")
    expect(viewModel.outputReadinessBlockedReason).toContain("整合")
  })

  // UT-EQV-009: 区別理由 - 専用取得廃止後も段階要約の事実から同じ判定が導ける（等価性確認）
  test("outputReadiness の事実状態（completedFieldCount/statusConsistent/outputCount）が viewModel に正しく反映される", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        outputReadiness: {
          completedFieldCount: 7,
          statusConsistent: true,
          outputCount: 7
        }
      }),
      true
    )

    // 段階要約から導出した事実状態が viewModel に正しく伝わる
    expect(viewModel.outputReadinessCompletedFieldCountLabel).toBe("7 件")
    expect(viewModel.outputReadinessStatusLabel).toBe("整合")
  })

  // UT-EQV-010: 操作可否 - running のとき canPause=true・canStart=false・canCheckOutputReadiness=false
  test("phaseState=running のとき一時停止可・開始不可・出力確認不可", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
        outputReadiness: {
          completedFieldCount: 5,
          statusConsistent: true,
          outputCount: 5
        }
      }),
      true
    )

    // 操作可否は事実状態から導出する
    expect(viewModel.screenActionEnablement.canPause).toBe(true)
    expect(viewModel.screenActionEnablement.canStart).toBe(false)
    // 段階未完了なので出力確認不可
    expect(viewModel.screenActionEnablement.canCheckOutputReadiness).toBe(false)
  })

  // UT-EQV-010: 操作可否 - completed かつ statusConsistent=true のとき canCheckOutputReadiness=true
  test("phaseState=completed かつ statusConsistent=true のとき出力確認可", () => {
    const viewModel = presenter.toViewModel(buildState(), true)

    expect(viewModel.screenActionEnablement.canCheckOutputReadiness).toBe(true)
    const checkCard = viewModel.actionCards.find((c) => c.id === "check-output-readiness")
    expect(checkCard?.disabled).toBe(false)
  })

  // UT-EQV-010: 操作可否 - recoverable_failed のとき canRetry=true・canCancel=true
  test("phaseState=recoverable_failed のとき再試行可かつキャンセル可", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "recoverable_failed",
        outputReadiness: {
          completedFieldCount: 5,
          statusConsistent: false,
          outputCount: 0
        }
      }),
      true
    )

    expect(viewModel.screenActionEnablement.canRetry).toBe(true)
    expect(viewModel.screenActionEnablement.canCancel).toBe(true)
    expect(viewModel.screenActionEnablement.canCheckOutputReadiness).toBe(false)
  })

  // UT-EQV-010: terminal_job のとき全操作不可
  test("startBlockedReason=terminal_job のとき全操作不可", () => {
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
        actionEnablement: {
          canStart: false,
          startBlockedReason: "terminal_job",
          canPause: false,
          canResume: false,
          canRetry: false,
          canCancel: false
        }
      }),
      true
    )

    expect(viewModel.screenActionEnablement.canStart).toBe(false)
    expect(viewModel.screenActionEnablement.canPause).toBe(false)
    expect(viewModel.screenActionEnablement.canResume).toBe(false)
    expect(viewModel.screenActionEnablement.canRetry).toBe(false)
    expect(viewModel.screenActionEnablement.canCancel).toBe(false)
  })
})

// U-PRES-002/003/004: execution field 不在時の presenter 判定（本文翻訳）
// 根拠: body-translation-phase-REQ-002「execution field 不在で isExecutionConfigured=false」
describe("BodyTranslationPhasePresenter - execution field 不在判定（U-PRES-002/003/004）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  function buildSummaryWithoutExecution(): BodyTranslationPhaseSummaryResponse {
    return {
      jobId: 9,
      currentPhase: "body_translation",
      phaseState: "ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 10,
        targetCount: 10,
        translatedCount: 0,
        skippedCount: 0,
        currentStep: "ready"
      },
      inputSummary: {
        targetCount: 10,
        dictionaryDigest: "sha256:dict",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:meta",
        promptDigest: "sha256:prompt"
      },
      // execution field を意図的に含めない（不在で「未設定」を表す仕様）
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: false,
        canCancel: false
      },
      outputReadiness: {
        completedFieldCount: 0,
        statusConsistent: false,
        outputCount: 0
      }
    }
  }

  function buildStateWithoutExecution() {
    return {
      jobId: 9,
      phase: "ready" as const,
      summary: buildSummaryWithoutExecution(),
      outputReadiness: null,
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

    expect(viewModel.isExecutionConfigured).toBe(false)
  })

  // U-PRES-003: execution field が不在の場合、modelLabel が「設定未完了」を返す（空文字 "" を返さない）
  test("execution field が不在のとき modelLabel が設定未完了を返す", () => {
    // 空文字フォールバックを防止し「設定未完了」相当の表示語を返すことを証明する。
    const viewModel = presenter.toViewModel(buildStateWithoutExecution(), true)

    expect(viewModel.modelLabel).not.toBe("")
    expect(viewModel.modelLabel).toBe("設定未完了")
  })

  // U-PRES-002（派生）: execution field が存在する場合と不在の場合を presenter が独立して判定できる
  test("execution field が存在しモデル値が入力済みのとき isExecutionConfigured が true を返す", () => {
    // 2 つの不在状態の混同を防ぐ。設定済み状態と未設定状態を独立して判定できることを証明する。
    const summaryWithExecution: BodyTranslationPhaseSummaryResponse = {
      ...buildSummaryWithoutExecution(),
      execution: {
        credentialRef: "cred",
        provider: "lm_studio",
        model: "local-model",
        executionMode: "sync",
        requestUnitCount: 5,
        outputCount: 0
      }
    }
    const viewModel = presenter.toViewModel(
      {
        jobId: 9,
        phase: "ready" as const,
        summary: summaryWithExecution,
        outputReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true
      },
      true
    )

    // execution field が存在し値が入力済み → isExecutionConfigured=true
    expect(viewModel.isExecutionConfigured).toBe(true)
    expect(viewModel.modelLabel).toBe("local-model")
  })
})
