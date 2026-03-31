<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type { JobCreateScreenInput, JobCreateScreenStore } from "@application/usecases/job-create";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import BootstrapStatusScreen from "@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte";
  import JobCreateScreen from "@ui/screens/job-create/JobCreateScreen.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;
</script>

<svelte:head>
  <title>AITranslationEngineJp Job Create</title>
</svelte:head>

<main class="shell">
  <JobCreateScreen {jobCreateStore} {jobCreateUsecase} />
  <BootstrapStatusScreen {bootstrapStatusStore} {bootstrapStatusUsecase} />
</main>

<style>
  .shell {
    min-height: 100vh;
    display: grid;
    place-items: center;
    gap: 1.5rem;
    padding: 2rem;
  }
</style>
