import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  ProviderSettingsScreenControllerContract,
  ProviderSettingsScreenViewModelListener
} from "@application/contract/provider-settings"
import { ProviderSettingsPresenter } from "@application/presenter/provider-settings"
import { ProviderSettingsStore } from "@application/store/provider-settings"
import { ProviderSettingsUseCase } from "@application/usecase/provider-settings"
import ProviderSettingsPage from "@ui/screens/provider-settings/ProviderSettingsPage.svelte"

class ProviderSettingsScreenControllerFake
  implements ProviderSettingsScreenControllerContract
{
  private readonly store = new ProviderSettingsStore()

  private readonly presenter = new ProviderSettingsPresenter()

  private readonly useCase = new ProviderSettingsUseCase(null, this.store)

  readonly mount = vi.fn(async () => {
    await this.useCase.load()
  })

  readonly dispose = vi.fn(() => {})

  subscribe(listener: ProviderSettingsScreenViewModelListener): () => void {
    return this.store.subscribe((state) => {
      listener(this.presenter.toViewModel(state, false))
    })
  }

  getViewModel() {
    return this.presenter.toViewModel(this.store.snapshot(), false)
  }

  selectProvider(providerId: string): void {
    if (providerId === "gemini" || providerId === "xai" || providerId === "lm_studio") {
      this.useCase.selectProvider(providerId)
    }
  }

  updateEndpoint(event: Event): void {
    const target = event.currentTarget
    if (target instanceof HTMLInputElement) {
      this.useCase.updateEndpoint(target.value)
    }
  }

  openApiKeyPanel(): void {
    this.useCase.openApiKeyPanel()
  }

  closeApiKeyPanel(): void {
    this.useCase.closeApiKeyPanel()
  }

  updateCredentialInput(nextValue: string): void {
    this.credentialInputValue = nextValue
  }

  clearCredentialInput(): void {
    this.credentialInputValue = ""
  }

  private credentialInputValue = ""

  async saveSettings(): Promise<void> {
    await this.useCase.saveSettings(() => this.credentialInputValue)
    this.clearCredentialInput()
  }

  async resetSettings(): Promise<void> {
    this.clearCredentialInput()
    await this.useCase.resetSettings()
  }

  async validateConnection(): Promise<void> {
    await this.useCase.validateConnection()
  }
}

describe("ProviderSettingsPage", () => {
  test("必要な UI だけを表示し、provider list は 3 件に限定する", async () => {
    const controller = new ProviderSettingsScreenControllerFake()
    render(ProviderSettingsPage, {
      props: {
        createController: () => controller
      }
    })

    await waitFor(() => {
      expect(
        screen.getByRole("heading", { level: 2, name: "AIサービス設定" })
      ).toBeInTheDocument()
    })

    expect(screen.getAllByText("Gemini").length).toBeGreaterThan(0)
    expect(screen.getAllByText("xAI").length).toBeGreaterThan(0)
    expect(screen.getAllByText("LM Studio").length).toBeGreaterThan(0)
    expect(screen.queryByText("OpenAI")).not.toBeInTheDocument()
    expect(screen.getByLabelText("エンドポイント")).toBeInTheDocument()
    expect(screen.getByText("APIキー状態")).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "接続を確認" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "リセット" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "設定を保存" })
    ).toBeInTheDocument()
    expect(screen.queryByText("Batch API")).not.toBeInTheDocument()
    expect(screen.queryByText("処理方法")).not.toBeInTheDocument()
    expect(screen.queryByText("モデル")).not.toBeInTheDocument()
  })

  test("APIキーは必要 provider にだけ表示し、保存後に DOM へ残さない", async () => {
    const user = userEvent.setup()
    const controller = new ProviderSettingsScreenControllerFake()
    render(ProviderSettingsPage, {
      props: {
        createController: () => controller
      }
    })

    await waitFor(() => {
      expect(screen.getAllByText("Gemini").length).toBeGreaterThan(0)
    })

    await user.click(screen.getAllByRole("button", { name: "設定" })[0])
    const apiKeyInput = screen.getByLabelText("APIキー")
    await user.type(apiKeyInput, "top-secret-value")
    await user.click(screen.getByRole("button", { name: "保存" }))

    await waitFor(() => {
      expect(screen.queryByDisplayValue("top-secret-value")).not.toBeInTheDocument()
    })

    await user.click(screen.getByText("LM Studio"))
    expect(
      screen.queryByRole("button", { name: "設定" })
    ).not.toBeInTheDocument()
  })
})
