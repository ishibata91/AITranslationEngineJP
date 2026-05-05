import { render, screen, waitFor } from "@testing-library/svelte"
import { describe, expect, test, vi } from "vitest"

import type {
  ListProviderSettingsResponse,
  ProviderSettingsGatewayContract,
  ProviderSettingsProviderId,
  ResetProviderSettingsRequest,
  SaveProviderSettingsRequest,
  ProviderSettingsSummary,
  ValidateProviderSettingsRequest,
  ValidateProviderSettingsResponse
} from "@application/gateway-contract/provider-settings"
/* eslint-disable local/enforce-layer-boundaries */
import {
  resolveReviewFakeApiRuntimeContext
} from "@controller/review-fake-api/review-fake-api-runtime"
import { createReviewFakeApiAppFactories } from "../bootstrap/app-screen-controller-factories"
/* eslint-enable local/enforce-layer-boundaries */
import App from "@ui/App.svelte"

const REVIEW_PROVIDER_IDS = ["gemini", "xai", "lm_studio"] as const
type ReviewFakeApiScenarioId =
  | "empty"
  | "loading"
  | "success"
  | "running"
  | "error"
  | "config-missing"

interface DeferredResponse {
  promise: Promise<ListProviderSettingsResponse>
  resolve: (response: ListProviderSettingsResponse) => void
}

function createDeferredResponse(): DeferredResponse {
  let resolveResponse:
    | ((response: ListProviderSettingsResponse) => void)
    | null = null
  const promise = new Promise<ListProviderSettingsResponse>((resolve) => {
    resolveResponse = resolve
  })

  if (!resolveResponse) {
    throw new Error("deferred response initializer failed")
  }

  return {
    promise,
    resolve: resolveResponse
  }
}

function createProvider(
  providerId: ProviderSettingsProviderId,
  overrides: Partial<ProviderSettingsSummary> = {}
): ProviderSettingsSummary {
  const baseProvider: ProviderSettingsSummary = {
    providerId,
    label:
      providerId === "lm_studio"
        ? "LM Studio"
        : providerId === "xai"
          ? "xAI"
          : "Gemini",
    endpoint:
      providerId === "lm_studio"
        ? "http://127.0.0.1:1234/v1"
        : "https://example.invalid",
    credentialState: providerId === "lm_studio" ? "not_required" : "configured",
    validationState: "validated",
    savedState: "configured",
    requestToken: `${providerId}-review`
  }

  return {
    ...baseProvider,
    ...overrides
  }
}

function createResponse(
  providers: ProviderSettingsSummary[]
): ListProviderSettingsResponse {
  return {
    route: {
      routeId: "provider-settings",
      label: "AIサービス設定",
      currentRouteState: "レビュー確認",
      dashboardEntryId: "provider-settings"
    },
    providers
  }
}

function createProvidersForScenario(
  scenarioId: ReviewFakeApiScenarioId
): ProviderSettingsSummary[] {
  if (scenarioId === "empty") {
    return []
  }

  if (scenarioId === "running") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        validationState: "pending"
      })
    )
  }

  if (scenarioId === "error") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        validationState: "failed",
        lastFailureKind: "provider_unreachable"
      })
    )
  }

  if (scenarioId === "config-missing") {
    return REVIEW_PROVIDER_IDS.map((providerId) =>
      createProvider(providerId, {
        credentialState:
          providerId === "lm_studio" ? "not_required" : "missing",
        validationState:
          providerId === "lm_studio" ? "validated" : "not_validated",
        savedState: providerId === "lm_studio" ? "configured" : "partial",
        lastFailureKind:
          providerId === "lm_studio" ? undefined : "credential_missing"
      })
    )
  }

  return REVIEW_PROVIDER_IDS.map((providerId) => createProvider(providerId))
}

function createProviderSettingsGateway(
  scenarioId: ReviewFakeApiScenarioId,
  deferredResponse: DeferredResponse | null = null
): {
  gateway: ProviderSettingsGatewayContract
  listProviderSettings: ReturnType<typeof vi.fn>
} {
  const listResponse =
    deferredResponse?.promise ??
    Promise.resolve(createResponse(createProvidersForScenario(scenarioId)))
  const listProviderSettings = vi.fn(() => listResponse)

  return {
    gateway: {
      ListProviderSettings: listProviderSettings,
      SaveProviderSettings: vi.fn((request: SaveProviderSettingsRequest) =>
        Promise.resolve({
          provider: createProvider(request.providerId)
        })
      ),
      ResetProviderSettings: vi.fn((request: ResetProviderSettingsRequest) =>
        Promise.resolve({
          provider: createProvider(request.providerId)
        })
      ),
      ValidateProviderSettings: vi.fn(
        (
          request: ValidateProviderSettingsRequest
        ): Promise<ValidateProviderSettingsResponse> =>
          Promise.resolve({
            providerId: request.providerId,
            validationState: "validated",
            requestToken: request.requestToken
          })
      )
    },
    listProviderSettings
  }
}

function renderProviderSettingsReviewUrl(
  scenarioId: ReviewFakeApiScenarioId,
  deferredResponse: DeferredResponse | null = null
): {
  providerSettingsGateway: ProviderSettingsGatewayContract
  listProviderSettings: ReturnType<typeof vi.fn>
} {
  window.history.replaceState(
    null,
    "",
    `/?fakeApi=1&fakeScenario=${scenarioId}#provider-settings`
  )

  const context = resolveReviewFakeApiRuntimeContext(
    new URLSearchParams(window.location.search),
    { reviewModeEnabled: true }
  )
  const { gateway: providerSettingsGateway, listProviderSettings } =
    createProviderSettingsGateway(context.scenarioId, deferredResponse)
  const factories = createReviewFakeApiAppFactories(context, {
    providerSettings: () => providerSettingsGateway
  })

  render(App, {
    props: {
      createBodyTranslationPhaseScreenController:
        factories.createBodyTranslationPhaseScreenController,
      createMasterDictionaryScreenController:
        factories.createMasterDictionaryScreenController,
      createMasterPersonaScreenController:
        factories.createMasterPersonaScreenController,
      createPersonaGenerationPhaseScreenController:
        factories.createPersonaGenerationPhaseScreenController,
      createProviderSettingsScreenController:
        factories.createProviderSettingsScreenController,
      createTermTranslationPhaseScreenController:
        factories.createTermTranslationPhaseScreenController,
      createTranslationInputScreenController:
        factories.createTranslationInputScreenController,
      createTranslationJobSetupScreenController:
        factories.createTranslationJobSetupScreenController,
      createTranslationOutputArtifactScreenController:
        factories.createTranslationOutputArtifactScreenController
    }
  })

  return {
    providerSettingsGateway,
    listProviderSettings
  }
}

describe("frontend fakeAPI review scenario", () => {
  test("SCN-FFARF-001: fakeAPI レビュー起動 URL で backend なしの provider settings 画面を開く", async () => {
    const { listProviderSettings } = renderProviderSettingsReviewUrl("success")

    expect(
      screen.getByRole("heading", { level: 1, name: "AIサービス設定" })
    ).toBeInTheDocument()

    await waitFor(() => {
      expect(listProviderSettings).toHaveBeenCalledTimes(1)
      expect(
        screen.getAllByText("Gateway: 接続準備済み")[0]
      ).toBeInTheDocument()
      expect(screen.getAllByText("利用可能")[0]).toBeInTheDocument()
    })

    expect(screen.queryByText("fakeAPI")).not.toBeInTheDocument()
  })

  test("SCN-FFARF-002: empty 状態は件数 0 と detail 非表示で確認できる", async () => {
    renderProviderSettingsReviewUrl("empty")

    await waitFor(() => {
      expect(
        screen.getByText("0 件の AIサービスを管理します。")
      ).toBeInTheDocument()
    })

    expect(
      screen.queryByRole("heading", { level: 3, name: "Gemini" })
    ).not.toBeInTheDocument()
  })

  test("SCN-FFARF-002: loading 状態は読込中として確認できる", async () => {
    const deferredResponse = createDeferredResponse()

    renderProviderSettingsReviewUrl("loading", deferredResponse)

    await waitFor(() => {
      expect(screen.getByText("読込中")).toBeInTheDocument()
    })

    deferredResponse.resolve(createResponse([]))
  })

  test.each([
    {
      scenarioId: "success",
      expectedText: "接続確認済みです。"
    },
    {
      scenarioId: "running",
      expectedText: "現在の入力内容で接続確認しています。"
    },
    {
      scenarioId: "error",
      expectedText: "エンドポイントを確認してから再実行してください。"
    },
    {
      scenarioId: "config-missing",
      expectedText: "保存後に接続確認してください。"
    }
  ] as const)(
    "SCN-FFARF-002: $scenarioId 状態は provider settings 表示で確認できる",
    async ({ scenarioId, expectedText }) => {
      renderProviderSettingsReviewUrl(scenarioId)

      await waitFor(() => {
        expect(screen.getAllByText(expectedText)[0]).toBeInTheDocument()
      })
    }
  )
})
