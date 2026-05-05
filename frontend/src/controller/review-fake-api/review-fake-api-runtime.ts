import type { BodyTranslationPhaseGatewayContract } from "@application/gateway-contract/body-translation-phase"
import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
import type { MasterPersonaGatewayContract } from "@application/gateway-contract/master-persona"
import type { PersonaGenerationPhaseGatewayContract } from "@application/gateway-contract/persona-generation-phase"
import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import type { TranslationInputGatewayContract } from "@application/gateway-contract/translation-input"
import type { TranslationJobSetupGatewayContract } from "@application/gateway-contract/translation-job-setup"
import type { TranslationOutputArtifactGatewayContract } from "@application/gateway-contract/translation-output-artifact"

const REVIEW_FAKE_API_SCENARIO_IDS = [
  "empty",
  "loading",
  "success",
  "running",
  "error",
  "config-missing"
] as const

const REVIEW_FAKE_API_TRIGGER_PARAM = "fakeApi"
const REVIEW_FAKE_API_SCENARIO_PARAM = "fakeScenario"
const REVIEW_FAKE_API_DEFAULT_SCENARIO_ID = "empty"
const REVIEW_FAKE_API_FALLBACK_SCENARIO_ID = "config-missing"

export type ReviewFakeApiScenarioId =
  (typeof REVIEW_FAKE_API_SCENARIO_IDS)[number]

export interface ReviewFakeApiRuntimeContext {
  enabled: boolean
  scenarioId: ReviewFakeApiScenarioId
  defaultScenarioId: ReviewFakeApiScenarioId
  triggerValue: string | null
  overrideValue: string | null
}

interface UrlSearchParamsLike {
  get(name: string): string | null
}

interface ReviewFakeApiRuntimeOptions {
  reviewModeEnabled: boolean
}

export interface ReviewFakeApiGatewayRegistry {
  bodyTranslationPhase?: (
    context: ReviewFakeApiRuntimeContext
  ) => BodyTranslationPhaseGatewayContract | null
  masterDictionary?: (
    context: ReviewFakeApiRuntimeContext
  ) => MasterDictionaryGatewayContract | null
  masterPersona?: (
    context: ReviewFakeApiRuntimeContext
  ) => MasterPersonaGatewayContract | null
  personaGenerationPhase?: (
    context: ReviewFakeApiRuntimeContext
  ) => PersonaGenerationPhaseGatewayContract | null
  providerSettings?: (
    context: ReviewFakeApiRuntimeContext
  ) => ProviderSettingsGatewayContract | null
  termTranslationPhase?: (
    context: ReviewFakeApiRuntimeContext
  ) => TermTranslationPhaseGatewayContract | null
  translationInput?: (
    context: ReviewFakeApiRuntimeContext
  ) => TranslationInputGatewayContract | null
  translationJobSetup?: (
    context: ReviewFakeApiRuntimeContext
  ) => TranslationJobSetupGatewayContract | null
  translationOutputArtifact?: (
    context: ReviewFakeApiRuntimeContext
  ) => TranslationOutputArtifactGatewayContract | null
}

function isReviewFakeApiScenarioId(
  value: string
): value is ReviewFakeApiScenarioId {
  return (REVIEW_FAKE_API_SCENARIO_IDS as readonly string[]).includes(value)
}

function normalizeScenarioId(
  value: string | null,
  fallbackScenarioId: ReviewFakeApiScenarioId
): ReviewFakeApiScenarioId {
  if (!value) {
    return fallbackScenarioId
  }

  return isReviewFakeApiScenarioId(value)
    ? value
    : REVIEW_FAKE_API_FALLBACK_SCENARIO_ID
}

function resolveDefaultScenarioId(
  triggerValue: string | null
): ReviewFakeApiScenarioId {
  if (!triggerValue || triggerValue === "1" || triggerValue === "true") {
    return REVIEW_FAKE_API_DEFAULT_SCENARIO_ID
  }

  return normalizeScenarioId(
    triggerValue,
    REVIEW_FAKE_API_DEFAULT_SCENARIO_ID
  )
}

export function resolveReviewFakeApiRuntimeContext(
  searchParams: UrlSearchParamsLike,
  options: ReviewFakeApiRuntimeOptions
): ReviewFakeApiRuntimeContext {
  const triggerValue = searchParams.get(REVIEW_FAKE_API_TRIGGER_PARAM)
  const enabled = options.reviewModeEnabled && triggerValue !== null
  const defaultScenarioId = enabled
    ? resolveDefaultScenarioId(triggerValue)
    : REVIEW_FAKE_API_DEFAULT_SCENARIO_ID
  const overrideValue = enabled
    ? searchParams.get(REVIEW_FAKE_API_SCENARIO_PARAM)
    : null

  return {
    enabled,
    scenarioId: enabled
      ? normalizeScenarioId(overrideValue, defaultScenarioId)
      : REVIEW_FAKE_API_DEFAULT_SCENARIO_ID,
    defaultScenarioId,
    triggerValue,
    overrideValue
  }
}
