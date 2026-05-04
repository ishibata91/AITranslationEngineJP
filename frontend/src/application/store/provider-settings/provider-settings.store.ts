import type {
  ProviderSettingsProviderState,
  ProviderSettingsScreenState
} from "@application/gateway-contract/provider-settings"

type Listener = (state: ProviderSettingsScreenState) => void

function cloneProvider(
  provider: ProviderSettingsProviderState
): ProviderSettingsProviderState {
  return { ...provider }
}

function createInitialState(): ProviderSettingsScreenState {
  return {
    phase: "idle",
    providers: [],
    selectedProviderId: null,
    apiKeyPanelOpen: false,
    saveNotice: "",
    errorMessage: ""
  }
}

export class ProviderSettingsStore {
  private state: ProviderSettingsScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): ProviderSettingsScreenState {
    return {
      ...this.state,
      providers: this.state.providers.map(cloneProvider)
    }
  }

  update(mutator: (draft: ProviderSettingsScreenState) => void): void {
    const nextState = this.snapshot()
    mutator(nextState)
    this.state = nextState
    this.emit()
  }

  private emit(): void {
    const snapshot = this.snapshot()
    for (const listener of this.listeners) {
      listener(snapshot)
    }
  }
}
