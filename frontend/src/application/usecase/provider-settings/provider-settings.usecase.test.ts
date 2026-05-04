import { describe, expect, test } from "vitest"

import type { ProviderSettingsScreenState } from "@application/gateway-contract/provider-settings"

import { ProviderSettingsUseCase } from "./provider-settings.usecase"

type StoreLike = {
  snapshot(): ProviderSettingsScreenState
  update(mutator: (draft: ProviderSettingsScreenState) => void): void
}

function createStore(): StoreLike {
  let state: ProviderSettingsScreenState = {
    phase: "idle",
    providers: [],
    selectedProviderId: null,
    apiKeyPanelOpen: false,
    saveNotice: "",
    errorMessage: ""
  }

  return {
    snapshot() {
      return {
        ...state,
        providers: state.providers.map((provider) => ({ ...provider }))
      }
    },
    update(mutator) {
      const nextState = {
        ...state,
        providers: state.providers.map((provider) => ({ ...provider }))
      }
      mutator(nextState)
      state = nextState
    }
  }
}

describe("ProviderSettingsUseCase", () => {
  test("APIキー保存後も screen state に入力値を残さない", async () => {
    const store = createStore()
    const useCase = new ProviderSettingsUseCase(null, store)

    await useCase.load()
    useCase.selectProvider("xai")
    useCase.openApiKeyPanel()
    await useCase.saveSettings(() => "super-secret-token")

    const snapshot = store.snapshot()
    expect(JSON.stringify(snapshot)).not.toContain("super-secret-token")
    expect(
      snapshot.providers.find((provider) => provider.providerId === "xai")
        ?.credentialState
    ).toBe("configured")
  })

  test("遅延した接続確認 response を現在入力へ混入させない", async () => {
    const store = createStore()
    const useCase = new ProviderSettingsUseCase(null, store)

    await useCase.load()
    useCase.selectProvider("gemini")
    const validationPromise = useCase.validateConnection()
    useCase.updateEndpoint("https://changed.example.test")
    await validationPromise

    const provider = store
      .snapshot()
      .providers.find((item) => item.providerId === "gemini")
    expect(provider?.endpointDraft).toBe("https://changed.example.test")
    expect(provider?.validationState).toBe("not_validated")
    expect(provider?.lastFailureKind).toBe("validation_stale")
  })
})
