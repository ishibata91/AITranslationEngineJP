import { describe, expect, test } from "vitest"

import { TranslationJobManagementPresenter } from "./translation-job-management.presenter"

type PresenterState = Parameters<
  TranslationJobManagementPresenter["toViewModel"]
>[0]

function createState(): PresenterState {
  return {
    phase: "ready",
    detailPhase: "ready",
    isReloading: false,
    selectedJobId: 10,
    filterId: "all",
    searchQuery: "",
    activeOperation: null,
    isDeleteConfirmationOpen: false,
    feedback: null,
    jobs: [
      {
        jobId: 10,
        jobState: "Running",
        jobStateLabel: "実行中",
        stateTone: "info",
        inputSource: {
          inputSourceId: 100,
          inputSourceLabel: "input.json",
          inputSourceKindLabel: "xEdit 抽出データ",
          sourcePath:
            "/very/long/path/to/source/that/should/not/be-truncated/by-presenter/input.json",
          pluginName:
            "VeryLongPluginNameThatShouldStayVisibleWithoutMutation.esp",
          extractedJsonLabel: "抽出データ #100"
        },
        progress: {
          currentPhase: "body_translation",
          currentPhaseLabel: "本文翻訳",
          percent: 45,
          progressLabel: "45%",
          lastUpdatedLabel: "2026-05-06 10:00"
        },
        stopAvailability: {
          kind: "stop",
          enabled: true,
          label: "停止",
          helperText: "停止可能"
        },
        resumeAvailability: {
          kind: "resume",
          enabled: false,
          label: "再開",
          helperText: "再開不可",
          reasonCategory: "terminal_state",
          reasonText: "完了扱いまたは回復不能な状態です"
        },
        deleteAvailability: {
          kind: "delete",
          enabled: false,
          label: "削除",
          helperText: "削除不可",
          reasonCategory: "running_delete_blocked",
          reasonText: "Running ジョブは削除できません"
        }
      }
    ],
    selectedJobDetail: {
      jobId: 10,
      jobState: "Running",
      jobStateLabel: "実行中",
      stateTone: "info",
      inputSource: {
        inputSourceId: 100,
        inputSourceLabel: "input.json",
        inputSourceKindLabel: "xEdit 抽出データ",
        sourcePath:
          "/very/long/path/to/source/that/should/not/be-truncated/by-presenter/input.json",
        pluginName: "VeryLongPluginNameThatShouldStayVisibleWithoutMutation.esp",
        extractedJsonLabel: "抽出データ #100"
      },
      progress: {
        currentPhase: "body_translation",
        currentPhaseLabel: "本文翻訳",
        percent: 45,
        progressLabel: "45%",
        lastUpdatedLabel: "2026-05-06 10:00"
      },
      stopAvailability: {
        kind: "stop",
        enabled: true,
        label: "停止",
        helperText: "停止可能"
      },
      resumeAvailability: {
        kind: "resume",
        enabled: false,
        label: "再開",
        helperText: "再開不可",
        reasonCategory: "cache_missing",
        reasonText: "入力キャッシュがありません"
      },
      deleteAvailability: {
        kind: "delete",
        enabled: false,
        label: "削除",
        helperText: "削除不可",
        reasonCategory: "running_delete_blocked",
        reasonText: "Running ジョブは削除できません"
      },
      cacheState: "missing",
      cacheStateLabel: "欠落",
      runtimeSummary: {
        providerLabel: "openai",
        modelLabel: "gpt-5",
        executionModeLabel: "batch",
        credentialState: "configured",
        credentialStateLabel: "設定済み"
      },
      resumeBlockedReasons: [
        {
          category: "cache_missing",
          title: "入力キャッシュが欠落しています",
          detail: "再開前に再構築してください"
        }
      ],
      warnings: [
        {
          category: "phase_progress_aggregation_failed",
          title: "進捗を集約できません",
          detail: "翻訳段階の進捗集約に失敗しました。再読込してください。"
        }
      ],
      deleteImpactLines: ["ジョブのみ削除", "入力データは残る"]
    }
  }
}

describe("TranslationJobManagementPresenter", () => {
  test("一覧用の表示項目、操作可否、無効理由、長文を view model に保持する", () => {
    const presenter = new TranslationJobManagementPresenter()

    const viewModel = presenter.toViewModel(createState(), true)

    expect(viewModel.jobs).toHaveLength(1)
    expect(viewModel.jobs[0].title).toBe("input.json")
    expect(viewModel.jobs[0].deleteOperation.enabled).toBe(false)
    expect(viewModel.jobs[0].deleteOperation.reasonText).toContain("削除できません")
    expect(viewModel.selectedJob?.cacheStateLabel).toBe("欠落")
    expect(viewModel.selectedJob?.resumeBlockedReasons[0]?.categoryLabel).toBe(
      "入力キャッシュ欠落"
    )
    expect(viewModel.selectedJob?.warnings[0]?.category).toBe(
      "phase_progress_aggregation_failed"
    )
    expect(viewModel.selectedJob?.warnings[0]?.categoryLabel).toBe(
      "進捗集約失敗"
    )
    expect(viewModel.jobRunTarget?.sourcePath).toContain("/very/long/path")
    expect(viewModel.selectedJob?.stopOperation.enabled).toBe(true)
  })

  test("フィルタと検索で一覧件数を変える", () => {
    const presenter = new TranslationJobManagementPresenter()
    const state = createState()
    state.filterId = "Running"
    state.searchQuery = "input.json"

    const viewModel = presenter.toViewModel(state, true)

    expect(viewModel.jobs).toHaveLength(1)
    expect(viewModel.headerCountLabel).toBe("1 件を表示")
  })

  test("Completed 用 filter を持たず、未完了選択時だけ phase page target を作る", () => {
    const presenter = new TranslationJobManagementPresenter()
    const state = createState()

    const viewModel = presenter.toViewModel(state, true)

    expect(viewModel.filterChips.some((chip) => chip.id === ("Completed" as never))).toBe(false)
    expect(viewModel.jobRunTarget).toMatchObject({
      jobId: 10,
      currentPhase: "body_translation",
      currentPhaseLabel: "本文翻訳"
    })
  })

  test("detail loading 中は一覧 summary から phase page target を維持する", () => {
    const presenter = new TranslationJobManagementPresenter()
    const state = createState()
    state.detailPhase = "loading"
    state.selectedJobDetail = null

    const viewModel = presenter.toViewModel(state, true)

    expect(viewModel.selectedJob).toBeNull()
    expect(viewModel.jobRunTarget).toMatchObject({
      jobId: 10,
      currentPhase: "body_translation",
      currentPhaseLabel: "本文翻訳"
    })
  })

  test("開けない job は理由を保持し、phase page target を作らない", () => {
    const presenter = new TranslationJobManagementPresenter()
    const state = createState()
    state.jobs[0].canOpenPhase = false
    state.jobs[0].openBlockedReason = {
      category: "phase_progress_aggregation_failed",
      title: "翻訳段階を開けません",
      detail: "進捗を確認できないため一覧で確認してください。"
    }
    if (state.selectedJobDetail) {
      state.selectedJobDetail.canOpenPhase = false
      state.selectedJobDetail.openBlockedReason =
        state.jobs[0].openBlockedReason
    }

    const viewModel = presenter.toViewModel(state, true)

    expect(viewModel.jobs[0].canOpenPhase).toBe(false)
    expect(viewModel.jobs[0].jobRunTarget).toBeNull()
    expect(viewModel.jobs[0].openBlockedReasonText).toContain("進捗")
    expect(viewModel.jobRunTarget).toBeNull()
  })

  test("action response 後の一覧状態でも phase page target を作らない", () => {
    const presenter = new TranslationJobManagementPresenter()
    const state = createState()
    const openBlockedReason = {
      category: "phase_progress_aggregation_failed" as const,
      title: "翻訳段階を開けません",
      detail: "進捗を確認できないため一覧で確認してください。"
    }
    state.jobs[0] = {
      ...state.jobs[0],
      canOpenPhase: false,
      openBlockedReason
    }
    if (state.selectedJobDetail) {
      state.selectedJobDetail = {
        ...state.selectedJobDetail,
        canOpenPhase: false,
        openBlockedReason
      }
    }

    const viewModel = presenter.toViewModel(state, true)

    expect(viewModel.jobs[0].canOpenPhase).toBe(false)
    expect(viewModel.jobs[0].openBlockedReason).toEqual(openBlockedReason)
    expect(viewModel.jobs[0].jobRunTarget).toBeNull()
    expect(viewModel.jobRunTarget).toBeNull()
  })
})
