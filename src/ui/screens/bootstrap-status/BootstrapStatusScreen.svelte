<script lang="ts">
  import { onMount } from "svelte";
  import type {
    BootstrapStatusField,
    BootstrapStatusScreenInput
  } from "@application/ports/input/bootstrap-status";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import { BootstrapStatusView } from "@ui/views/bootstrap-status";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;

  onMount(() => {
    void bootstrapStatusUsecase.initialize();
  });

  function handleRefresh(): void {
    void bootstrapStatusUsecase.refresh();
  }

  function handleRetry(): void {
    void bootstrapStatusUsecase.retry();
  }

  function handleSelect(event: CustomEvent<BootstrapStatusField>): void {
    const nextSelection =
      $bootstrapStatusStore.selection === event.detail ? null : event.detail;

    bootstrapStatusUsecase.select(nextSelection);
  }
</script>

<BootstrapStatusView
  state={$bootstrapStatusStore}
  on:refresh={handleRefresh}
  on:retry={handleRetry}
  on:select={handleSelect}
/>
