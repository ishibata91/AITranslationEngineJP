<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type { JobCreateScreenInput, JobCreateScreenStore } from "@application/usecases/job-create";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import AppShell from "@ui/app-shell/AppShell.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;
</script>

<AppShell
  {bootstrapStatusStore}
  {bootstrapStatusUsecase}
  {jobCreateStore}
  {jobCreateUsecase}
/>
