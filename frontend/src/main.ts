import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"
import { createTermTranslationPhaseScreenControllerFactory } from "@controller/term-translation-phase"
import { createTermTranslationPhaseGateway } from "@controller/wails/term-translation-phase.gateway"
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

mount(App, {
  target,
  props: {
    createMasterDictionaryScreenController:
      masterDictionaryScreenControllerFactory,
    createTermTranslationPhaseScreenController:
      termTranslationPhaseScreenControllerFactory
  }
})
