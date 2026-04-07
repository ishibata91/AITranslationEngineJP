<script lang="ts">
  import { onMount } from "svelte";
  import type {
    ExecutionObserveScreenInput,
    ExecutionObserveScreenStore,
  } from "@application/usecases/execution-observe";
  import { ExecutionObserveView } from "@ui/views/execution-observe";

  export let executionObserveStore =
    undefined as unknown as ExecutionObserveScreenStore;
  export let executionObserveUsecase =
    undefined as unknown as ExecutionObserveScreenInput;

  onMount(() => {
    void executionObserveUsecase.initialize();
  });
</script>

<ExecutionObserveView
  state={$executionObserveStore}
  on:refresh={() => void executionObserveUsecase.refresh()}
/>
