import { mount } from "svelte";
import App from "./App.svelte";
import "./app.css";
import { createBootstrapStatusScreenUsecase } from "@application/usecases/bootstrap-status";
import { createJobListScreenStore, createJobListScreenUsecase } from "@application/usecases/job-list";
import {
  createJobCreateScreenStore,
  createJobCreateScreenUsecase
} from "@application/usecases/job-create";
import { createTauriBootstrapStatusGateway } from "@gateway/tauri/bootstrap-status";
import { createBootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";

let previewJobCounter = 0;
const previewJobs = [
  {
    jobId: "job-observe-101",
    state: "Ready" as const
  },
  {
    jobId: "job-observe-202",
    state: "Running" as const
  }
];

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
const jobListStore = createJobListScreenStore();
const jobListUsecase = createJobListScreenUsecase({
  executor: async () => ({
    jobs: previewJobs
  }),
  store: jobListStore
});

mount(App, {
  props: {
    bootstrapStatusStore,
    bootstrapStatusUsecase,
    jobCreateStore,
    jobCreateUsecase,
    jobListStore,
    jobListUsecase
  },
  target: document.getElementById("app")!
});
