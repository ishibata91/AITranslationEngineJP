import { describe, expect, test } from "vitest"
import type {
  PersonaGenerationPhaseSummaryResponse
} from "@application/gateway-contract/persona-generation-phase"
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
            canCancel: false
          }
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
          actionEnablement: {
            canStart: true,
            canPause: false,
            canResume: false,
            canRetry: false,
            canCancel: false
          }
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
          actionEnablement: {
            canStart: false,
            canPause: true,
            canResume: false,
            canRetry: false,
            canCancel: true
          }
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
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: false,
        canCancel: false
      },
      ...overrides
    }
  }

  function buildState(summaryOverrides: Partial<PersonaGenerationPhaseSummaryResponse> = {}) {
    return {
      jobId: 10,
      phase: "ready" as const,
      summary: buildBaseSummary(summaryOverrides),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      initialFetchDone: true
    }
  }

  // UT-EQV-008: 成立条件 - ジョブ終端でない・resultSummary.bodyReadiness=true のとき本文翻訳段階を開始できる
  test("ジョブ終端でなく bodyReadiness=true のとき本文翻訳段階を開始できる", () => {
    // 成立条件: isPersonaTerminalJob=false かつ resultSummary.bodyReadiness=true
    const vm = presenter.toViewModel(buildState(), true)

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(false)
    expect(vm.bodyReadinessLabel).toBe("Ready")
  })

  // UT-EQV-008: bodyReadiness=false のとき本文翻訳段階を開始できない（ペルソナ snapshot 参照未準備）
  test("resultSummary.bodyReadiness=false のとき本文翻訳段階を開始できない", () => {
    const vm = presenter.toViewModel(
      buildState({
        resultSummary: {
          generatedCount: 5,
          failedCount: 0,
          personaCount: 5,
          missingCount: 0,
          snapshotId: "persona-snapshot",
          snapshotDigest: "sha256:persona",
          snapshotReferenceStatus: "missing",
          bodyReadiness: false
        }
      }),
      true
    )

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(true)
    expect(vm.bodyReadinessLabel).toBe("Blocked")
    expect(vm.bodyReadinessBlockedReason).toContain("ペルソナ snapshot")
  })

  // UT-EQV-008: ジョブが終端状態のとき本文翻訳段階を開始できない（terminal_job 相当）
  test("startBlockedReason=terminal_job のとき本文翻訳段階を開始できずジョブ終端理由を返す", () => {
    const vm = presenter.toViewModel(
      buildState({
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

    const startBodyCard = vm.actionCards.find((c) => c.id === "start-body-phase")
    expect(startBodyCard?.disabled).toBe(true)
    expect(vm.bodyReadinessBlockedReason).toContain("終端")
  })

  // UT-EQV-008: 操作可否 - running のとき canPause=true かつ canResume=false
  test("phaseState=running のとき一時停止可かつ再開不可", () => {
    const vm = presenter.toViewModel(buildState({ phaseState: "running" }), true)

    const pauseCard = vm.actionCards.find((c) => c.id === "pause")
    const resumeCard = vm.actionCards.find((c) => c.id === "resume")
    expect(pauseCard?.disabled).toBe(false)
    expect(resumeCard?.disabled).toBe(true)
  })

  // UT-EQV-008: 操作可否 - paused のとき canResume=true かつ canPause=false
  test("phaseState=paused のとき再開可かつ一時停止不可", () => {
    const vm = presenter.toViewModel(buildState({ phaseState: "paused" }), true)

    const pauseCard = vm.actionCards.find((c) => c.id === "pause")
    const resumeCard = vm.actionCards.find((c) => c.id === "resume")
    expect(pauseCard?.disabled).toBe(true)
    expect(resumeCard?.disabled).toBe(false)
  })

  // UT-EQV-008: 操作可否 - recoverable_failed のとき canRetry=true かつ canCancel=true
  test("phaseState=recoverable_failed のとき再試行可かつキャンセル可", () => {
    const vm = presenter.toViewModel(buildState({ phaseState: "recoverable_failed" }), true)

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
      actionEnablement: {
        canStart: false,
        canPause: false,
        canResume: false,
        canRetry: false,
        canCancel: false
      }
    }
  }

  function buildStateWithoutExecution() {
    return {
      jobId: 10,
      phase: "ready" as const,
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
