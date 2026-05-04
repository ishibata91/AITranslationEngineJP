import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { createBodyTranslationPhaseScreenControllerFactory } from "@controller/body-translation-phase"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"
import { createPersonaGenerationPhaseScreenControllerFactory } from "@controller/persona-generation-phase"
import { createProviderSettingsScreenControllerFactory } from "@controller/provider-settings"
import { createTermTranslationPhaseScreenControllerFactory } from "@controller/term-translation-phase"
import { createTranslationInputScreenControllerFactory } from "@controller/translation-input"
import { createTranslationOutputArtifactScreenControllerFactory } from "@controller/translation-output-artifact"
import { createTranslationJobSetupScreenControllerFactory } from "@controller/translation-job-setup"
import { createBodyTranslationPhaseGateway } from "@controller/wails/body-translation-phase.gateway"
import { createProviderSettingsGateway } from "@controller/wails/provider-settings.gateway"
import { createTranslationOutputArtifactGateway } from "@controller/wails/translation-output-artifact.gateway"
import { createTranslationInputGateway } from "@controller/wails/translation-input.gateway"
import { createTranslationJobSetupGateway } from "@controller/wails/translation-job-setup.gateway"
import { createTermTranslationPhaseGateway } from "@controller/wails/term-translation-phase.gateway"
import { createPersonaGenerationPhaseGateway } from "@controller/wails/persona-generation-phase.gateway"
import { mount } from "svelte"
import App from "@ui/App.svelte"

const target = document.getElementById("app")

if (!target) {
  throw new Error("app root not found")
}

const masterDictionaryGateway = createMasterDictionaryGateway()
const masterDictionaryScreenControllerFactory =
  createMasterDictionaryScreenControllerFactory(masterDictionaryGateway)
const termTranslationPhaseGateway = createTermTranslationPhaseGateway()
const termTranslationPhaseScreenControllerFactory =
  createTermTranslationPhaseScreenControllerFactory(termTranslationPhaseGateway)
const bodyTranslationPhaseGateway = createBodyTranslationPhaseGateway()
const bodyTranslationPhaseScreenControllerFactory =
  createBodyTranslationPhaseScreenControllerFactory(bodyTranslationPhaseGateway)
const personaGenerationPhaseGateway = createPersonaGenerationPhaseGateway()
const personaGenerationPhaseScreenControllerFactory =
  createPersonaGenerationPhaseScreenControllerFactory(
    personaGenerationPhaseGateway
  )
const providerSettingsGateway = createProviderSettingsGateway()
const providerSettingsScreenControllerFactory =
  createProviderSettingsScreenControllerFactory(providerSettingsGateway)
const translationInputGateway = createTranslationInputGateway()
const translationInputScreenControllerFactory =
  createTranslationInputScreenControllerFactory(translationInputGateway)
const translationJobSetupGateway = createTranslationJobSetupGateway()
const translationJobSetupScreenControllerFactory =
  createTranslationJobSetupScreenControllerFactory(translationJobSetupGateway)
const translationOutputArtifactGateway =
  createTranslationOutputArtifactGateway()
const translationOutputArtifactScreenControllerFactory =
  createTranslationOutputArtifactScreenControllerFactory(
    translationOutputArtifactGateway
  )

mount(App, {
  target,
  props: {
    createBodyTranslationPhaseScreenController:
      bodyTranslationPhaseScreenControllerFactory,
    createMasterDictionaryScreenController:
      masterDictionaryScreenControllerFactory,
    createPersonaGenerationPhaseScreenController:
      personaGenerationPhaseScreenControllerFactory,
    createProviderSettingsScreenController:
      providerSettingsScreenControllerFactory,
    createTermTranslationPhaseScreenController:
      termTranslationPhaseScreenControllerFactory,
    createTranslationInputScreenController:
      translationInputScreenControllerFactory,
    createTranslationOutputArtifactScreenController:
      translationOutputArtifactScreenControllerFactory,
    createTranslationJobSetupScreenController:
      translationJobSetupScreenControllerFactory
  }
})
