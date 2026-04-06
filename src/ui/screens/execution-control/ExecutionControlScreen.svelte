<script lang="ts">
  import { onMount } from "svelte";
  import type {
    ExecutionControlScreenInput,
    ExecutionControlScreenStore,
  } from "@application/usecases/execution-control";
  import { ExecutionControlView } from "@ui/views/execution-control";

  export let executionControlStore =
    undefined as unknown as ExecutionControlScreenStore;
  export let executionControlUsecase =
    undefined as unknown as ExecutionControlScreenInput;

  onMount(() => {
    void executionControlUsecase.initialize();
  });
</script>

<ExecutionControlView
  state={$executionControlStore}
  on:cancel={() => void executionControlUsecase.cancel()}
  on:pause={() => void executionControlUsecase.pause()}
  on:resume={() => void executionControlUsecase.resume()}
  on:retry={() => void executionControlUsecase.retry()}
/>
