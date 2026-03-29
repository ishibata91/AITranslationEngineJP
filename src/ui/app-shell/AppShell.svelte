<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import BootstrapStatusScreen from "@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
</script>

<svelte:head>
  <title>AITranslationEngineJp Bootstrap</title>
</svelte:head>

<main class="shell">
  <BootstrapStatusScreen {bootstrapStatusStore} {bootstrapStatusUsecase} />
</main>

<style>
  .shell {
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: 2rem;
  }
</style>
