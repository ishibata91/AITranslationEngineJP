import type {
  ProviderSettingsScreenControllerContract,
  ProviderSettingsScreenViewModelListener
} from "@application/contract/provider-settings"
import type {
  ProviderSettingsScreenState,
  ProviderSettingsScreenViewModel
} from "@application/gateway-contract/provider-settings"

interface ProviderSettingsStoreLike {
  subscribe(listener: (state: ProviderSettingsScreenState) => void): () => void
  snapshot(): ProviderSettingsScreenState
}

interface ProviderSettingsPresenterLike {
  toViewModel(
    state: ProviderSettingsScreenState,
    isGatewayConnected: boolean
  ): ProviderSettingsScreenViewModel
}

interface ProviderSettingsUseCaseLike {
  load(): Promise<void>
  selectProvider(providerId: "gemini" | "xai" | "lm_studio"): void
  openApiKeyPanel(): void
  closeApiKeyPanel(): void
  updateEndpoint(nextValue: string): void
  saveSettings(readCredentialInput: () => string): Promise<void>
  resetSettings(): Promise<void>
  validateConnection(): Promise<void>
}

interface ProviderSettingsScreenControllerDependencies {
  isGatewayConnected: boolean
  store: ProviderSettingsStoreLike
  presenter: ProviderSettingsPresenterLike
  useCase: ProviderSettingsUseCaseLike
}

export class ProviderSettingsScreenController implements ProviderSettingsScreenControllerContract {
  private credentialInputValue = ""

  constructor(
    private readonly dependencies: ProviderSettingsScreenControllerDependencies
  ) {}

  mount(): Promise<void> {
    return this.dependencies.useCase.load()
  }

  dispose(): void {
    return
  }

  subscribe(listener: ProviderSettingsScreenViewModelListener): () => void {
    return this.dependencies.store.subscribe((state) => {
      listener(
        this.dependencies.presenter.toViewModel(
          state,
          this.dependencies.isGatewayConnected
        )
      )
    })
  }

  getViewModel(): ProviderSettingsScreenViewModel {
    return this.dependencies.presenter.toViewModel(
      this.dependencies.store.snapshot(),
      this.dependencies.isGatewayConnected
    )
  }

  selectProvider(providerId: string): void {
    if (
      providerId === "gemini" ||
      providerId === "xai" ||
      providerId === "lm_studio"
    ) {
      this.dependencies.useCase.selectProvider(providerId)
    }
  }

  updateEndpoint(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    this.dependencies.useCase.updateEndpoint(target.value)
  }

  openApiKeyPanel(): void {
    this.dependencies.useCase.openApiKeyPanel()
  }

  closeApiKeyPanel(): void {
    this.dependencies.useCase.closeApiKeyPanel()
  }

  updateCredentialInput(nextValue: string): void {
    this.credentialInputValue = nextValue
  }

  clearCredentialInput(): void {
    this.credentialInputValue = ""
  }

  async saveSettings(): Promise<void> {
    await this.dependencies.useCase.saveSettings(
      () => this.credentialInputValue
    )
    this.clearCredentialInput()
  }

  async resetSettings(): Promise<void> {
    this.clearCredentialInput()
    await this.dependencies.useCase.resetSettings()
  }

  async validateConnection(): Promise<void> {
    await this.dependencies.useCase.validateConnection()
  }
}
