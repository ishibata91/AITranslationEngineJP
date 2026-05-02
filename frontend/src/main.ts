import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"
import { createPersonaGenerationPhaseScreenControllerFactory } from "@controller/persona-generation-phase"
import { createTermTranslationPhaseScreenControllerFactory } from "@controller/term-translation-phase"
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
  createTermTranslationPhaseScreenControllerFactory(
    termTranslationPhaseGateway
  )
const personaGenerationPhaseGateway = createPersonaGenerationPhaseGateway()
const personaGenerationPhaseScreenControllerFactory =
  createPersonaGenerationPhaseScreenControllerFactory(
    personaGenerationPhaseGateway
  )

mount(App, {
  target,
  props: {
    createMasterDictionaryScreenController:
      masterDictionaryScreenControllerFactory,
    createPersonaGenerationPhaseScreenController:
      personaGenerationPhaseScreenControllerFactory,
    createTermTranslationPhaseScreenController:
      termTranslationPhaseScreenControllerFactory
  }
})
