import {
  createProductionAppFactories,
  createReviewFakeApiAppFactories
} from "./bootstrap/app-screen-controller-factories"
import { createDefaultReviewFakeApiGatewayRegistry } from "@controller/review-fake-api/default-review-fake-api-gateway-registry"
import { resolveReviewFakeApiRuntimeContext } from "@controller/review-fake-api/review-fake-api-runtime"
import { mount } from "svelte"
import App from "@ui/App.svelte"

const target = document.getElementById("app")

if (!target) {
  throw new Error("app root not found")
}

const reviewFakeApiContext = resolveReviewFakeApiRuntimeContext(
  new URLSearchParams(window.location.search),
  {
    reviewModeEnabled: import.meta.env.DEV
  }
)
const appFactories = reviewFakeApiContext.enabled
  ? createReviewFakeApiAppFactories(
      reviewFakeApiContext,
      createDefaultReviewFakeApiGatewayRegistry()
    )
  : createProductionAppFactories()

mount(App, {
  target,
  props: {
    createBodyTranslationPhaseScreenController:
      appFactories.createBodyTranslationPhaseScreenController,
    createMasterDictionaryScreenController:
      appFactories.createMasterDictionaryScreenController,
    createPersonaGenerationPhaseScreenController:
      appFactories.createPersonaGenerationPhaseScreenController,
    createProviderSettingsScreenController:
      appFactories.createProviderSettingsScreenController,
    createTermTranslationPhaseScreenController:
      appFactories.createTermTranslationPhaseScreenController,
    createTranslationInputScreenController:
      appFactories.createTranslationInputScreenController,
    createTranslationOutputArtifactScreenController:
      appFactories.createTranslationOutputArtifactScreenController,
    createTranslationJobSetupScreenController:
      appFactories.createTranslationJobSetupScreenController,
    createMasterPersonaScreenController:
      appFactories.createMasterPersonaScreenController
  }
})
