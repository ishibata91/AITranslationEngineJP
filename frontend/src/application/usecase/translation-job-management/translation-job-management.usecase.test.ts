import { describe, expect, test, vi } from "vitest"

import type { TranslationJobManagementGatewayContract } from "@application/gateway-contract/translation-job-management"

import { TranslationJobManagementUseCase } from "./translation-job-management.usecase"

type UseCaseStoreLike = ConstructorParameters<
  typeof TranslationJobManagementUseCase
>[1]
type UseCaseState = ReturnType<UseCaseStoreLike["snapshot"]>

function createInitialState(): UseCaseState {
  return {
    phase: "idle",
    jobs: [],
    selectedJobId: null,
    selectedJobDetail: null,
    detailPhase: "idle",
    filterId: "all",
    searchQuery: "",
    isReloading: false,
    activeOperation: null,
    isDeleteConfirmationOpen: false,
    feedback: null
  }
}

function createStore(initial = createInitialState()) {
  const state: UseCaseState = structuredClone(initial)
  return {
    snapshot: () => structuredClone(state),
    update: (mutator: (draft: UseCaseState) => void) => {
      const draft = structuredClone(state)
      mutator(draft)
      Object.assign(state, draft)
    }
  }
}

function createJobSummary(jobId: number) {
  return {
    jobId,
    jobState: "Paused" as const,
    jobStateLabel: "中断中",
    stateTone: "warning" as const,
    inputSource: {
      inputSourceId: 1,
      inputSourceLabel: "input.json",
      inputSourceKindLabel: "xEdit 抽出データ",
      sourcePath: "/mods/input.json",
      pluginName: "Skyrim.esm",
      extractedJsonLabel: "抽出データ #1"
    },
    progress: {
      currentPhaseLabel: "本文翻訳",
      percent: 20,
      progressLabel: "20%",
      lastUpdatedLabel: "now"
    },
    stopAvailability: { kind: "stop" as const, enabled: false, label: "停止", helperText: "" },
    resumeAvailability: {
      kind: "resume" as const,
      enabled: false,
      label: "再開",
      helperText: "",
      reasonCategory: "cache_missing" as const,
      reasonText: "入力キャッシュがありません"
    },
    deleteAvailability: { kind: "delete" as const, enabled: true, label: "削除", helperText: "" }
  }
}

function createJobSummaryWithState(
  jobId: number,
  jobState: "Ready" | "Running" | "Paused" | "RecoverableFailed" | "Failed" | "Canceled"
) {
  const summary = createJobSummary(jobId)
  return {
    ...summary,
    jobState,
    jobStateLabel: jobState
  }
}

describe("TranslationJobManagementUseCase", () => {
  test("gateway 未接続で selectJob すると warning を設定する", async () => {
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(null, store)

    await usecase.selectJob(99)

    const state = store.snapshot()
    expect(state.selectedJobId).toBe(99)
    expect(state.feedback?.tone).toBe("warning")
  })

  test("setFilter と setSearchQuery は状態を更新する", () => {
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(null, store)

    usecase.setFilter("Running")
    usecase.setSearchQuery("input.json")

    const state = store.snapshot()
    expect(state.filterId).toBe("Running")
    expect(state.searchQuery).toBe("input.json")
  })

  test("gateway 未接続で load すると list_load_failure を設定する", async () => {
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(null, store)

    await usecase.load()

    const state = store.snapshot()
    expect(state.phase).toBe("error")
    expect(state.feedback?.category).toBe("list_load_failure")
  })

  test("load は gateway が返した未完了 job 一覧をそのまま表示対象へ保持する", async () => {
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({
        jobs: [
          createJobSummaryWithState(10, "Ready"),
          createJobSummaryWithState(11, "Running")
        ]
      }),
      GetJobDetail: vi.fn(),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()

    const state = store.snapshot()
    expect(state.phase).toBe("ready")
    expect(state.jobs.map((job) => job.jobId)).toEqual([10, 11])
    expect(state.jobs.map((job) => job.jobState)).toEqual(["Ready", "Running"])
  })

  test("reload で選択済み job が消えた時 stale_selection を設定する", async () => {
    const detail = {
      ...createJobSummary(10),
      cacheState: "available" as const,
      cacheStateLabel: "利用可能",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured" as const,
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    }
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi
        .fn()
        .mockResolvedValueOnce({ jobs: [createJobSummary(10)] })
        .mockResolvedValueOnce({ jobs: [] }),
      GetJobDetail: vi.fn().mockResolvedValue(detail),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(10)
    await usecase.reload()

    const state = store.snapshot()
    expect(state.detailPhase).toBe("stale")
    expect(state.feedback?.category).toBe("stale_selection")
    expect(state.isReloading).toBe(false)
  })

  test("stop 失敗時に stop_failed を設定する", async () => {
    const detail = {
      ...createJobSummary(11),
      cacheState: "available" as const,
      cacheStateLabel: "利用可能",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured" as const,
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    }
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({ jobs: [createJobSummary(11)] }),
      GetJobDetail: vi.fn().mockResolvedValue(detail),
      RequestStop: vi.fn().mockRejectedValue(new Error("stop failed")),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(11)
    await usecase.requestStop()

    const state = store.snapshot()
    expect(state.activeOperation).toBeNull()
    expect(state.feedback?.category).toBe("stop_failed")
  })

  test("一覧読み込み失敗で list_load_failure を設定する", async () => {
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockRejectedValue(new Error("network")),
      GetJobDetail: vi.fn(),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()

    const state = store.snapshot()
    expect(state.phase).toBe("error")
    expect(state.feedback?.category).toBe("list_load_failure")
  })

  test("選択後に detail 読み込み失敗で stale_selection を設定する", async () => {
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({ jobs: [createJobSummary(10)] }),
      GetJobDetail: vi.fn().mockRejectedValue(new Error("missing")),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(10)

    const state = store.snapshot()
    expect(state.detailPhase).toBe("stale")
    expect(state.feedback?.category).toBe("stale_selection")
  })

  test("削除成功で選択と詳細をクリアし delete_success を設定する", async () => {
    const detail = {
      ...createJobSummary(11),
      cacheState: "available" as const,
      cacheStateLabel: "利用可能",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured" as const,
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    }
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({ jobs: [createJobSummary(11)] }),
      GetJobDetail: vi.fn().mockResolvedValue(detail),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn().mockResolvedValue({
        message: "削除しました",
        tone: "success",
        deletedJobId: 11
      })
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(11)
    await usecase.deleteSelectedJob()

    const state = store.snapshot()
    expect(state.selectedJobId).toBeNull()
    expect(state.selectedJobDetail).toBeNull()
    expect(state.feedback?.category).toBe("delete_success")
  })

  test("削除拒否で detail を同期し reasonCategory を返す", async () => {
    const detail = {
      ...createJobSummary(22),
      jobState: "Running" as const,
      cacheState: "available" as const,
      cacheStateLabel: "利用可能",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured" as const,
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    }
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({ jobs: [createJobSummary(22)] }),
      GetJobDetail: vi.fn().mockResolvedValue(detail),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn(),
      DeleteJob: vi.fn().mockResolvedValue({
        message: "Running job は削除できません",
        tone: "warning",
        detail,
        reasonCategory: "running_delete_blocked"
      })
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(22)
    await usecase.deleteSelectedJob()

    const state = store.snapshot()
    expect(state.selectedJobDetail?.jobState).toBe("Running")
    expect(state.feedback?.category).toBe("running_delete_blocked")
  })

  test("resume 成功で reasonCategory 未指定時に resume_success を設定する", async () => {
    const detail = {
      ...createJobSummary(33),
      cacheState: "available" as const,
      cacheStateLabel: "利用可能",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured" as const,
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: ["job のみ削除"]
    }
    const gateway: TranslationJobManagementGatewayContract = {
      ListIncompleteJobs: vi.fn().mockResolvedValue({ jobs: [createJobSummary(33)] }),
      GetJobDetail: vi.fn().mockResolvedValue(detail),
      RequestStop: vi.fn(),
      ResumeJob: vi.fn().mockResolvedValue({
        message: "再開を受け付けました",
        tone: "success",
        detail
      }),
      DeleteJob: vi.fn()
    }
    const store = createStore()
    const usecase = new TranslationJobManagementUseCase(gateway, store)

    await usecase.load()
    await usecase.selectJob(33)
    await usecase.requestResume()

    const state = store.snapshot()
    expect(state.feedback?.category).toBe("resume_success")
  })
})
