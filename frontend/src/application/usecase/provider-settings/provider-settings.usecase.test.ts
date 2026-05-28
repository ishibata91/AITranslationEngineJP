import { describe, expect, test } from "vitest"

import type {
  ProviderSettingsGatewayContract,
  ProviderSettingsScreenState
} from "@application/gateway-contract/provider-settings"

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

  test("保存拒否された不正 endpoint は入力不正表示へ写像し保存成功 notice を出さない", async () => {
    const store = createStore()
    const gateway: ProviderSettingsGatewayContract = {
      ListProviderSettings: () =>
        Promise.resolve({
          route: {
            routeId: "provider-settings",
            label: "Provider Settings",
            currentRouteState: "ready",
            dashboardEntryId: "provider-settings"
          },
          providers: [
            {
              providerId: "gemini",
              label: "Gemini",
              endpoint: "https://gemini.example/v1",
              credentialState: "configured",
              validationState: "validated",
              savedState: "configured",
              requestToken: "gemini-initial"
            }
          ]
        }),
      SaveProviderSettings: () =>
        Promise.reject(
          new Error(
            "save provider settings: save provider settings usecase: provider settings validation: validation_failed endpoint format is invalid"
          )
        ),
      ResetProviderSettings: () => Promise.reject(new Error("unused")),
      ValidateProviderSettings: () => Promise.reject(new Error("unused"))
    }
    const useCase = new ProviderSettingsUseCase(gateway, store)

    await useCase.load()
    useCase.updateEndpoint("invalid-endpoint")
    await useCase.saveSettings(() => "super-secret-token")

    const snapshot = store.snapshot()
    const provider = snapshot.providers[0]
    expect(snapshot.phase).toBe("ready")
    expect(snapshot.errorMessage).toMatch(/入力|不正/)
    expect(snapshot.saveNotice).toBe("")
    expect(provider.endpointDraft).toBe("invalid-endpoint")
    expect(provider.persistedEndpoint).toBe("https://gemini.example/v1")
    expect(provider.validationState).toBe("failed")
    expect(provider.lastFailureKind).toBe("validation_failed")
    expect(JSON.stringify(snapshot)).not.toContain("super-secret-token")
  })

  test("非検証エラーの保存失敗は入力不正状態と新 requestToken へ進めない", async () => {
    const store = createStore()
    const gateway: ProviderSettingsGatewayContract = {
      ListProviderSettings: () =>
        Promise.resolve({
          route: {
            routeId: "provider-settings",
            label: "Provider Settings",
            currentRouteState: "ready",
            dashboardEntryId: "provider-settings"
          },
          providers: [
            {
              providerId: "gemini",
              label: "Gemini",
              endpoint: "https://gemini.example/v1",
              credentialState: "configured",
              validationState: "validated",
              savedState: "configured",
              requestToken: "gemini-initial"
            }
          ]
        }),
      SaveProviderSettings: () =>
        Promise.reject(
          new Error(
            "save provider settings: secret store save failed after rollback"
          )
        ),
      ResetProviderSettings: () => Promise.reject(new Error("unused")),
      ValidateProviderSettings: () => Promise.reject(new Error("unused"))
    }
    const useCase = new ProviderSettingsUseCase(gateway, store)

    await useCase.load()
    useCase.updateEndpoint("https://changed.example/v1")
    const beforeSave = store.snapshot().providers[0]

    await useCase.saveSettings(() => "super-secret-token")

    const snapshot = store.snapshot()
    const provider = snapshot.providers[0]
    expect(snapshot.phase).toBe("ready")
    expect(snapshot.errorMessage).toBe("設定の保存に失敗しました。")
    expect(snapshot.saveNotice).toBe("")
    expect(provider.endpointDraft).toBe("https://changed.example/v1")
    expect(provider.persistedEndpoint).toBe("https://gemini.example/v1")
    expect(provider.savedState).toBe("configured")
    expect(provider.validationState).toBe(beforeSave.validationState)
    expect(provider.lastFailureKind).toBe(beforeSave.lastFailureKind)
    expect(provider.requestToken).toBe(beforeSave.requestToken)
    expect(JSON.stringify(snapshot)).not.toContain("super-secret-token")
  })
})
