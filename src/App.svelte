<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type {
    DictionaryObserveScreenInput,
    DictionaryObserveScreenStore,
  } from "@application/usecases/dictionary-observe";
  import type {
    PersonaObserveScreenInput,
    PersonaObserveScreenStore,
  } from "@application/usecases/persona-observe";
  import type {
    ExecutionControlScreenInput,
    ExecutionControlScreenStore,
  } from "@application/usecases/execution-control";
  import type {
    ExecutionObserveScreenInput,
    ExecutionObserveScreenStore,
  } from "@application/usecases/execution-observe";
  import type { JobListScreenInput, JobListScreenStore } from "@application/usecases/job-list";
  import type { JobCreateScreenInput, JobCreateScreenStore } from "@application/usecases/job-create";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import AppShell from "@ui/app-shell/AppShell.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
  export let dictionaryObserveStore =
    undefined as unknown as DictionaryObserveScreenStore | undefined;
  export let dictionaryObserveUsecase =
    undefined as unknown as DictionaryObserveScreenInput | undefined;
  export let executionControlStore =
    undefined as unknown as ExecutionControlScreenStore | undefined;
  export let executionControlUsecase =
    undefined as unknown as ExecutionControlScreenInput | undefined;
  export let executionObserveStore =
    undefined as unknown as ExecutionObserveScreenStore | undefined;
  export let executionObserveUsecase =
    undefined as unknown as ExecutionObserveScreenInput | undefined;
  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;
  export let jobListStore = undefined as unknown as JobListScreenStore;
  export let jobListUsecase = undefined as unknown as JobListScreenInput;
  export let personaObserveStore =
    undefined as unknown as PersonaObserveScreenStore | undefined;
  export let personaObserveUsecase =
    undefined as unknown as PersonaObserveScreenInput | undefined;
</script>

<AppShell
  {bootstrapStatusStore}
  {bootstrapStatusUsecase}
  {dictionaryObserveStore}
  {dictionaryObserveUsecase}
  {executionControlStore}
  {executionControlUsecase}
  {executionObserveStore}
  {executionObserveUsecase}
  {jobCreateStore}
  {jobCreateUsecase}
  {jobListStore}
  {jobListUsecase}
  {personaObserveStore}
  {personaObserveUsecase}
/>
