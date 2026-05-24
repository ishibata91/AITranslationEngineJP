import { afterEach, describe, expect, test, vi } from "vitest"

import {
  ListProviderSettings,
  ResetProviderSettings,
  SaveProviderSettings,
  ValidateProviderSettings
} from "../../../wailsjs/go/wails/AppController.js"
import { createProviderSettingsGateway } from "./provider-settings.gateway"

type ListProviderSettingsResponse = Awaited<
  ReturnType<typeof ListProviderSettings>
>
type SaveProviderSettingsResponse = Awaited<
  ReturnType<typeof SaveProviderSettings>
>
type ResetProviderSettingsResponse = Awaited<
  ReturnType<typeof ResetProviderSettings>
>
type ValidateProviderSettingsResponse = Awaited<
  ReturnType<typeof ValidateProviderSettings>
>

vi.mock("../../../wailsjs/go/wails/AppController.js", () => ({
  ListProviderSettings: vi.fn(),
  SaveProviderSettings: vi.fn(),
  ResetProviderSettings: vi.fn(),
  ValidateProviderSettings: vi.fn()
}))

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
})

describe("createProviderSettingsGateway", () => {
  test("ListProviderSettings は request 省略時も空 request で公開 binding wrapper を呼ぶ", async () => {
    // 公開 seam: generated binding wrapper 経由で request を送る。
    vi.mocked(ListProviderSettings).mockResolvedValue({
      route: {
        routeId: "provider-settings",
        label: "Provider Settings",
        currentRouteState: "ready",
        dashboardEntryId: "provider-settings"
      },
      providers: [providerSummary()]
    } as unknown as ListProviderSettingsResponse)

    const gateway = createProviderSettingsGateway()

    await expect(gateway.ListProviderSettings()).resolves.toMatchObject({
      providers: [{ providerId: "gemini" }]
    })
    expect(ListProviderSettings).toHaveBeenCalledTimes(1)
    expect(ListProviderSettings).toHaveBeenCalledWith({})
  })

  test("各 mutation は request を公開 binding wrapper へ渡す", async () => {
    // 公開 seam: request payload は gateway で加工せず binding へ渡す。
    vi.mocked(SaveProviderSettings).mockResolvedValue({
      provider: providerSummary()
    } as unknown as SaveProviderSettingsResponse)
    vi.mocked(ResetProviderSettings).mockResolvedValue({
      provider: providerSummary()
    } as unknown as ResetProviderSettingsResponse)
    vi.mocked(ValidateProviderSettings).mockResolvedValue({
      providerId: "gemini",
      validationState: "validated",
      requestToken: "token-a"
    } as unknown as ValidateProviderSettingsResponse)

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

    expect(SaveProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      apiKeyInputPresent: true,
      credentialInput: "secret"
    })
    expect(ResetProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini"
    })
    expect(ValidateProviderSettings).toHaveBeenCalledWith({
      providerId: "gemini",
      endpoint: "https://example.invalid",
      credentialState: "configured",
      credentialReferenceId: "credential-a",
      requestToken: "token-a"
    })
  })

  test("binding が未接続なら wrapper 例外をそのまま返す", async () => {
    // 未接続時は公開 seam で返った例外を返す。
    vi.mocked(ListProviderSettings).mockRejectedValue(
      new Error("Wails binding is not wired yet: ListProviderSettings")
    )

    const gateway = createProviderSettingsGateway()

    await expect(gateway.ListProviderSettings()).rejects.toThrow(
      "Wails binding is not wired yet: ListProviderSettings"
    )
  })

  test("runtime shape 検証失敗時は診断に secret 平文を含めない", async () => {
    // runtime shape 失敗時の公開値に secret が漏れないことを確認する。
    vi.mocked(ListProviderSettings).mockResolvedValue({
      route: {
        routeId: "provider-settings",
        label: "Provider Settings",
        currentRouteState: "ready",
        dashboardEntryId: "provider-settings"
      },
      providers: [
        {
          ...providerSummary(),
          providerId: "unknown-provider",
          credentialInput: "raw-secret-value"
        }
      ]
    } as unknown as ListProviderSettingsResponse)

    const gateway = createProviderSettingsGateway()

    await expect(gateway.ListProviderSettings()).rejects.toMatchObject({
      name: "GatewayResponseShapeError",
      userFacingMessage: "Gateway response shape is invalid."
    })

    try {
      await gateway.ListProviderSettings()
    } catch (error: unknown) {
      expect(error).toBeInstanceOf(Error)
      expect(typeof (error as Error & { internalDiagnostic?: unknown }).internalDiagnostic).toBe("string")
      const diagnostic = JSON.stringify(error)
      expect(diagnostic).not.toContain("raw-secret-value")
      expect(diagnostic).not.toContain("credentialInput")
      expect(diagnostic).not.toContain("apiKey")
    }
  })
})
