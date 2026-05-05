import { beforeEach, describe, expect, test, vi } from "vitest"

const mocks = vi.hoisted(() => ({
  createBodyTranslationPhaseScreenControllerFactory: vi.fn(() => vi.fn()),
  createMasterDictionaryScreenControllerFactory: vi.fn(() => vi.fn()),
  createMasterPersonaScreenControllerFactory: vi.fn(() => vi.fn()),
  createPersonaGenerationPhaseScreenControllerFactory: vi.fn(() => vi.fn()),
  createProviderSettingsScreenControllerFactory: vi.fn(() => vi.fn()),
  createTermTranslationPhaseScreenControllerFactory: vi.fn(() => vi.fn()),
  createTranslationInputScreenControllerFactory: vi.fn(() => vi.fn()),
  createTranslationJobSetupScreenControllerFactory: vi.fn(() => vi.fn()),
  createTranslationOutputArtifactScreenControllerFactory: vi.fn(() => vi.fn())
}))

vi.mock("@controller/body-translation-phase", () => ({
  createBodyTranslationPhaseScreenControllerFactory:
    mocks.createBodyTranslationPhaseScreenControllerFactory
}))
vi.mock("@controller/master-dictionary", () => ({
  createMasterDictionaryScreenControllerFactory:
    mocks.createMasterDictionaryScreenControllerFactory
}))
vi.mock("@controller/master-persona", () => ({
  createMasterPersonaScreenControllerFactory:
    mocks.createMasterPersonaScreenControllerFactory
}))
vi.mock("@controller/persona-generation-phase", () => ({
  createPersonaGenerationPhaseScreenControllerFactory:
    mocks.createPersonaGenerationPhaseScreenControllerFactory
}))
vi.mock("@controller/provider-settings", () => ({
  createProviderSettingsScreenControllerFactory:
    mocks.createProviderSettingsScreenControllerFactory
}))
vi.mock("@controller/term-translation-phase", () => ({
  createTermTranslationPhaseScreenControllerFactory:
    mocks.createTermTranslationPhaseScreenControllerFactory
}))
vi.mock("@controller/translation-input", () => ({
  createTranslationInputScreenControllerFactory:
    mocks.createTranslationInputScreenControllerFactory
}))
vi.mock("@controller/translation-job-setup", () => ({
  createTranslationJobSetupScreenControllerFactory:
    mocks.createTranslationJobSetupScreenControllerFactory
}))
vi.mock("@controller/translation-output-artifact", () => ({
  createTranslationOutputArtifactScreenControllerFactory:
    mocks.createTranslationOutputArtifactScreenControllerFactory
}))

import { createReviewFakeApiAppFactories } from "../../bootstrap/app-screen-controller-factories"

type ReviewFakeApiRuntimeContext =
  import("./review-fake-api-runtime").ReviewFakeApiRuntimeContext

function createContext(
  partial: Partial<ReviewFakeApiRuntimeContext> = {}
): ReviewFakeApiRuntimeContext {
  return {
    enabled: true,
    scenarioId: "success",
    defaultScenarioId: "empty",
    triggerValue: "1",
    overrideValue: "success",
    ...partial
  }
}

describe("createReviewFakeApiAppFactories", () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  test("fakeAPI 起動時は registry から渡した gateway で DI 差し替えする", () => {
    const context = createContext()
    const fakeGateway = { load: vi.fn() }
    const registry = {
      translationInput: vi.fn(() => fakeGateway as never)
    }

    createReviewFakeApiAppFactories(context, registry)

    expect(registry.translationInput).toHaveBeenCalledWith(context)
    expect(
      mocks.createTranslationInputScreenControllerFactory
    ).toHaveBeenCalledWith(fakeGateway)
  })

  test("モックデータ欠落時は null を渡し Wails gateway fallback へ流さない", () => {
    const context = createContext({ scenarioId: "config-missing" })
    const registry = {
      translationInput: vi.fn(() => null)
    }

    createReviewFakeApiAppFactories(context, registry)

    expect(mocks.createTranslationInputScreenControllerFactory).toHaveBeenCalledWith(
      null
    )
  })

  test("registry 未登録時も null を渡し本番初期状態へ切り替えない", () => {
    createReviewFakeApiAppFactories(createContext(), {})

    expect(mocks.createTranslationInputScreenControllerFactory).toHaveBeenCalledWith(
      null
    )
  })
})
