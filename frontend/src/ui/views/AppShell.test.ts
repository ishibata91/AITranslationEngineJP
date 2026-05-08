import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import { createShellState } from "@ui/stores/shell-state"
import AppShell from "@ui/views/AppShell.svelte"
import type { TranslationJobManagementScreenViewModel } from "@application/contract/translation-job-management/translation-job-management-screen-types"

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
            pageLead: "AIサービスごとのエンドポイントと APIキー状態を管理します。",
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
})
