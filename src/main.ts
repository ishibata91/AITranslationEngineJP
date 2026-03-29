import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { createBootstrapStatusScreenUsecase } from "@application/usecases/bootstrap-status";
import { createTauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status";
import { createBootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";

const bootstrapStatusStore = createBootstrapStatusScreenStore();
const bootstrapStatusUsecase = createBootstrapStatusScreenUsecase({
  gateway: createTauriBootstrapStatusGateway(),
  store: bootstrapStatusStore
});

mount(App, {
  props: {
    bootstrapStatusStore,
    bootstrapStatusUsecase
  },
  target: document.getElementById("app")!
});
