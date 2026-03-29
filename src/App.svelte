<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import AppShell from "@ui/app-shell/AppShell.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
</script>

<AppShell {bootstrapStatusStore} {bootstrapStatusUsecase} />
