import { describe, expect, test } from "vitest"

import type {
  BodyTranslationPhaseProjection,
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
          projection: {
            ...base.summary!.projection,
            phaseLifecycle: "pending"
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
      projection: {
        phaseLifecycle: "completed",
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
    // 新実装は projection を唯一の入力とするため、projection.phaseLifecycle も running に揃える
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "running",
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
          completedFieldCount: 5,
          statusConsistent: true,
          outputCount: 5
        }
      }),
      true
    )

    // 操作可否は projection から導出する
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
    // 新実装は projection を唯一の入力とするため、projection.phaseLifecycle も recoverable_failed に揃える
    const viewModel = presenter.toViewModel(
      buildState({
        phaseState: "recoverable_failed",
        projection: {
          phaseLifecycle: "recoverable_failed",
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
        projection: {
          phaseLifecycle: "running",
          jobLifecycle: "completed", // terminal
          errorKind: "none",
          aiSettingsConfigured: true,
          targetCount: 10,
          previousPhaseLifecycle: "completed",
          personaBodyReadiness: {
            bodyReadiness: true,
            snapshotReferenceStatus: "available"
          }
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
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: false,
        targetCount: 10,
        previousPhaseLifecycle: "completed",
        personaBodyReadiness: {
          bodyReadiness: true,
          snapshotReferenceStatus: "available"
        }
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

// U-PRES-MODEL-003: buildBodyModelOptions の placeholder 重複解消（本文翻訳）
// 根拠: fix-phase-ai-settings-pill-update-after-model-select の secondary 問題修正
// AIModelSelectionCard が固定 placeholder を持つため、buildBodyModelOptions が重複 placeholder を返さないことを証明する
describe("BodyTranslationPhasePresenter - buildBodyModelOptions placeholder 重複解消（U-PRES-MODEL-003）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  function buildState(): TestScreenState {
    return {
      jobId: 9,
      phase: "ready",
      summary: {
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
        projection: {
          phaseLifecycle: "ready",
          jobLifecycle: "running",
          errorKind: "none",
          aiSettingsConfigured: false,
          targetCount: 10,
          previousPhaseLifecycle: "completed",
          personaBodyReadiness: {
            bodyReadiness: true,
            snapshotReferenceStatus: "available"
          }
        },
        outputReadiness: {
          completedFieldCount: 0,
          statusConsistent: false,
          outputCount: 0
        }
      },
      outputReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-MODEL-003: availableModels に要素がある時、modelOptions に placeholder が含まれない
  test("availableModels に要素がある時 modelOptions に空 value の placeholder が含まれない", () => {
    // AIModelSelectionCard が固定 placeholder を描画するため、buildBodyModelOptions は重複 placeholder を返さない。
    const availableModels = [{ value: "fake-model", label: "fake-model" }]

    const viewModel = presenter.toViewModel(buildState(), true, [], availableModels)

    // availableModels を渡した時 placeholder が含まれないことを証明する
    const hasPlaceholder = viewModel.modelOptions.some((opt) => opt.value === "")
    expect(hasPlaceholder).toBe(false)
  })

  // U-PRES-MODEL-003: availableModels に要素がある時、modelOptions は availableModels の内容を返す
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

// RAEF-UNIT-015〜018: deriveBodyActionEnablement の境界値証明
// 根拠: design-diff.md H-12

// body 画面の state を projection 差し替えで構築するヘルパー
function buildBodyStateWithProjection(
  projectionOverrides: Partial<BodyTranslationPhaseProjection> = {}
) {
  const baseProjection: BodyTranslationPhaseProjection = {
    phaseLifecycle: "idle_ready",
    jobLifecycle: "running",
    errorKind: "none",
    aiSettingsConfigured: true,
    targetCount: 1,
    previousPhaseLifecycle: "completed",
    personaBodyReadiness: {
      bodyReadiness: true,
      snapshotReferenceStatus: "available"
    },
    ...projectionOverrides
  }
  return {
    jobId: 9,
    phase: "ready" as const,
    summary: {
      jobId: 9,
      currentPhase: "body_translation",
      phaseState: "idle_ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 10,
        targetCount: 1,
        translatedCount: 0,
        skippedCount: 0,
        currentStep: "idle_ready"
      },
      inputSummary: {
        targetCount: 1,
        dictionaryDigest: "sha256:dictionary",
        personaDigest: "sha256:persona",
        metadataDigest: "sha256:metadata",
        promptDigest: "sha256:prompt"
      },
      requestSummary: {
        providerTargetCount: 1,
        exactDictionaryExclusionCount: 0,
        partialDictionaryConstraintCount: 0
      },
      execution: {
        credentialRef: "cred",
        provider: "fake",
        model: "model",
        executionMode: "batch",
        requestUnitCount: 0,
        outputCount: 0
      },
      outputReadiness: {
        completedFieldCount: 0,
        statusConsistent: true,
        outputCount: 0
      },
      projection: baseProjection
    },
    outputReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    initialFetchDone: true
  }
}

describe("deriveBodyActionEnablement — H-12 start 境界値（RAEF-UNIT-015〜018）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  // RAEF-UNIT-015: previousPhaseLifecycle ∉ COMPLETED_PHASE のとき canStart=false かつ ペルソナ生成未完了理由を返す
  test("previousPhaseLifecycle が running のとき開始不可でペルソナ生成段階未完了理由を返す（H-12 previousPhaseCompleted=false 境界）", () => {
    // H-12: previousPhaseLifecycle ∉ COMPLETED_PHASE → canStart=false、startBlockedReason=「ペルソナ生成段階が完了していないため本文翻訳を開始できません。」
    const vm = presenter.toViewModel(buildBodyStateWithProjection({ previousPhaseLifecycle: "running" }), true)
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ペルソナ生成段階が完了していないため本文翻訳を開始できません。")
  })

  // RAEF-UNIT-016: previousPhaseCompleted=true だが personaBodyReadiness.bodyReadiness=false かつ snapshotReferenceStatus≠available のとき canStart=false かつ snapshot 参照未準備理由を返す
  test("previousPhaseLifecycle が completed だが personaBodyReadiness 両条件とも false のとき開始不可でスナップショット未準備理由を返す（H-12 personaBodyReady=false 境界）", () => {
    // H-12: previousPhaseCompleted=true ∧ bodyReadiness=false ∧ snapshotReferenceStatus≠available → canStart=false、startBlockedReason=「ペルソナ snapshot 参照が準備できていないため本文翻訳を開始できません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({
        previousPhaseLifecycle: "completed",
        personaBodyReadiness: { bodyReadiness: false, snapshotReferenceStatus: "missing" }
      }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ペルソナ snapshot 参照が準備できていないため本文翻訳を開始できません。")
  })

  // RAEF-UNIT-017: personaBodyReadiness.bodyReadiness=true のとき personaBodyReady=true として canStart 成立。snapshotReferenceStatus=available のとき同様に canStart 成立
  test("personaBodyReadiness.bodyReadiness=true のとき開始条件が personaBodyReady=true になる（H-12 bodyReadiness 経路）", () => {
    // H-12: personaBodyReadiness.bodyReadiness=true → personaBodyReady=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({
        personaBodyReadiness: { bodyReadiness: true, snapshotReferenceStatus: "" }
      }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    // personaBodyReady=true なので他の条件が満たされれば開始可
    expect(card?.disabled).toBe(false)
  })

  test("snapshotReferenceStatus=available のとき開始条件が personaBodyReady=true になる（H-12 snapshotReferenceStatus 経路）", () => {
    // H-12: personaBodyReadiness.snapshotReferenceStatus=available → personaBodyReady=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({
        personaBodyReadiness: { bodyReadiness: false, snapshotReferenceStatus: "available" }
      }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    // personaBodyReady=true（snapshotReferenceStatus 経路）なので開始可
    expect(card?.disabled).toBe(false)
  })

  // RAEF-UNIT-018: 全条件満足（¬terminal ∧ idleReady ∧ aiSettingsConfigured ∧ targetCount=1 ∧ previousPhaseCompleted ∧ personaBodyReady）のとき canStart=true かつ startBlockedReason が空文字
  test("全条件満足のとき開始可かつ startBlockedReason が空文字になる（H-12 正常パス、targetCount 最小値 1）", () => {
    // H-12: 全有効化条件を満たす最小境界値で canStart=true、startBlockedReason=""
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({
        phaseLifecycle: "idle_ready",
        jobLifecycle: "running",
        aiSettingsConfigured: true,
        targetCount: 1,
        previousPhaseLifecycle: "completed",
        personaBodyReadiness: { bodyReadiness: true, snapshotReferenceStatus: "available" }
      }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })
})

// RAEF-UNIT-EXT-010〜018: deriveBodyActionEnablement の H-13〜H-16 境界値証明
// 根拠: design-diff.md H-13（pause）、H-14（resume）、H-15（retry）、H-16（cancel）

describe("deriveBodyActionEnablement — H-13 pause 境界値（RAEF-UNIT-EXT-010〜012）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  // RAEF-UNIT-EXT-010: phaseLifecycle ∈ RUNNING_PHASE かつ ¬terminal のとき canPause=true
  test("phaseLifecycle が running かつ ¬terminal のとき中断可になる（H-13 RUNNING_PHASE 境界）", () => {
    // H-13: jobLifecycle ∉ TERMINAL_JOB ∧ phaseLifecycle ∈ RUNNING_PHASE → canPause=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "running", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-011: jobLifecycle ∈ TERMINAL_JOB のとき canPause=false かつ pauseBlockedReason が終端文を返す
  test("jobLifecycle が completed のとき中断不可でジョブ終端理由を返す（H-13 TERMINAL_JOB 境界）", () => {
    // H-13: jobLifecycle ∈ TERMINAL_JOB → canPause=false、pauseBlockedReason=「ジョブが終端状態のため中断できません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "running", jobLifecycle: "completed" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のため中断できません。")
  })

  // RAEF-UNIT-EXT-012: phaseLifecycle ∉ RUNNING_PHASE かつ ¬terminal のとき canPause=false かつ pauseBlockedReason が実行中でない文を返す
  test("phaseLifecycle が paused かつ ¬terminal のとき中断不可でフェーズ実行中でない理由を返す（H-13 ¬RUNNING_PHASE 境界）", () => {
    // H-13: phaseLifecycle ∉ RUNNING_PHASE ∧ ¬terminal → canPause=false、pauseBlockedReason=「フェーズが実行中ではありません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "paused", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが実行中ではありません。")
  })
})

describe("deriveBodyActionEnablement — H-14 resume 境界値（RAEF-UNIT-EXT-013〜016）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  // RAEF-UNIT-EXT-013: phaseLifecycle ∈ PAUSED_PHASE かつ ¬terminal のとき canResume=true
  test("phaseLifecycle が paused かつ ¬terminal のとき再開可になる（H-14 PAUSED_PHASE 境界）", () => {
    // H-14: phaseLifecycle ∈ PAUSED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "paused", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-014: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE かつ ¬terminal のとき canResume=true
  test("phaseLifecycle が recoverable_failed かつ ¬terminal のとき再開可になる（H-14 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-14: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "recoverable_failed", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-015: errorKind=recoverable かつ ¬terminal のとき canResume=true
  test("errorKind が recoverable かつ ¬terminal のとき再開可になる（H-14 errorKind=recoverable 境界）", () => {
    // H-14: errorKind=recoverable ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "running", errorKind: "recoverable", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
  })

  // RAEF-UNIT-EXT-016: jobLifecycle ∈ TERMINAL_JOB のとき canResume=false かつ resumeBlockedReason が終端文を返す
  test("jobLifecycle が failed のとき再開不可でジョブ終端理由を返す（H-14 TERMINAL_JOB 境界）", () => {
    // H-14: jobLifecycle ∈ TERMINAL_JOB → canResume=false、resumeBlockedReason=「ジョブが終端状態のため再開できません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "paused", jobLifecycle: "failed" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のため再開できません。")
  })
})

describe("deriveBodyActionEnablement — H-15 retry 境界値（RAEF-UNIT-EXT-017〜018）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  // RAEF-UNIT-EXT-017: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE かつ ¬terminal のとき canRetry=true
  test("phaseLifecycle が recoverable_failed かつ ¬terminal のとき再試行可になる（H-15 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-15: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canRetry=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "recoverable_failed", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-018: phaseLifecycle ∈ PAUSED_PHASE かつ errorKind=none のとき canRetry=false かつ retryBlockedReason が再試行不可文を返す
  test("phaseLifecycle が paused かつ errorKind=none のとき再試行不可でフェーズ再試行不可理由を返す（H-15 ¬recoverableFailed 境界）", () => {
    // H-15: phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE ∧ errorKind ≠ recoverable → canRetry=false、retryBlockedReason=「フェーズが再試行可能な状態ではありません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "paused", errorKind: "none", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが再試行可能な状態ではありません。")
  })
})

describe("deriveBodyActionEnablement — H-16 cancel 境界値（RAEF-UNIT-EXT-019〜023）", () => {
  const presenter = new BodyTranslationPhasePresenter()

  // RAEF-UNIT-EXT-019: phaseLifecycle ∈ RUNNING_PHASE かつ ¬terminal のとき canCancel=true
  test("phaseLifecycle が running かつ ¬terminal のときキャンセル可になる（H-16 RUNNING_PHASE 境界）", () => {
    // H-16: phaseLifecycle ∈ RUNNING_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canCancel=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "running", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-020: phaseLifecycle ∈ PAUSED_PHASE かつ ¬terminal のとき canCancel=true
  test("phaseLifecycle が paused かつ ¬terminal のときキャンセル可になる（H-16 PAUSED_PHASE 境界）", () => {
    // H-16: phaseLifecycle ∈ PAUSED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canCancel=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "paused", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-021: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE かつ ¬terminal のとき canCancel=true
  test("phaseLifecycle が recoverable_failed かつ ¬terminal のときキャンセル可になる（H-16 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-16: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canCancel=true
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "recoverable_failed", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-022: phaseLifecycle が idle_ready（非アクティブ）かつ ¬terminal のとき canCancel=false かつ cancelBlockedReason がキャンセル不可文を返す
  test("phaseLifecycle が idle_ready かつ ¬terminal のときキャンセル不可でフェーズキャンセル不可理由を返す（H-16 ¬active 境界）", () => {
    // H-16: phaseLifecycle ∉ RUNNING_PHASE ∧ phaseLifecycle ∉ PAUSED_PHASE ∧ phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE ∧ errorKind ≠ recoverable → canCancel=false、cancelBlockedReason=「フェーズがキャンセル可能な状態ではありません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "idle_ready", errorKind: "none", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズがキャンセル可能な状態ではありません。")
  })

  // RAEF-UNIT-EXT-023: jobLifecycle ∈ TERMINAL_JOB のとき canCancel=false かつ cancelBlockedReason が終端文を返す
  test("jobLifecycle が canceled のときキャンセル不可でジョブ終端理由を返す（H-16 TERMINAL_JOB 境界）", () => {
    // H-16: jobLifecycle ∈ TERMINAL_JOB → canCancel=false、cancelBlockedReason=「ジョブが終端状態のためキャンセルできません。」
    const vm = presenter.toViewModel(
      buildBodyStateWithProjection({ phaseLifecycle: "running", jobLifecycle: "canceled" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のためキャンセルできません。")
  })
})
