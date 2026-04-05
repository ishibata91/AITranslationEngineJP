<script lang="ts">
  import { onMount } from "svelte";
  import type {
    TranslationPreviewScreenInput,
    TranslationPreviewScreenStore,
  } from "@application/usecases/translation-preview";
  import { TranslationPreviewView } from "@ui/views/translation-preview";

  export let translationPreviewStore =
    undefined as unknown as TranslationPreviewScreenStore;
  export let translationPreviewUsecase =
    undefined as unknown as TranslationPreviewScreenInput;

  onMount(() => {
    void translationPreviewUsecase.initialize();
  });
</script>

<TranslationPreviewView
  state={$translationPreviewStore}
  on:observe={() => void translationPreviewUsecase.observe()}
  on:refresh={() => void translationPreviewUsecase.refresh()}
  on:retry={() => void translationPreviewUsecase.retry()}
  on:select={(event) => translationPreviewUsecase.select(event.detail)}
  on:updateJobId={(event) =>
    void translationPreviewUsecase.updateFilters({
      ...$translationPreviewStore.filters,
      jobId: event.detail,
    })}
/>
