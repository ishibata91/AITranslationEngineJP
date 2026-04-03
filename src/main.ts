import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { createBootstrapStatusScreenUsecase } from "@application/usecases/bootstrap-status";
import {
  createJobListScreenStore,
  createJobListScreenUsecase,
} from "@application/usecases/job-list";
import {
  createJobCreateScreenStore,
  createJobCreateScreenUsecase,
} from "@application/usecases/job-create";
import { createTauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status";
import { createTauriJobCreateExecutor } from "@gateway/tauri/job-create";
import { createTauriJobListExecutor } from "@gateway/tauri/job-list";
import { createBootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";

const bootstrapStatusStore = createBootstrapStatusScreenStore();
const bootstrapStatusUsecase = createBootstrapStatusScreenUsecase({
  gateway: createTauriBootstrapStatusGateway(),
  store: bootstrapStatusStore,
});
const jobCreateStore = createJobCreateScreenStore();
const jobCreateUsecase = createJobCreateScreenUsecase({
  executor: createTauriJobCreateExecutor(),
  store: jobCreateStore,
});
const jobListStore = createJobListScreenStore();
const jobListUsecase = createJobListScreenUsecase({
  executor: createTauriJobListExecutor(),
  store: jobListStore,
});

mount(App, {
  props: {
    bootstrapStatusStore,
    bootstrapStatusUsecase,
    jobCreateStore,
    jobCreateUsecase,
    jobListStore,
    jobListUsecase,
  },
  target: document.getElementById("app")!,
});
