import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { mount } from "svelte"
import App from "@ui/App.svelte"

const target = document.getElementById("app")

if (!target) {
  throw new Error("app root not found")
}

const masterDictionaryGateway = createMasterDictionaryGateway()

mount(App, {
  target,
  props: {
    masterDictionaryGateway
  }
})
