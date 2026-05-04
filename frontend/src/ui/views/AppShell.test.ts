import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import { createShellState } from "@ui/stores/shell-state"
import AppShell from "@ui/views/AppShell.svelte"

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
})
