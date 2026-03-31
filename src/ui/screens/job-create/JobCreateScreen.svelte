<script lang="ts">
  import { onMount } from "svelte";
  import type { JobCreateScreenInput, JobCreateScreenStore } from "@application/usecases/job-create";
  import { JobCreateView } from "@ui/views/job-create";

  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;

  onMount(() => {
    void jobCreateUsecase.initialize();
  });
</script>

<JobCreateView
  state={$jobCreateStore}
  on:resetResult={() => jobCreateUsecase.resetResult()}
  on:submit={() => void jobCreateUsecase.submit()}
  on:updateSourceGroupField={(event) =>
    jobCreateUsecase.updateSourceGroupField(
      event.detail.groupIndex,
      event.detail.field,
      event.detail.value
    )}
  on:updateTranslationUnitField={(event) =>
    jobCreateUsecase.updateTranslationUnitField(
      event.detail.groupIndex,
      event.detail.unitIndex,
      event.detail.field,
      event.detail.value
    )}
/>
