import type {
  ProviderSettingsGatewayContract,
  ProviderSettingsErrorKind,
  ProviderSettingsProviderId,
  ProviderSettingsProviderState,
  ProviderSettingsScreenState,
  ProviderSettingsSummary,
  ValidateProviderSettingsResponse
} from "@application/gateway-contract/provider-settings"

interface ProviderSettingsStoreLike {
  snapshot(): ProviderSettingsScreenState
  update(mutator: (draft: ProviderSettingsScreenState) => void): void
}

function providerLabel(providerId: ProviderSettingsProviderId): string {
  switch (providerId) {
    case "gemini":
      return "Gemini"
    case "xai":
      return "xAI"
    case "lm_studio":
      return "LM Studio"
  }
}

function defaultEndpoint(providerId: ProviderSettingsProviderId): string {
  switch (providerId) {
    case "gemini":
      return "https://generativelanguage.googleapis.com"
    case "xai":
      return "https://api.x.ai/v1"
    case "lm_studio":
      return "http://127.0.0.1:1234/v1"
  }
}

function requiresCredential(providerId: ProviderSettingsProviderId): boolean {
  return providerId !== "lm_studio"
}

function isSaveValidationFailure(error: unknown): boolean {
  if (error instanceof Error) {
    return error.message.includes("validation_failed")
  }

  return typeof error === "string" && error.includes("validation_failed")
}

function providerSettingsSaveErrorMessage(error: unknown): string {
  if (isSaveValidationFailure(error)) {
    return "入力内容が不正です。エンドポイントを確認してください。"
  }

  return "設定の保存に失敗しました。"
}

function buildSavedState(provider: {
  endpoint: string
  credentialConfigured: boolean
  requiresCredential: boolean
}): ProviderSettingsProviderState["savedState"] {
  const hasEndpoint = provider.endpoint.trim().length > 0
  if (
    !hasEndpoint &&
    (!provider.requiresCredential || !provider.credentialConfigured)
  ) {
    return "not_saved"
  }

  if (!provider.requiresCredential) {
    return hasEndpoint ? "configured" : "not_saved"
  }

  return hasEndpoint && provider.credentialConfigured ? "configured" : "partial"
}

function cloneProviderSummary(
  provider: ProviderSettingsSummary
): ProviderSettingsProviderState {
  const endpoint = provider.endpoint ?? ""
  return {
    providerId: provider.providerId,
    label: provider.label,
    endpointDraft: endpoint,
    persistedEndpoint: endpoint,
    credentialState: provider.credentialState,
    credentialReferenceId: provider.credentialReferenceId,
    validationState: provider.validationState,
    savedState: provider.savedState,
    requestToken: provider.requestToken ?? ""
  }
}

function createDefaultProviders(): ProviderSettingsProviderState[] {
  return (["gemini", "xai", "lm_studio"] as const).map((providerId) => {
    const endpoint = defaultEndpoint(providerId)
    const credentialConfigured = providerId === "gemini"
    return {
      providerId,
      label: providerLabel(providerId),
      endpointDraft: endpoint,
      persistedEndpoint: endpoint,
      credentialState: requiresCredential(providerId)
        ? credentialConfigured
          ? "configured"
          : "missing"
        : "not_required",
      validationState:
        credentialConfigured || providerId === "lm_studio"
          ? "validated"
          : "not_validated",
      savedState: buildSavedState({
        endpoint,
        credentialConfigured,
        requiresCredential: requiresCredential(providerId)
      }),
      requestToken: `${providerId}-initial`,
      lastFailureKind:
        credentialConfigured || providerId === "lm_studio"
          ? undefined
          : "credential_missing"
    }
  })
}

function delay(milliseconds: number): Promise<void> {
  return new Promise((resolve) => {
    window.setTimeout(resolve, milliseconds)
  })
}

export class ProviderSettingsUseCase {
  private requestSequence = 0

  constructor(
    private readonly gateway: ProviderSettingsGatewayContract | null,
    private readonly store: ProviderSettingsStoreLike
  ) {}

  async load(): Promise<void> {
    this.store.update((draft) => {
      draft.phase = "loading"
      draft.errorMessage = ""
    })

    if (this.gateway) {
      const response = await this.gateway.ListProviderSettings({})
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.providers = response.providers
          .map(cloneProviderSummary)
          .filter(
            (provider) =>
              provider.providerId === "gemini" ||
              provider.providerId === "xai" ||
              provider.providerId === "lm_studio"
          )
        draft.selectedProviderId = draft.providers[0]?.providerId ?? null
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "ready"
      draft.providers = createDefaultProviders()
      draft.selectedProviderId = draft.providers[0]?.providerId ?? null
    })
  }

  selectProvider(providerId: ProviderSettingsProviderId): void {
    this.store.update((draft) => {
      draft.selectedProviderId = providerId
      draft.apiKeyPanelOpen = false
      draft.errorMessage = ""
    })
  }

  openApiKeyPanel(): void {
    this.store.update((draft) => {
      draft.apiKeyPanelOpen = true
      draft.errorMessage = ""
    })
  }

  closeApiKeyPanel(): void {
    this.store.update((draft) => {
      draft.apiKeyPanelOpen = false
    })
  }

  updateEndpoint(nextValue: string): void {
    const state = this.store.snapshot()
    const providerId = state.selectedProviderId
    if (!providerId) {
      return
    }

    this.store.update((draft) => {
      const provider = draft.providers.find(
        (candidate) => candidate.providerId === providerId
      )
      if (!provider) {
        return
      }

      provider.endpointDraft = nextValue
      provider.validationState = "not_validated"
      provider.lastFailureKind = "validation_stale"
      provider.requestToken = this.nextRequestToken(provider.providerId)
      draft.saveNotice = ""
      draft.errorMessage = ""
    })
  }

  async saveSettings(readCredentialInput: () => string): Promise<void> {
    const state = this.store.snapshot()
    const provider = this.selectedProvider(state)
    if (!provider) {
      return
    }

    const credentialInput = readCredentialInput().trim()
    const credentialInputPresent = credentialInput.length > 0
    const endpoint = provider.endpointDraft.trim()

    this.store.update((draft) => {
      draft.phase = "saving"
      draft.errorMessage = ""
      draft.saveNotice = ""
    })

    if (this.gateway) {
      try {
        const response = await this.gateway.SaveProviderSettings({
          providerId: provider.providerId,
          endpoint: endpoint || undefined,
          apiKeyInputPresent: credentialInputPresent,
          credentialInput: credentialInputPresent ? credentialInput : undefined
        })
        this.store.update((draft) => {
          const current = draft.providers.find(
            (candidate) =>
              candidate.providerId === response.provider.providerId
          )
          if (!current) {
            draft.phase = "ready"
            return
          }
          this.applySavedProvider(
            current,
            response.provider,
            credentialInputPresent
          )
          draft.phase = "ready"
          draft.apiKeyPanelOpen = false
          draft.saveNotice = "設定を保存しました。"
        })
      } catch (error) {
        this.store.update((draft) => {
          const current = draft.providers.find(
            (candidate) => candidate.providerId === provider.providerId
          )
          if (current && isSaveValidationFailure(error)) {
            this.applySaveFailure(current, "validation_failed")
          }
          draft.phase = "ready"
          draft.saveNotice = ""
          draft.errorMessage = providerSettingsSaveErrorMessage(error)
        })
      }
      return
    }

    await delay(10)

    this.store.update((draft) => {
      const current = draft.providers.find(
        (candidate) => candidate.providerId === provider.providerId
      )
      if (!current) {
        draft.phase = "ready"
        return
      }

      const credentialConfigured =
        current.credentialState === "configured" || credentialInputPresent
      current.endpointDraft = endpoint
      current.persistedEndpoint = endpoint
      current.credentialState = requiresCredential(current.providerId)
        ? credentialConfigured
          ? "configured"
          : "missing"
        : "not_required"
      current.savedState = buildSavedState({
        endpoint,
        credentialConfigured,
        requiresCredential: requiresCredential(current.providerId)
      })
      current.credentialReferenceId = credentialConfigured
        ? (current.credentialReferenceId ?? `${current.providerId}-local`)
        : undefined
      current.validationState = "not_validated"
      current.lastFailureKind = endpoint
        ? "validation_stale"
        : "endpoint_missing"
      current.requestToken = this.nextRequestToken(current.providerId)
      draft.phase = "ready"
      draft.apiKeyPanelOpen = false
      draft.saveNotice = "設定を保存しました。"
    })
  }

  async resetSettings(): Promise<void> {
    const state = this.store.snapshot()
    const provider = this.selectedProvider(state)
    if (!provider) {
      return
    }

    this.store.update((draft) => {
      draft.phase = "resetting"
      draft.errorMessage = ""
    })

    if (this.gateway) {
      const response = await this.gateway.ResetProviderSettings({
        providerId: provider.providerId
      })
      this.store.update((draft) => {
        const current = draft.providers.find(
          (candidate) => candidate.providerId === response.provider.providerId
        )
        if (!current) {
          draft.phase = "ready"
          return
        }
        this.applySavedProvider(current, response.provider, false)
        draft.phase = "ready"
        draft.apiKeyPanelOpen = false
        draft.saveNotice = "リセットしました。"
      })
      return
    }

    await delay(10)

    this.store.update((draft) => {
      const current = draft.providers.find(
        (candidate) => candidate.providerId === provider.providerId
      )
      if (!current) {
        draft.phase = "ready"
        return
      }
      current.endpointDraft = ""
      current.persistedEndpoint = ""
      current.credentialState = requiresCredential(current.providerId)
        ? "missing"
        : "not_required"
      current.credentialReferenceId = undefined
      current.validationState = "not_validated"
      current.savedState = "not_saved"
      current.lastFailureKind = "endpoint_missing"
      current.requestToken = this.nextRequestToken(current.providerId)
      draft.phase = "ready"
      draft.apiKeyPanelOpen = false
      draft.saveNotice = "リセットしました。"
    })
  }

  async validateConnection(): Promise<void> {
    const state = this.store.snapshot()
    const provider = this.selectedProvider(state)
    if (!provider) {
      return
    }

    if (
      requiresCredential(provider.providerId) &&
      provider.credentialState !== "configured"
    ) {
      this.store.update((draft) => {
        const current = draft.providers.find(
          (candidate) => candidate.providerId === provider.providerId
        )
        if (!current) {
          return
        }
        current.validationState = "not_validated"
        current.lastFailureKind = "credential_missing"
        current.requestToken = this.nextRequestToken(current.providerId)
        draft.errorMessage = "APIキーを設定してから接続確認してください。"
      })
      return
    }

    if (!provider.endpointDraft.trim()) {
      this.store.update((draft) => {
        const current = draft.providers.find(
          (candidate) => candidate.providerId === provider.providerId
        )
        if (!current) {
          return
        }
        current.validationState = "not_validated"
        current.lastFailureKind = "endpoint_missing"
        current.requestToken = this.nextRequestToken(current.providerId)
        draft.errorMessage =
          "エンドポイントを入力してから接続確認してください。"
      })
      return
    }

    const requestToken =
      provider.requestToken || this.nextRequestToken(provider.providerId)
    const endpoint = provider.endpointDraft.trim()

    this.store.update((draft) => {
      const current = draft.providers.find(
        (candidate) => candidate.providerId === provider.providerId
      )
      if (!current) {
        return
      }
      if (!current.requestToken) {
        current.requestToken = requestToken
      }
      current.validationState = "pending"
      current.lastFailureKind = undefined
      draft.phase = "validating"
      draft.errorMessage = ""
      draft.saveNotice = ""
    })

    const response = this.gateway
      ? await this.gateway.ValidateProviderSettings({
          providerId: provider.providerId,
          endpoint,
          credentialState: provider.credentialState,
          credentialReferenceId: provider.credentialReferenceId,
          requestToken
        })
      : await this.validateLocally(provider.providerId, requestToken, endpoint)

    this.store.update((draft) => {
      const current = draft.providers.find(
        (candidate) => candidate.providerId === response.providerId
      )
      if (!current) {
        draft.phase = "ready"
        return
      }

      if (current.requestToken !== response.requestToken) {
        draft.phase = "ready"
        return
      }

      current.validationState = response.validationState
      current.lastFailureKind = response.failureKind
      draft.phase = "ready"
    })
  }

  private async validateLocally(
    providerId: ProviderSettingsProviderId,
    requestToken: string,
    endpoint: string
  ): Promise<ValidateProviderSettingsResponse> {
    await delay(40)

    if (!endpoint.startsWith("http://") && !endpoint.startsWith("https://")) {
      return {
        providerId,
        validationState: "failed",
        requestToken,
        failureKind: "provider_unreachable"
      }
    }

    return {
      providerId,
      validationState: "validated",
      requestToken
    }
  }

  private selectedProvider(
    state: ProviderSettingsScreenState
  ): ProviderSettingsProviderState | null {
    return (
      state.providers.find(
        (provider) => provider.providerId === state.selectedProviderId
      ) ?? null
    )
  }

  private applySaveFailure(
    provider: ProviderSettingsProviderState,
    failureKind: ProviderSettingsErrorKind
  ): void {
    provider.validationState = "failed"
    provider.lastFailureKind = failureKind
    provider.requestToken = this.nextRequestToken(provider.providerId)
  }

  private applySavedProvider(
    current: ProviderSettingsProviderState,
    provider: ProviderSettingsSummary,
    apiKeyInputPresent: boolean
  ): void {
    current.endpointDraft = provider.endpoint ?? ""
    current.persistedEndpoint = provider.endpoint ?? ""
    current.credentialState =
      apiKeyInputPresent && current.credentialState !== "not_required"
        ? "configured"
        : provider.credentialState
    current.credentialReferenceId = provider.credentialReferenceId
    current.validationState = provider.validationState
    current.savedState = provider.savedState
    current.requestToken =
      provider.requestToken ?? this.nextRequestToken(provider.providerId)
    current.lastFailureKind = provider.lastFailureKind
  }

  private nextRequestToken(providerId: ProviderSettingsProviderId): string {
    this.requestSequence += 1
    return `${providerId}-${this.requestSequence}`
  }
}
