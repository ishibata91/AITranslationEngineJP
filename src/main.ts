import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { createBootstrapStatusScreenUsecase } from "@application/usecases/bootstrap-status";
import {
  createJobCreateScreenStore,
  createJobCreateScreenUsecase
} from "@application/usecases/job-create";
import { createTauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status";
import { createBootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";

let previewJobCounter = 0;

const bootstrapStatusStore = createBootstrapStatusScreenStore();
const bootstrapStatusUsecase = createBootstrapStatusScreenUsecase({
  gateway: createTauriBootstrapStatusGateway(),
  store: bootstrapStatusStore
});
const jobCreateStore = createJobCreateScreenStore();
const jobCreateUsecase = createJobCreateScreenUsecase({
  executor: async () => {
    previewJobCounter += 1;

    return {
      jobId: `job-preview-${previewJobCounter}`,
      state: "Ready"
    };
  },
  store: jobCreateStore
});

mount(App, {
  props: {
    bootstrapStatusStore,
    bootstrapStatusUsecase,
    jobCreateStore,
    jobCreateUsecase
  },
  target: document.getElementById("app")!
});
