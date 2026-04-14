import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"
import { mount } from "svelte"
import App from "@ui/App.svelte"

const target = document.getElementById("app")

if (!target) {
  throw new Error("app root not found")
}

const masterDictionaryGateway = createMasterDictionaryGateway()
const masterDictionaryScreenControllerFactory =
  createMasterDictionaryScreenControllerFactory(masterDictionaryGateway)

mount(App, {
  target,
  props: {
    createMasterDictionaryScreenController:
      masterDictionaryScreenControllerFactory
  }
})
