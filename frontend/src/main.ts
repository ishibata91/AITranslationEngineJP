import { createProductionAppFactories } from "./bootstrap/app-screen-controller-factories"
import { mount } from "svelte"
import App from "@ui/App.svelte"

const target = document.getElementById("app")

if (!target) {
  throw new Error("app root not found")
}

const appFactories = createProductionAppFactories()

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
    createTranslationJobManagementScreenController:
      appFactories.createTranslationJobManagementScreenController,
    createTranslationInputScreenController:
      appFactories.createTranslationInputScreenController,
    createTranslationOutputArtifactScreenController:
      appFactories.createTranslationOutputArtifactScreenController,
    createMasterPersonaScreenController:
      appFactories.createMasterPersonaScreenController
  }
})
