import { afterEach, describe, expect, test, vi } from "vitest"

import { createProviderSettingsGateway } from "./provider-settings.gateway"

type ProviderSettingsBindings = {
  ListProviderSettings?: ReturnType<typeof vi.fn>
  SaveProviderSettings?: ReturnType<typeof vi.fn>
  ResetProviderSettings?: ReturnType<typeof vi.fn>
  ValidateProviderSettings?: ReturnType<typeof vi.fn>
}

type GoRecord = {
  wails: {
    ProviderSettingsController?: ProviderSettingsBindings
    AppController?: ProviderSettingsBindings
  }
}

const originalGo: unknown = Reflect.get(globalThis as object, "go")

function installGo(record: GoRecord): void {
  Object.defineProperty(globalThis, "go", {
    value: record,
    configurable: true,
    writable: true
  })
}

function providerSummary() {
  return {
    providerId: "gemini" as const,
    label: "Gemini",
    endpoint: "https://example.invalid",
    credentialState: "configured" as const,
    credentialReferenceId: "credential-a",
    validationState: "validated" as const,
    savedState: "configured" as const,
    requestToken: "token-a"
  }
}

afterEach(() => {
  vi.restoreAllMocks()
  Object.defineProperty(globalThis, "go", {
    value: originalGo,
    configurable: true,
    writable: true
  })
})

describe("createProviderSettingsGateway", () => {
  test("ListProviderSettings は request 省略時も空 request で Wails binding を呼ぶ", async () => {
    const listProviderSettings = vi.fn(() =>
      Promise.resolve({
        route: {
          routeId: "provider-settings",
          label: "Provider Settings",
          currentRouteState: "ready",
          dashboardEntryId: "provider-settings"
        },
        providers: [providerSummary()]
      })
    )
    installGo({
      wails: {
        ProviderSettingsController: {
          ListProviderSettings: listProviderSettings
        }
      }
    })

    const gateway = createProviderSettingsGateway()

    await expect(gateway.ListProviderSettings()).resolves.toMatchObject({
      providers: [{ providerId: "gemini" }]
    })
    expect(listProviderSettings).toHaveBeenCalledTimes(1)
    expect(listProviderSettings).toHaveBeenCalledWith({})
  })

  test("各 mutation は request を Wails binding へ渡す", async () => {
    const saveProviderSettings = vi.fn(() =>
      Promise.resolve({ provider: providerSummary() })
    )
    const resetProviderSettings = vi.fn(() =>
      Promise.resolve({ provider: providerSummary() })
    )
    const validateProviderSettings = vi.fn(() =>
      Promise.resolve({
        providerId: "gemini",
        validationState: "validated",
        requestToken: "token-a"
      })
    )
    installGo({
      wails: {
        ProviderSettingsController: {
          SaveProviderSettings: saveProviderSettings,
          ResetProviderSettings: resetProviderSettings,
          ValidateProviderSettings: validateProviderSettings
        }
      }
    })

    const gateway = createProviderSettingsGateway()

    await gateway.SaveProviderSettings({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      apiKeyInputPresent: true,
      credentialInput: "secret"
    })
    await gateway.ResetProviderSettings({ providerId: "gemini" })
    await gateway.ValidateProviderSettings({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      credentialState: "configured",
      credentialReferenceId: "credential-a",
      requestToken: "token-a"
    })

    expect(saveProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      apiKeyInputPresent: true,
      credentialInput: "secret"
    })
    expect(resetProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini"
    })
    expect(validateProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      credentialState: "configured",
      credentialReferenceId: "credential-a",
      requestToken: "token-a"
    })
  })

  test("ProviderSettingsController が未接続なら AppController の binding を使う", async () => {
    const listProviderSettings = vi.fn(() =>
      Promise.resolve({
        route: {
          routeId: "provider-settings",
          label: "Provider Settings",
          currentRouteState: "ready",
          dashboardEntryId: "provider-settings"
        },
        providers: [providerSummary()]
      })
    )
    installGo({
      wails: {
        AppController: {
          ListProviderSettings: listProviderSettings
        }
      }
    })

    const gateway = createProviderSettingsGateway()

    await gateway.ListProviderSettings({})

    expect(listProviderSettings).toHaveBeenCalledTimes(1)
    expect(listProviderSettings).toHaveBeenCalledWith({})
  })

  test("binding が未接続なら reject する", async () => {
    installGo({
      wails: {}
    })

    const gateway = createProviderSettingsGateway()

    await expect(gateway.ListProviderSettings()).rejects.toThrow(
      "Wails binding is not wired yet: ListProviderSettings"
    )
  })
})
