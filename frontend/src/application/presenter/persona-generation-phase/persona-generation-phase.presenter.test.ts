import { describe, expect, test } from "vitest"
import type {
  PersonaGenerationPhaseProjection,
  PersonaGenerationPhaseSummaryResponse
} from "@application/gateway-contract/persona-generation-phase"
import { PersonaGenerationPhasePresenter } from "./persona-generation-phase.presenter"

const presenter = new PersonaGenerationPhasePresenter()

const defaultProjection: PersonaGenerationPhaseProjection = {
  phaseLifecycle: "pending",
  jobLifecycle: "running",
  errorKind: "none",
  aiSettingsConfigured: true,
  targetCount: 5,
  previousPhaseLifecycle: "completed"
}

describe("PersonaGenerationPhasePresenter", () => {
  test("未接続状態を view model に反映する", () => {
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready",
        projection: null,
        summary: null,
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: false,
        initialFetchDone: false
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
        projection: { ...defaultProjection, phaseLifecycle: "completed" },
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
        },
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true
      },
      true
    )

    expect(vm.viewState).toBe("snapshot_missing")
    expect(vm.errorKindLabel).toContain("snapshot_missing")
  })

  test("pending phase state は開始待ちとして表示する", () => {
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready",
        projection: defaultProjection,
        summary: {
          jobId: 10,
          currentPhase: "persona_generation",
          phaseState: "pending",
          phaseRunId: 20,
          progress: {
            percent: 0,
            processedCount: 0,
            totalCount: 1,
            targetCount: 1,
            currentStep: "pending"
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
        },
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true
      },
      true
    )

    expect(vm.viewState).toBe("not_started")
    expect(vm.phaseStateLabel).toBe("開始待ち")
    expect(vm.progressDetail).toContain("開始待ち")
    expect(vm.progressDetail).not.toContain("pending")
  })

  test("進行詳細は処理対象件数と一覧 total を別主語として検出できる表示を返す", () => {
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready",
        projection: { ...defaultProjection, phaseLifecycle: "running" },
        summary: {
          jobId: 10,
          currentPhase: "persona_generation",
          phaseState: "running",
          progress: {
            percent: 60,
            processedCount: 12,
            totalCount: 30,
            targetCount: 18,
            currentStep: "running"
          },
          targetSummary: {
            targetCount: 18,
            commonPersonaHitCount: 2,
            commonPersonaMissCount: 16,
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
            inputCount: 18,
            outputCount: 12,
            evidenceRefs: []
          },
        },
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true,
        processingTargetPageState: {
          items: [],
          metadata: [],
          page: 1,
          pageSize: 20,
          totalCount: 18,
          searchQuery: "",
          busy: false
        }
      },
      true
    )

    expect(vm.progressDetail).toBe("12 / 30 件 / 対象 18 件 / 実行中")
    expect(vm.targetCountLabel).toBe("18 件")
  })
})

// UT-EQV-008: persona の本文翻訳段階開始可否と操作可否の等価性テスト
// 期待値は detail-spec persona-generation-phase-REQ-008 成立条件・区別理由から固定する
// 注記: FE-persona の実装変更により、canStartBodyPhase の判定は projection.phaseLifecycle に基づく（design-diff G-2-b / H-11）
describe("PersonaGenerationPhasePresenter - 本文翻訳段階開始可否・操作可否の等価性（UT-EQV-008）", () => {
  const presenter = new PersonaGenerationPhasePresenter()

  function buildBaseSummary(
    overrides: Partial<PersonaGenerationPhaseSummaryResponse> = {}
  ): PersonaGenerationPhaseSummaryResponse {
    return {
      jobId: 10,
      currentPhase: "persona_generation",
      phaseState: "completed",
      progress: {
        percent: 100,
        processedCount: 5,
        totalCount: 5,
        targetCount: 5,
        currentStep: "completed"
      },
      targetSummary: {
        targetCount: 5,
        commonPersonaHitCount: 0,
        commonPersonaMissCount: 5,
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
        inputCount: 5,
        outputCount: 5,
        evidenceRefs: []
      },
      resultSummary: {
        generatedCount: 5,
        failedCount: 0,
        personaCount: 5,
        missingCount: 0,
        snapshotId: "persona-snapshot",
        snapshotDigest: "sha256:persona",
        snapshotReferenceStatus: "ready",
        bodyReadiness: true
      },
      ...overrides
    }
  }

  function buildBaseProjection(
    overrides: Partial<PersonaGenerationPhaseProjection> = {}
  ): PersonaGenerationPhaseProjection {
    return {
      phaseLifecycle: "completed",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      targetCount: 5,
      previousPhaseLifecycle: "completed",
      ...overrides
    }
  }

  function buildState(
    summaryOverrides: Partial<PersonaGenerationPhaseSummaryResponse> = {},
    projectionOverrides: Partial<PersonaGenerationPhaseProjection> = {}
  ) {
    return {
      jobId: 10,
      phase: "ready" as const,
      projection: buildBaseProjection(projectionOverrides),
      summary: buildBaseSummary(summaryOverrides),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // UT-EQV-008: 成立条件 - ジョブ終端でない・phaseLifecycle=completed のとき本文翻訳段階を開始できる
  // 注記: 新設計では canStartBodyPhase の判定は projection.phaseLifecycle ∈ COMPLETED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB（H-11）
  test("ジョブ終端でなく phaseLifecycle=completed のとき本文翻訳段階を開始できる", () => {
    const vm = presenter.toViewModel(buildState(), true)

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(false)
    expect(vm.bodyReadinessLabel).toBe("Ready")
  })

  // UT-EQV-008: phaseLifecycle が完了系でないとき本文翻訳段階を開始できない
  // 注記: 新設計では bodyReadiness は body 側 projection が担うため、persona 側は phaseLifecycle だけを確認する（H-11）
  test("phaseLifecycle=running のとき本文翻訳段階を開始できない", () => {
    const vm = presenter.toViewModel(
      buildState({}, { phaseLifecycle: "running" }),
      true
    )

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(true)
    expect(vm.bodyReadinessLabel).toBe("Blocked")
    expect(vm.bodyReadinessBlockedReason).toContain("未完了")
  })

  // UT-EQV-008: ジョブが終端状態のとき本文翻訳段階を開始できない（terminal_job 相当）
  test("jobLifecycle=failed のとき本文翻訳段階を開始できずジョブ終端理由を返す", () => {
    const vm = presenter.toViewModel(
      buildState({}, { jobLifecycle: "failed" }),
      true
    )

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(true)
    expect(vm.bodyReadinessBlockedReason).toContain("終端")
  })

  // UT-EQV-008: 操作可否 - running のとき canPause=true かつ canResume=false
  test("phaseLifecycle=running のとき一時停止可かつ再開不可", () => {
    const vm = presenter.toViewModel(buildState({}, { phaseLifecycle: "running" }), true)

    const pauseCard = vm.actionCards.find((c) => c.id === "pause")
    const resumeCard = vm.actionCards.find((c) => c.id === "resume")
    expect(pauseCard?.disabled).toBe(false)
    expect(resumeCard?.disabled).toBe(true)
  })

  // UT-EQV-008: 操作可否 - paused のとき canResume=true かつ canPause=false
  test("phaseLifecycle=paused のとき再開可かつ一時停止不可", () => {
    const vm = presenter.toViewModel(buildState({}, { phaseLifecycle: "paused" }), true)

    const pauseCard = vm.actionCards.find((c) => c.id === "pause")
    const resumeCard = vm.actionCards.find((c) => c.id === "resume")
    expect(pauseCard?.disabled).toBe(true)
    expect(resumeCard?.disabled).toBe(false)
  })

  // UT-EQV-008: 操作可否 - recoverable_failed のとき canRetry=true かつ canCancel=true
  test("phaseLifecycle=recoverable_failed のとき再試行可かつキャンセル可", () => {
    const vm = presenter.toViewModel(buildState({}, { phaseLifecycle: "recoverable_failed" }), true)

    const retryCard = vm.actionCards.find((c) => c.id === "retry")
    const cancelCard = vm.actionCards.find((c) => c.id === "cancel")
    expect(retryCard?.disabled).toBe(false)
    expect(cancelCard?.disabled).toBe(false)
  })
})

// U-PRES-002/003/004: execution field 不在時の presenter 判定（ペルソナ生成）
// 根拠: persona-generation-phase-REQ-002「execution field 不在で isExecutionConfigured=false」
describe("PersonaGenerationPhasePresenter - execution field 不在判定（U-PRES-002/003/004）", () => {
  const presenter = new PersonaGenerationPhasePresenter()

  function buildSummaryWithoutExecution(): PersonaGenerationPhaseSummaryResponse {
    return {
      jobId: 10,
      currentPhase: "persona_generation",
      phaseState: "ready",
      progress: {
        percent: 0,
        processedCount: 0,
        totalCount: 5,
        targetCount: 5,
        currentStep: "ready"
      },
      targetSummary: {
        targetCount: 5,
        commonPersonaHitCount: 0,
        commonPersonaMissCount: 5,
        skippedCount: 0,
        skippedReasons: [],
        targetSnapshotDigest: "sha256:1"
      },
      // execution field を意図的に含めない（不在で「未設定」を表す仕様）
    }
  }

  function buildStateWithoutExecution() {
    return {
      jobId: 10,
      phase: "ready" as const,
      projection: null,
      summary: buildSummaryWithoutExecution(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-002: execution field が不在の場合、isExecutionConfigured が false を返す
  test("execution field が不在のとき isExecutionConfigured が false を返す", () => {
    // 空文字ではなく field 不在で未設定を判定できることを証明する。
    const vm = presenter.toViewModel(buildStateWithoutExecution(), true)

    expect(vm.isExecutionConfigured).toBe(false)
  })

  // U-PRES-003: execution field が不在の場合、modelLabel が「設定未完了」を返す（空文字 "" を返さない）
  test("execution field が不在のとき modelLabel が設定未完了を返す", () => {
    // 空文字フォールバックを防止し「設定未完了」相当の表示語を返すことを証明する。
    const vm = presenter.toViewModel(buildStateWithoutExecution(), true)

    expect(vm.modelLabel).not.toBe("")
    expect(vm.modelLabel).toBe("設定未完了")
  })

  // U-PRES-002（派生）: execution field が不在の場合と存在する場合を presenter が独立して判定できる
  test("execution field が存在しモデル値が入力済みのとき isExecutionConfigured が true を返す", () => {
    // 2 つの不在状態の混同を防ぐ。設定済み状態と未設定状態を独立して判定できることを証明する。
    const summaryWithExecution: PersonaGenerationPhaseSummaryResponse = {
      ...buildSummaryWithoutExecution(),
      execution: {
        credentialRef: "cred",
        provider: "fake",
        model: "m",
        executionMode: "single_request",
        promptDigest: "sha256:1",
        inputCount: 5,
        outputCount: 0,
        evidenceRefs: []
      }
    }
    const vm = presenter.toViewModel(
      {
        jobId: 10,
        phase: "ready" as const,
        projection: null,
        summary: summaryWithExecution,
        bodyReadiness: null,
        errorMessage: "",
        pendingAction: null,
        hasLoaded: true,
        initialFetchDone: true
      },
      true
    )

    // execution field が存在し値が入力済み → isExecutionConfigured=true
    expect(vm.isExecutionConfigured).toBe(true)
    expect(vm.modelLabel).toBe("m")
  })
})

// U-PRES-MODEL-002: buildPersonaModelOptions の placeholder 重複解消（ペルソナ生成）
// 根拠: fix-phase-ai-settings-pill-update-after-model-select の secondary 問題修正
// AIModelSelectionCard が固定 placeholder を持つため、buildPersonaModelOptions が重複 placeholder を返さないことを証明する
describe("PersonaGenerationPhasePresenter - buildPersonaModelOptions placeholder 重複解消（U-PRES-MODEL-002）", () => {
  const presenter = new PersonaGenerationPhasePresenter()

  function buildState() {
    return {
      jobId: 10,
      phase: "ready" as const,
      projection: {
        phaseLifecycle: "ready",
        jobLifecycle: "running",
        errorKind: "none",
        aiSettingsConfigured: false,
        targetCount: 5,
        previousPhaseLifecycle: "completed"
      },
      summary: {
        jobId: 10,
        currentPhase: "persona_generation" as const,
        phaseState: "ready" as const,
        progress: {
          percent: 0,
          processedCount: 0,
          totalCount: 5,
          targetCount: 5,
          currentStep: "ready"
        },
        targetSummary: {
          targetCount: 5,
          commonPersonaHitCount: 0,
          commonPersonaMissCount: 5,
          skippedCount: 0,
          skippedReasons: [],
          targetSnapshotDigest: "sha256:1"
        }
      },
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // U-PRES-MODEL-002: availableModels に要素がある時、modelOptions に placeholder が含まれない
  test("availableModels に要素がある時 modelOptions に空 value の placeholder が含まれない", () => {
    // AIModelSelectionCard が固定 placeholder を描画するため、buildPersonaModelOptions は重複 placeholder を返さない。
    const availableModels = [{ value: "fake-model", label: "fake-model" }]

    const vm = presenter.toViewModel(buildState(), true, [], availableModels)

    // availableModels を渡した時 placeholder が含まれないことを証明する
    const hasPlaceholder = vm.modelOptions.some((opt) => opt.value === "")
    expect(hasPlaceholder).toBe(false)
  })

  // U-PRES-MODEL-002: availableModels に要素がある時、modelOptions は availableModels の内容を返す
  test("availableModels に要素がある時 modelOptions が availableModels をそのまま返す", () => {
    // placeholder を挿入せず availableModels の内容だけを返すことを証明する。
    const availableModels = [
      { value: "model-a", label: "Model A" },
      { value: "model-b", label: "Model B" }
    ]

    const vm = presenter.toViewModel(buildState(), true, [], availableModels)

    expect(vm.modelOptions).toHaveLength(2)
    expect(vm.modelOptions).toContainEqual({ value: "model-a", label: "Model A" })
    expect(vm.modelOptions).toContainEqual({ value: "model-b", label: "Model B" })
  })
})

// RAEF-UNIT-010〜014: derivePersonaActionEnablement および derivePersonaCanStartBodyPhase の境界値証明
// 根拠: design-diff.md H-6〜H-11

// persona 画面の state を構築するヘルパー。projection のみ差し替えて検証する。
function buildPersonaState(
  projectionOverrides: Partial<PersonaGenerationPhaseProjection> = {}
) {
  return {
    jobId: 10,
    phase: "ready" as const,
    projection: { ...defaultProjection, ...projectionOverrides },
    summary: null as PersonaGenerationPhaseSummaryResponse | null,
    bodyReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    initialFetchDone: true
  }
}

describe("derivePersonaActionEnablement — H-6 start 境界値（RAEF-UNIT-010〜011）", () => {
  // RAEF-UNIT-010: previousPhaseLifecycle ∉ COMPLETED_PHASE のとき canStart=false かつ startBlockedReason に段階未完了文を返す
  test("previousPhaseLifecycle が running のとき開始不可で単語翻訳段階未完了理由を返す（H-6 previousPhaseCompleted=false 境界）", () => {
    // H-6: previousPhaseLifecycle ∉ COMPLETED_PHASE → canStart=false、startBlockedReason=「単語翻訳段階が完了していないためペルソナ生成を開始できません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ previousPhaseLifecycle: "running", phaseLifecycle: "idle_ready" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("単語翻訳段階が完了していないためペルソナ生成を開始できません。")
  })

  // RAEF-UNIT-011: 全条件満足（¬terminal ∧ idleReady ∧ aiSettingsConfigured ∧ targetCount=1 ∧ previousPhaseCompleted）のとき canStart=true
  test("全条件満足のとき開始可になる（H-6 正常パス、targetCount 最小値 1）", () => {
    // H-6: 全有効化条件を満たす境界値（targetCount=1）で canStart=true
    const vm = presenter.toViewModel(
      buildPersonaState({
        phaseLifecycle: "idle_ready",
        jobLifecycle: "running",
        aiSettingsConfigured: true,
        targetCount: 1,
        previousPhaseLifecycle: "completed"
      }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start")

    expect(card?.disabled).toBe(false)
  })
})

describe("derivePersonaActionEnablement — H-10 cancel 境界値（RAEF-UNIT-012）", () => {
  // RAEF-UNIT-012: RUNNING_PHASE のとき canCancel=true、jobLifecycle ∈ TERMINAL_JOB のとき canCancel=false かつ cancelBlockedReason を返す
  test("phaseLifecycle が running のときキャンセル可になる（H-10 RUNNING_PHASE 境界）", () => {
    // H-10: phaseLifecycle ∈ RUNNING_PHASE → canCancel=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(false)
  })

  test("jobLifecycle が completed のときキャンセル不可でジョブ終端理由を返す（H-10 TERMINAL_JOB 境界）", () => {
    // H-10: jobLifecycle ∈ TERMINAL_JOB → canCancel=false、cancelBlockedReason=「ジョブが終端状態のためキャンセルできません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", jobLifecycle: "completed" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "cancel")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のためキャンセルできません。")
  })
})

describe("derivePersonaActionEnablement — H-7 pause 境界値（RAEF-UNIT-EXT-001〜003）", () => {
  // RAEF-UNIT-EXT-001: phaseLifecycle ∈ RUNNING_PHASE かつ ¬terminal のとき canPause=true
  test("phaseLifecycle が running かつ ¬terminal のとき中断可になる（H-7 RUNNING_PHASE 境界）", () => {
    // H-7: jobLifecycle ∉ TERMINAL_JOB ∧ phaseLifecycle ∈ RUNNING_PHASE → canPause=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-002: jobLifecycle ∈ TERMINAL_JOB のとき canPause=false かつ pauseBlockedReason が終端文を返す
  test("jobLifecycle が completed のとき中断不可でジョブ終端理由を返す（H-7 TERMINAL_JOB 境界）", () => {
    // H-7: jobLifecycle ∈ TERMINAL_JOB → canPause=false、pauseBlockedReason=「ジョブが終端状態のため中断できません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", jobLifecycle: "completed" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のため中断できません。")
  })

  // RAEF-UNIT-EXT-003: phaseLifecycle ∉ RUNNING_PHASE かつ ¬terminal のとき canPause=false かつ pauseBlockedReason が実行中でない文を返す
  test("phaseLifecycle が paused かつ ¬terminal のとき中断不可でフェーズ実行中でない理由を返す（H-7 ¬RUNNING_PHASE 境界）", () => {
    // H-7: phaseLifecycle ∉ RUNNING_PHASE ∧ ¬terminal → canPause=false、pauseBlockedReason=「フェーズが実行中ではありません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "paused", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "pause")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが実行中ではありません。")
  })
})

describe("derivePersonaActionEnablement — H-8 resume 境界値（RAEF-UNIT-EXT-004〜007）", () => {
  // RAEF-UNIT-EXT-004: phaseLifecycle ∈ PAUSED_PHASE かつ ¬terminal のとき canResume=true
  test("phaseLifecycle が paused かつ ¬terminal のとき再開可になる（H-8 PAUSED_PHASE 境界）", () => {
    // H-8: phaseLifecycle ∈ PAUSED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "paused", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-005: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE かつ ¬terminal のとき canResume=true
  test("phaseLifecycle が recoverable_failed かつ ¬terminal のとき再開可になる（H-8 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-8: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "recoverable_failed", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-006: errorKind=recoverable かつ ¬terminal のとき canResume=true（phaseLifecycle が RUNNING で errorKind が recoverable の場合）
  test("errorKind が recoverable かつ ¬terminal のとき再開可になる（H-8 errorKind=recoverable 境界）", () => {
    // H-8: errorKind=recoverable ∧ jobLifecycle ∉ TERMINAL_JOB → canResume=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", errorKind: "recoverable", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(false)
  })

  // RAEF-UNIT-EXT-007: jobLifecycle ∈ TERMINAL_JOB のとき canResume=false かつ resumeBlockedReason が終端文を返す
  test("jobLifecycle が failed のとき再開不可でジョブ終端理由を返す（H-8 TERMINAL_JOB 境界）", () => {
    // H-8: jobLifecycle ∈ TERMINAL_JOB → canResume=false、resumeBlockedReason=「ジョブが終端状態のため再開できません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "paused", jobLifecycle: "failed" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "resume")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ジョブが終端状態のため再開できません。")
  })
})

describe("derivePersonaActionEnablement — H-9 retry 境界値（RAEF-UNIT-EXT-008〜009）", () => {
  // RAEF-UNIT-EXT-008: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE かつ ¬terminal のとき canRetry=true
  test("phaseLifecycle が recoverable_failed かつ ¬terminal のとき再試行可になる（H-9 RECOVERABLE_FAILED_PHASE 境界）", () => {
    // H-9: phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canRetry=true
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "recoverable_failed", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-EXT-009: phaseLifecycle ∈ PAUSED_PHASE かつ errorKind=none のとき canRetry=false かつ retryBlockedReason が再試行不可文を返す
  test("phaseLifecycle が paused かつ errorKind=none のとき再試行不可でフェーズ再試行不可理由を返す（H-9 ¬recoverableFailed 境界）", () => {
    // H-9: phaseLifecycle ∉ RECOVERABLE_FAILED_PHASE ∧ errorKind ≠ recoverable → canRetry=false、retryBlockedReason=「フェーズが再試行可能な状態ではありません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "paused", errorKind: "none", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "retry")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("フェーズが再試行可能な状態ではありません。")
  })
})

describe("derivePersonaCanStartBodyPhase — H-11 境界値（RAEF-UNIT-013〜014）", () => {
  // RAEF-UNIT-013: COMPLETED_PHASE かつ ¬terminal のとき canStartNextPhase=true かつ personaBodyReadiness は persona 側 projection に含まれない
  test("phaseLifecycle が completed かつ ¬terminal のとき本文翻訳開始可になる（H-11 正常パス）", () => {
    // H-11: phaseLifecycle ∈ COMPLETED_PHASE ∧ jobLifecycle ∉ TERMINAL_JOB → canStartNextPhase=true
    // persona 側 projection に personaBodyReadiness が存在しないことを confirm する（選択 A）
    const projection = { ...defaultProjection, phaseLifecycle: "completed", jobLifecycle: "running" }
    expect("personaBodyReadiness" in projection).toBe(false)

    const vm = presenter.toViewModel(buildPersonaState({ phaseLifecycle: "completed", jobLifecycle: "running" }), true)
    const card = vm.actionCards.find((c) => c.id === "start-body-phase")

    expect(card?.disabled).toBe(false)
    expect(card?.blockedReason).toBe("")
  })

  // RAEF-UNIT-014: phaseLifecycle ∉ COMPLETED_PHASE のとき canStartNextPhase=false かつ 未完了理由を返す
  test("phaseLifecycle が running のとき本文翻訳開始不可でペルソナ生成未完了理由を返す（H-11 completed=false 境界）", () => {
    // H-11: phaseLifecycle ∉ COMPLETED_PHASE → canStartNextPhase=false、BlockedReason=「ペルソナ生成段階が未完了のため本文翻訳を開始できません。」
    const vm = presenter.toViewModel(
      buildPersonaState({ phaseLifecycle: "running", jobLifecycle: "running" }),
      true
    )
    const card = vm.actionCards.find((c) => c.id === "start-body-phase")

    expect(card?.disabled).toBe(true)
    expect(card?.blockedReason).toBe("ペルソナ生成段階が未完了のため本文翻訳を開始できません。")
  })
})
