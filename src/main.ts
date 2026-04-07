import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { createBootstrapStatusScreenUsecase } from "@application/usecases/bootstrap-status";
import {
  createDictionaryObserveScreenStore,
  createDictionaryObserveScreenUsecase,
} from "@application/usecases/dictionary-observe";
import {
  createJobListScreenStore,
  createJobListScreenUsecase,
} from "@application/usecases/job-list";
import {
  createJobCreateScreenStore,
  createJobCreateScreenUsecase,
} from "@application/usecases/job-create";
import {
  createPersonaObserveScreenStore,
  createPersonaObserveScreenUsecase,
} from "@application/usecases/persona-observe";
import {
  createExecutionControlScreenStore,
  createExecutionControlScreenUsecase,
} from "@application/usecases/execution-control";
import {
  createExecutionObserveScreenStore,
  createExecutionObserveScreenUsecase,
} from "@application/usecases/execution-observe";
import { createTauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status";
import { createTauriDictionaryObserveExecutor } from "@gateway/tauri/dictionary-observe";
import { createTauriExecutionControlGateway } from "@gateway/tauri/execution-control";
import { createTauriExecutionObserveLoader } from "@gateway/tauri/execution-observe";
import { createTauriJobCreateExecutor } from "@gateway/tauri/job-create";
import { createTauriJobListExecutor } from "@gateway/tauri/job-list";
import { createTauriPersonaObserveExecutor } from "@gateway/tauri/persona-observe";
import { createBootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";

const bootstrapStatusStore = createBootstrapStatusScreenStore();
const bootstrapStatusUsecase = createBootstrapStatusScreenUsecase({
  gateway: createTauriBootstrapStatusGateway(),
  store: bootstrapStatusStore,
});
const dictionaryObserveStore = createDictionaryObserveScreenStore();
const dictionaryObserveUsecase = createDictionaryObserveScreenUsecase({
  executor: createTauriDictionaryObserveExecutor(),
  store: dictionaryObserveStore,
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
const personaObserveStore = createPersonaObserveScreenStore();
const personaObserveUsecase = createPersonaObserveScreenUsecase({
  executor: createTauriPersonaObserveExecutor(),
  store: personaObserveStore,
});
const executionControlStore = createExecutionControlScreenStore();
const executionControlUsecase = createExecutionControlScreenUsecase({
  ...createTauriExecutionControlGateway(),
  store: executionControlStore,
});
const executionObserveStore = createExecutionObserveScreenStore();
const executionObserveUsecase = createExecutionObserveScreenUsecase({
  loadSnapshot: createTauriExecutionObserveLoader(),
  store: executionObserveStore,
});

mount(App, {
  props: {
    bootstrapStatusStore,
    bootstrapStatusUsecase,
    dictionaryObserveStore,
    dictionaryObserveUsecase,
    executionControlStore,
    executionControlUsecase,
    executionObserveStore,
    executionObserveUsecase,
    jobCreateStore,
    jobCreateUsecase,
    jobListStore,
    jobListUsecase,
    personaObserveStore,
    personaObserveUsecase,
  },
  target: document.getElementById("app")!,
});
