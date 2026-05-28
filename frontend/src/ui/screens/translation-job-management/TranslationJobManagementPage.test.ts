import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  TranslationJobManagementScreenControllerContract,
  TranslationJobManagementScreenViewModelListener
} from "@application/contract/translation-job-management/translation-job-management-screen-contract"
import type {
  TranslationJobManagementJobRunTarget,
  TranslationJobManagementJobState,
  TranslationJobManagementScreenViewModel
} from "@application/contract/translation-job-management/translation-job-management-screen-types"

import TranslationJobManagementPage from "./TranslationJobManagementPage.svelte"

const resumeTarget: TranslationJobManagementJobRunTarget = {
  jobId: 41,
  stateLabel: "Running",
  stateDescription: "実行中",
  currentPhase: "body_translation",
  currentPhaseLabel: "本文翻訳",
  progressLabel: "20%",
  inputSourceLabel: "input.json",
  sourcePath: "/mods/input.json"
}

function createViewModel(
  feedback: TranslationJobManagementScreenViewModel["feedback"] = null,
  jobRunTarget: TranslationJobManagementJobRunTarget | null = null,
  jobOverride: Partial<
    TranslationJobManagementScreenViewModel["jobs"][number]
  > = {}
): TranslationJobManagementScreenViewModel {
  const cardJobRunTarget =
    jobOverride.jobRunTarget === undefined
      ? resumeTarget
      : jobOverride.jobRunTarget

  return {
    gatewayStatus: "接続済み",
    pageTitle: "未完了ジョブ一覧",
    pageLead:
      "新規翻訳を開始するか、未完了ジョブを選んで現在の翻訳段階から再開します。",
    headerCountLabel: "1 件を表示",
    listEmptyTitle: "管理対象がありません",
    listEmptyDescription: "未完了ジョブはありません。",
    listErrorTitle: "一覧を読み込めません",
    listErrorDescription: "未完了ジョブの一覧取得に失敗しました。",
    detailPlaceholderTitle: "job を選択してください",
    detailPlaceholderDescription: "一覧から 1 件選びます。",
    phase: "ready",
    detailPhase: "ready",
    isReloading: false,
    searchQuery: "",
    filterChips: [{ id: "all", label: "すべて", count: 1, selected: true }],
    jobs: [
      {
        jobId: 41,
        title: "input.json",
        jobState: jobOverride.jobState ?? "Paused",
        stateLabel: jobOverride.stateLabel ?? "中断中",
        stateDescription: jobOverride.stateDescription ?? "中断中",
        stateTone: "warning",
        inputSourceLabel: "input.json",
        sourcePath: "/mods/input.json",
        currentPhase: jobOverride.currentPhase ?? "body_translation",
        currentPhaseLabel: jobOverride.currentPhaseLabel ?? "本文翻訳",
        progressLabel: "20%",
        lastUpdatedLabel: "now",
        canOpenPhase: jobOverride.canOpenPhase ?? true,
        openBlockedReason: jobOverride.openBlockedReason ?? null,
        openBlockedReasonText: jobOverride.openBlockedReasonText ?? "",
        jobRunTarget: cardJobRunTarget,
        isSelected: true,
        stopOperation: {
          kind: "stop",
          label: "停止",
          helperText: "停止できます。",
          enabled: false,
          reasonText: "実行中ではありません。",
          busy: false
        },
        resumeOperation: {
          kind: "resume",
          label: "再開",
          helperText: "再開できます。",
          enabled: true,
          reasonText: "",
          busy: false
        },
        deleteOperation: {
          kind: "delete",
          label: "削除",
          helperText: "削除できます。",
          enabled: true,
          reasonText: "",
          busy: false
        },
        ...jobOverride
      }
    ],
    feedback,
    selectedJob: {
      jobId: 41,
      jobIdLabel: "ジョブ #41",
      jobState: "Paused",
      stateLabel: "中断中",
      stateDescription: "中断中",
      stateTone: "warning",
      inputSourceLabel: "input.json",
      inputSourceKindLabel: "xEdit 抽出データ",
      sourcePath: "/mods/input.json",
      pluginName: "Skyrim.esm",
      extractedJsonLabel: "抽出データ #1",
      currentPhase: "body_translation",
      currentPhaseLabel: "本文翻訳",
      progressLabel: "20%",
      lastUpdatedLabel: "now",
      cacheStateLabel: "利用可能",
      providerLabel: "openai",
      modelLabel: "gpt-5",
      executionModeLabel: "batch",
      credentialStateLabel: "設定済み",
      stopOperation: {
        kind: "stop",
        label: "停止",
        helperText: "停止できます。",
        enabled: false,
        reasonText: "実行中ではありません。",
        busy: false
      },
      resumeOperation: {
        kind: "resume",
        label: "再開",
        helperText: "再開できます。",
        enabled: true,
        reasonText: "",
        busy: false
      },
      deleteOperation: {
        kind: "delete",
        label: "削除",
        helperText: "削除できます。",
        enabled: true,
        reasonText: "",
        busy: false
      },
      resumeBlockedReasons: [],
      warnings: [],
      deleteImpactLines: []
    },
    deleteConfirmation: null,
    jobRunTarget
  }
}

class TranslationJobManagementPageControllerFake
  implements TranslationJobManagementScreenControllerContract
{
  readonly requestResume = vi.fn(() => {
    this.viewModel = this.nextViewModel
    for (const listener of this.listeners) {
      listener(this.viewModel)
    }
    return Promise.resolve()
  })

  private viewModel: TranslationJobManagementScreenViewModel
  private readonly listeners =
    new Set<TranslationJobManagementScreenViewModelListener>()

  constructor(
    private readonly nextViewModel: TranslationJobManagementScreenViewModel,
    initialViewModel: TranslationJobManagementScreenViewModel = createViewModel()
  ) {
    this.viewModel = initialViewModel
  }

  async mount(): Promise<void> {}
  dispose(): void {}
  subscribe(
    listener: TranslationJobManagementScreenViewModelListener
  ): () => void {
    this.listeners.add(listener)
    return () => this.listeners.delete(listener)
  }
  getViewModel(): TranslationJobManagementScreenViewModel {
    return this.viewModel
  }
  async reload(): Promise<void> {}
  readonly selectJob = vi.fn(() => Promise.resolve())
  setFilter(): void {}
  setSearchQuery(): void {}
  async requestStop(): Promise<void> {}
  openDeleteConfirmation(): void {}
  closeDeleteConfirmation(): void {}
  async deleteSelectedJob(): Promise<void> {}
}

describe("TranslationJobManagementPage", () => {
  test.each([
    ["ready", "Ready"],
    ["running", "Running"],
    ["failed", "Failed"]
  ] satisfies Array<[string, TranslationJobManagementJobState]>)(
    "%s job は現在の翻訳段階を開く",
    async (_caseName, jobState) => {
      const user = userEvent.setup()
      const onOpenJobRun = vi.fn()
      const target = {
        ...resumeTarget,
        stateLabel: jobState,
        stateDescription: jobState
      }
      const controller = new TranslationJobManagementPageControllerFake(
        createViewModel(),
        createViewModel(null, null, {
          jobState,
          stateLabel: jobState,
          stateDescription: jobState,
          canOpenPhase: true,
          jobRunTarget: target
        })
      )

      render(TranslationJobManagementPage, {
        props: {
          createController: () => controller,
          onOpenJobRun
        }
      })

      await user.click(
        screen.getByRole("button", { name: "現在の翻訳段階へ進む" })
      )

      expect(onOpenJobRun).toHaveBeenCalledWith(target)
      expect(controller.selectJob).toHaveBeenCalledWith(41)
    }
  )

  test("action 不可の job は shell を開かず選択だけ行う", async () => {
    const user = userEvent.setup()
    const onOpenJobRun = vi.fn()
    const controller = new TranslationJobManagementPageControllerFake(
      createViewModel(),
      createViewModel(null, null, {
        canOpenPhase: false,
        openBlockedReasonText: "現在の翻訳段階へ進む前に再開してください。",
        jobRunTarget: null
      })
    )

    render(TranslationJobManagementPage, {
      props: {
        createController: () => controller,
        onOpenJobRun
      }
    })

    const jobCard = screen.getByTestId("translation-job-management-job-card")
    const button = screen.getByRole("button", {
      name: "現在の翻訳段階へ進む"
    })

    expect(jobCard).toHaveAttribute("data-open-phase-state", "blocked")
    expect(button).toBeDisabled()
    expect(button).toHaveAttribute("data-current-phase-action-state", "blocked")

    await user.click(
      screen.getByTestId("translation-job-management-job-selection-region")
    )

    expect(controller.selectJob).toHaveBeenCalledWith(41)
    expect(onOpenJobRun).not.toHaveBeenCalled()
  })

  test("action 可能でも shell target がない job は shell を開かない", async () => {
    const user = userEvent.setup()
    const onOpenJobRun = vi.fn()
    const controller = new TranslationJobManagementPageControllerFake(
      createViewModel(),
      createViewModel(null, null, {
        canOpenPhase: true,
        jobRunTarget: null
      })
    )

    render(TranslationJobManagementPage, {
      props: {
        createController: () => controller,
        onOpenJobRun
      }
    })

    const jobCard = screen.getByTestId("translation-job-management-job-card")
    const button = screen.getByRole("button", {
      name: "現在の翻訳段階へ進む"
    })

    expect(jobCard).toHaveAttribute("data-open-phase-state", "missing-target")
    expect(button).toBeDisabled()
    expect(button).toHaveAttribute(
      "data-current-phase-action-state",
      "missing-target"
    )

    await user.click(
      screen.getByTestId("translation-job-management-job-selection-region")
    )

    expect(controller.selectJob).toHaveBeenCalledWith(41)
    expect(onOpenJobRun).not.toHaveBeenCalled()
  })

  test("resume success feedback と job run target がある時だけ shell を開く", async () => {
    vi.useFakeTimers()
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime })
    try {
      const onOpenJobRun = vi.fn()
      const controller = new TranslationJobManagementPageControllerFake(
        createViewModel(
          {
            title: "再開結果を更新しました",
            message: "再開を受け付けました",
            tone: "success",
            category: "resume_success"
          },
          resumeTarget
        )
      )

      render(TranslationJobManagementPage, {
        props: {
          createController: () => controller,
          onOpenJobRun
        }
      })

      await user.click(
        screen.getByTestId("translation-job-management-resume-button")
      )

      await waitFor(() =>
        expect(
          screen.getByTestId("translation-job-management-feedback-notification")
        ).toHaveTextContent("再開を受け付けました")
      )
      expect(onOpenJobRun).not.toHaveBeenCalled()

      await vi.advanceTimersByTimeAsync(250)

      await waitFor(() =>
        expect(onOpenJobRun).toHaveBeenCalledWith(resumeTarget)
      )
    } finally {
      vi.useRealTimers()
    }
  })

  test("resume warning feedback では job run target があっても shell を開かない", async () => {
    const user = userEvent.setup()
    const onOpenJobRun = vi.fn()
    const controller = new TranslationJobManagementPageControllerFake(
      createViewModel(
        {
          title: "再開結果を更新しました",
          message: "再開できません",
          tone: "warning",
          category: "resume_failed"
        },
        resumeTarget
      )
    )

    render(TranslationJobManagementPage, {
      props: {
        createController: () => controller,
        onOpenJobRun
      }
    })

    await user.click(
      screen.getByTestId("translation-job-management-resume-button")
    )

    await waitFor(() => expect(controller.requestResume).toHaveBeenCalledTimes(1))
    expect(onOpenJobRun).not.toHaveBeenCalled()
  })
})
