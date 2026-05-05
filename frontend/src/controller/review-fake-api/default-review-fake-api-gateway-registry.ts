import type {
  ListProviderSettingsResponse,
  ProviderSettingsGatewayContract,
  ProviderSettingsProviderId,
  ProviderSettingsSummary,
  ResetProviderSettingsRequest,
  SaveProviderSettingsRequest,
  ValidateProviderSettingsRequest,
  ValidateProviderSettingsResponse
} from "@application/gateway-contract/provider-settings"

import type {
  ReviewFakeApiGatewayRegistry,
  ReviewFakeApiScenarioId
} from "./review-fake-api-runtime"

const REVIEW_PROVIDER_IDS = ["gemini", "xai", "lm_studio"] as const

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

function createPendingProviderSettingsGateway(): ProviderSettingsGatewayContract {
  return {
    ListProviderSettings: () => new Promise(() => undefined),
    SaveProviderSettings: (request: SaveProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ResetProviderSettings: (request: ResetProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ValidateProviderSettings: (
      request: ValidateProviderSettingsRequest
    ): Promise<ValidateProviderSettingsResponse> =>
      Promise.resolve({
        providerId: request.providerId,
        validationState: "validated",
        requestToken: request.requestToken
      })
  }
}

function createProviderSettingsGateway(
  scenarioId: ReviewFakeApiScenarioId
): ProviderSettingsGatewayContract {
  if (scenarioId === "loading") {
    return createPendingProviderSettingsGateway()
  }

  return {
    ListProviderSettings: () =>
      Promise.resolve(createResponse(createProvidersForScenario(scenarioId))),
    SaveProviderSettings: (request: SaveProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ResetProviderSettings: (request: ResetProviderSettingsRequest) =>
      Promise.resolve({
        provider: createProvider(request.providerId)
      }),
    ValidateProviderSettings: (
      request: ValidateProviderSettingsRequest
    ): Promise<ValidateProviderSettingsResponse> =>
      Promise.resolve({
        providerId: request.providerId,
        validationState: "validated",
        requestToken: request.requestToken
      })
  }
}

export function createDefaultReviewFakeApiGatewayRegistry(): ReviewFakeApiGatewayRegistry {
  return {
    providerSettings: (context) => createProviderSettingsGateway(context.scenarioId)
  }
}
