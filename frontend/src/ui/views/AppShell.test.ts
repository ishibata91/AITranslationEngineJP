import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  PersonaGenerationPhaseScreenControllerContract,
  PersonaGenerationPhaseScreenViewModel
} from "@application/contract/persona-generation-phase"
import type {
  TermTranslationPhaseScreenControllerContract,
  TermTranslationPhaseScreenViewModel
} from "@application/contract/term-translation-phase"
import type {
  TranslationJobManagementScreenControllerContract,
  TranslationJobManagementScreenViewModelListener
} from "@application/contract/translation-job-management/translation-job-management-screen-contract"
import { createShellState } from "@ui/stores/shell-state"
import AppShell from "@ui/views/AppShell.svelte"
import type {
  TranslationJobManagementJobRunTarget,
  TranslationJobManagementScreenViewModel
} from "@application/contract/translation-job-management/translation-job-management-screen-types"

function createTranslationJobManagementViewModel(): TranslationJobManagementScreenViewModel {
  return {
    gatewayStatus: "接続済み",
    pageTitle: "未完了ジョブ一覧",
    pageLead:
      "新規翻訳を開始するか、未完了ジョブを選んで現在の翻訳段階から再開します。",
    headerCountLabel: "0 件を表示",
    listEmptyTitle: "管理対象がありません",
    listEmptyDescription: "未完了ジョブはありません。",
    listErrorTitle: "一覧を読み込めません",
    listErrorDescription: "未完了ジョブの一覧取得に失敗しました。",
    detailPlaceholderTitle: "job を選択してください",
    detailPlaceholderDescription: "一覧から 1 件選びます。",
    phase: "empty",
    detailPhase: "idle",
    isReloading: false,
    searchQuery: "",
    filterChips: [{ id: "all", label: "すべて", count: 0, selected: true }],
    jobs: [],
    feedback: null,
    selectedJob: null,
    deleteConfirmation: null,
    jobRunTarget: null
  }
}

const jobOneRunTarget: TranslationJobManagementJobRunTarget = {
  jobId: 1,
  stateLabel: "Ready",
  stateDescription: "実行待ち",
  currentPhase: "term_translation",
  currentPhaseLabel: "単語翻訳",
  progressLabel: "0% / 未開始",
  inputSourceLabel: "jobID1.esp",
  sourcePath: "/tmp/jobID1.esp"
}

function createTranslationJobManagementViewModelWithJob(
  jobRunTarget: TranslationJobManagementJobRunTarget | null = null
): TranslationJobManagementScreenViewModel {
  return {
    ...createTranslationJobManagementViewModel(),
    headerCountLabel: "1 件を表示",
    phase: "ready",
    jobs: [
      {
        jobId: 1,
        title: "jobID1",
        jobState: "Ready",
        stateLabel: "Ready",
        stateDescription: "実行待ち",
        stateTone: "info",
        inputSourceLabel: "jobID1.esp",
        sourcePath: "/tmp/jobID1.esp",
        currentPhase: "term_translation",
        currentPhaseLabel: "単語翻訳",
        progressLabel: "0% / 未開始",
        lastUpdatedLabel: "たった今",
        canOpenPhase: true,
        openBlockedReason: null,
        openBlockedReasonText: "",
        jobRunTarget: jobOneRunTarget,
        isSelected: false,
        stopOperation: {
          kind: "stop",
          label: "停止",
          helperText: "実行中ではありません。",
          enabled: false,
          reasonText: "停止できません。実行中ではありません。",
          busy: false
        },
        resumeOperation: {
          kind: "resume",
          label: "再開",
          helperText: "まだ実行されていません。",
          enabled: false,
          reasonText: "再開できません。実行前のジョブです。",
          busy: false
        },
        deleteOperation: {
          kind: "delete",
          label: "削除",
          helperText: "ジョブの DB 情報だけを削除します。",
          enabled: true,
          reasonText: "",
          busy: false
        }
      }
    ],
    jobRunTarget
  }
}

function createTermTranslationPhaseViewModel(): TermTranslationPhaseScreenViewModel {
  return {
    jobId: null,
    phase: "ready",
    summary: null,
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    gatewayStatus: "接続済み",
    viewState: "blocked",
    isLoading: false,
    isRefreshing: false,
    isSubmitting: false,
    hasJobSelection: true,
    currentPhaseLabel: "単語翻訳",
    phaseStateLabel: "開始待ち",
    statusTitle: "単語翻訳を開始できます",
    statusText: "現在のジョブを対象に単語翻訳を確認します。",
    progressPercent: 0,
    progressLabel: "0%",
    progressDetail: "未開始",
    startedAtLabel: "-",
    finishedAtLabel: "-",
    totalTermCountLabel: "0 件",
    dictionaryHitCountLabel: "0 件",
    aiTargetCountLabel: "0 件",
    confirmedCountLabel: "0 件",
    jobDictionaryAppliedCountLabel: "0 件",
    replacementTargetCountLabel: "0 件",
    unmatchedCountLabel: "0 件",
    providerLabel: "-",
    modelLabel: "-",
    executionModeLabel: "-",
    credentialRefLabel: "-",
    snapshotLabel: "-",
    errorKindLabel: "-",
    errorReasonLabel: "-",
    retryableLabel: "-",
    nextPhaseStatusLabel: "開始不可",
    nextPhaseBlockedReason: "",
    providerSkippedLabel: "-",
    actionCards: [],
    lastErrorSummary: null,
    actionEnablement: null,
    latestProgressSummary: null,
    latestResultSummary: null,
    latestExecutionSummary: null,
    latestErrorKind: null
  }
}

function createTermTranslationPhaseController(): TermTranslationPhaseScreenControllerContract {
  const viewModel = createTermTranslationPhaseViewModel()
  return {
    mount: vi.fn(async () => {}),
    dispose: vi.fn(() => {}),
    subscribe: vi.fn(() => () => {}),
    getViewModel: vi.fn(() => viewModel),
    setJobId: vi.fn(async () => {}),
    refresh: vi.fn(async () => {}),
    startPhase: vi.fn(async () => {}),
    pausePhase: vi.fn(async () => {}),
    resumePhase: vi.fn(async () => {}),
    retryPhase: vi.fn(async () => {})
  }
}

function createPersonaGenerationPhaseController(): PersonaGenerationPhaseScreenControllerContract {
  const viewModel = {
    bodyReadiness: null,
    latestBodyReadiness: null,
    bodyReadinessBlockedReason: ""
  } as PersonaGenerationPhaseScreenViewModel
  return {
    mount: vi.fn(async () => {}),
    dispose: vi.fn(() => {}),
    subscribe: vi.fn(() => () => {}),
    getViewModel: vi.fn(() => viewModel),
    setJobId: vi.fn(async () => {}),
    refresh: vi.fn(async () => {}),
    startPhase: vi.fn(async () => {}),
    pausePhase: vi.fn(async () => {}),
    resumePhase: vi.fn(async () => {}),
    retryPhase: vi.fn(async () => {}),
    cancelPhase: vi.fn(async () => {}),
    checkBodyReadiness: vi.fn(async () => {}),
    startBodyPhase: vi.fn(async () => {})
  }
}

function createTranslationJobManagementScenarioController(): TranslationJobManagementScreenControllerContract {
  let listener: TranslationJobManagementScreenViewModelListener | null = null
  return {
    mount: vi.fn(async () => {}),
    dispose: vi.fn(() => {}),
    subscribe: vi.fn(
      (nextListener: TranslationJobManagementScreenViewModelListener) => {
        listener = nextListener
        return () => {
          listener = null
        }
      }
    ),
    getViewModel: vi.fn(() => createTranslationJobManagementViewModelWithJob()),
    reload: vi.fn(async () => {}),
    selectJob: vi.fn((() =>
      Promise.resolve().then(() => {
        listener?.(
          createTranslationJobManagementViewModelWithJob(jobOneRunTarget)
        )
      })) as TranslationJobManagementScreenControllerContract["selectJob"]),
    setFilter: vi.fn(),
    setSearchQuery: vi.fn(),
    requestStop: vi.fn(async () => {}),
    requestResume: vi.fn(async () => {}),
    openDeleteConfirmation: vi.fn(),
    closeDeleteConfirmation: vi.fn(),
    deleteSelectedJob: vi.fn(async () => {})
  }
}

describe("AppShell", () => {
  test("dashboard entry から AIサービス設定 へ到達できる", async () => {
    window.history.replaceState(null, "", "#dashboard")

    const user = userEvent.setup()
    const shellState = createShellState()
    render(AppShell, {
      props: {
        defaultRouteId: shellState.defaultRouteId,
        routes: shellState.routes,
        defaultTranslationManagementViewId:
          shellState.defaultTranslationManagementViewId,
        translationManagementViews: shellState.translationManagementViews,
        createBodyTranslationPhaseScreenController: null,
        createMasterDictionaryScreenController: null,
        createMasterPersonaScreenController: null,
        createPersonaGenerationPhaseScreenController: null,
        createProviderSettingsScreenController: vi.fn(() => ({
          mount: vi.fn(async () => {}),
          dispose: vi.fn(() => {}),
          subscribe: vi.fn(() => () => {}),
          getViewModel: vi.fn(() => ({
            gatewayStatus: "未接続",
            pageTitle: "AIサービス設定",
            pageLead:
              "AIサービスごとのエンドポイントと APIキー状態を管理します。",
            providerCountLabel: "3 件の AIサービスを管理します。",
            phaseLabel: "待機中",
            saveNotice: "",
            errorMessage: "",
            providerList: [],
            selectedProvider: null
          })),
          selectProvider: vi.fn(),
          updateEndpoint: vi.fn(),
          openApiKeyPanel: vi.fn(),
          closeApiKeyPanel: vi.fn(),
          saveSettings: vi.fn(async () => {}),
          resetSettings: vi.fn(async () => {}),
          validateConnection: vi.fn(async () => {})
        })),
        createTermTranslationPhaseScreenController: null,
        createTranslationJobSetupScreenController: null,
        createTranslationOutputArtifactScreenController: null,
        createTranslationInputScreenController: null
      }
    })

    await user.click(screen.getAllByText("AIサービス設定")[0])

    await waitFor(() => {
      expect(
        screen.getByRole("heading", { level: 1, name: "AIサービス設定" })
      ).toBeInTheDocument()
    })
  })

  test("翻訳管理は未完了ジョブ一覧を初期表示し旧実行タブを出さない", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const shellState = createShellState()
    render(AppShell, {
      props: {
        defaultRouteId: shellState.defaultRouteId,
        routes: shellState.routes,
        defaultTranslationManagementViewId:
          shellState.defaultTranslationManagementViewId,
        translationManagementViews: shellState.translationManagementViews,
        createBodyTranslationPhaseScreenController: null,
        createMasterDictionaryScreenController: null,
        createMasterPersonaScreenController: null,
        createPersonaGenerationPhaseScreenController: null,
        createProviderSettingsScreenController: null,
        createTermTranslationPhaseScreenController: null,
        createTranslationJobManagementScreenController: vi.fn(() => ({
          mount: vi.fn(async () => {}),
          dispose: vi.fn(() => {}),
          subscribe: vi.fn(() => () => {}),
          getViewModel: vi.fn(createTranslationJobManagementViewModel),
          reload: vi.fn(async () => {}),
          selectJob: vi.fn(async () => {}),
          setFilter: vi.fn(),
          setSearchQuery: vi.fn(),
          requestStop: vi.fn(async () => {}),
          requestResume: vi.fn(async () => {}),
          openDeleteConfirmation: vi.fn(),
          closeDeleteConfirmation: vi.fn(),
          deleteSelectedJob: vi.fn(async () => {})
        })),
        createTranslationJobSetupScreenController: null,
        createTranslationOutputArtifactScreenController: null,
        createTranslationInputScreenController: null
      }
    })

    await waitFor(() => {
      expect(
        screen.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
      ).toBeInTheDocument()
    })
    expect(
      screen.queryByRole("tab", { name: /Job Run|実行/ })
    ).not.toBeInTheDocument()
  })

  test("翻訳段階ページの hash 直リンクは未完了一覧へ戻す", async () => {
    window.history.replaceState(
      null,
      "",
      "#translation-management/job-run?jobId=10&phase=本文翻訳"
    )

    const shellState = createShellState()
    render(AppShell, {
      props: {
        defaultRouteId: shellState.defaultRouteId,
        routes: shellState.routes,
        defaultTranslationManagementViewId:
          shellState.defaultTranslationManagementViewId,
        translationManagementViews: shellState.translationManagementViews,
        createBodyTranslationPhaseScreenController: null,
        createMasterDictionaryScreenController: null,
        createMasterPersonaScreenController: null,
        createPersonaGenerationPhaseScreenController: null,
        createProviderSettingsScreenController: null,
        createTermTranslationPhaseScreenController: null,
        createTranslationJobManagementScreenController: vi.fn(() => ({
          mount: vi.fn(async () => {}),
          dispose: vi.fn(() => {}),
          subscribe: vi.fn(() => () => {}),
          getViewModel: vi.fn(createTranslationJobManagementViewModel),
          reload: vi.fn(async () => {}),
          selectJob: vi.fn(async () => {}),
          setFilter: vi.fn(),
          setSearchQuery: vi.fn(),
          requestStop: vi.fn(async () => {}),
          requestResume: vi.fn(async () => {}),
          openDeleteConfirmation: vi.fn(),
          closeDeleteConfirmation: vi.fn(),
          deleteSelectedJob: vi.fn(async () => {})
        })),
        createTranslationJobSetupScreenController: null,
        createTranslationOutputArtifactScreenController: null,
        createTranslationInputScreenController: null
      }
    })

    await waitFor(() => {
      expect(window.location.hash).toBe("#translation-management")
    })
    expect(
      screen.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
    ).toBeInTheDocument()
  })

  test("未完了ジョブ一覧で現在の翻訳段階へ進む操作は初回と再実行で単語翻訳を表示する", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const user = userEvent.setup()
    const shellState = createShellState()
    const translationJobManagementController =
      createTranslationJobManagementScenarioController()
    render(AppShell, {
      props: {
        defaultRouteId: shellState.defaultRouteId,
        routes: shellState.routes,
        defaultTranslationManagementViewId:
          shellState.defaultTranslationManagementViewId,
        translationManagementViews: shellState.translationManagementViews,
        createBodyTranslationPhaseScreenController: null,
        createMasterDictionaryScreenController: null,
        createMasterPersonaScreenController: null,
        createPersonaGenerationPhaseScreenController: vi.fn(
          createPersonaGenerationPhaseController
        ),
        createProviderSettingsScreenController: null,
        createTermTranslationPhaseScreenController: vi.fn(
          createTermTranslationPhaseController
        ),
        createTranslationJobManagementScreenController: vi.fn(
          () => translationJobManagementController
        ),
        createTranslationJobSetupScreenController: null,
        createTranslationOutputArtifactScreenController: null,
        createTranslationInputScreenController: null
      }
    })

    await user.click(
      screen.getByRole("button", { name: "現在の翻訳段階へ進む" })
    )

    await waitFor(() => {
      expect(window.location.hash).toBe("#translation-management/job-run")
    })
    expect(
      screen.getByRole("heading", { level: 3, name: "ジョブ #1" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("heading", { level: 2, name: "単語翻訳" })
    ).toBeInTheDocument()
    expect(
      screen.queryByText("未完了ジョブ一覧でジョブを選んでください")
    ).not.toBeInTheDocument()

    await user.click(screen.getByRole("button", { name: "未完了一覧へ戻る" }))
    await user.click(
      screen.getByRole("button", { name: "現在の翻訳段階へ進む" })
    )

    await waitFor(() => {
      expect(window.location.hash).toBe("#translation-management/job-run")
    })
    expect(
      screen.getByRole("heading", { level: 3, name: "ジョブ #1" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("heading", { level: 2, name: "単語翻訳" })
    ).toBeInTheDocument()
  })
})
