<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type {
    DictionaryObserveScreenInput,
    DictionaryObserveScreenStore,
  } from "@application/usecases/dictionary-observe";
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
  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;
  export let jobListStore = undefined as unknown as JobListScreenStore;
  export let jobListUsecase = undefined as unknown as JobListScreenInput;
</script>

<AppShell
  {bootstrapStatusStore}
  {bootstrapStatusUsecase}
  {dictionaryObserveStore}
  {dictionaryObserveUsecase}
  {jobCreateStore}
  {jobCreateUsecase}
  {jobListStore}
  {jobListUsecase}
/>
