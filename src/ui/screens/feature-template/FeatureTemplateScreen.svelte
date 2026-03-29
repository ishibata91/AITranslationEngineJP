<script lang="ts">
  import { onMount } from "svelte";
  import type { FeatureTemplateScreenInput } from "@application/ports/input/feature-template";
  import type { FeatureTemplateScreenStore } from "@ui/stores/feature-template";
  import { FeatureTemplateView } from "@ui/views/feature-template";

  export let featureTemplateStore = undefined as unknown as FeatureTemplateScreenStore;
  export let featureTemplateUsecase = undefined as unknown as FeatureTemplateScreenInput;

  onMount(() => {
    void featureTemplateUsecase.initialize();
  });

  function handleRefresh(): void {
    void featureTemplateUsecase.refresh();
  }

  function handleRetry(): void {
    void featureTemplateUsecase.retry();
  }

  function handleSelect(event: CustomEvent<string>): void {
    const nextSelection = $featureTemplateStore.selection === event.detail ? null : event.detail;
    featureTemplateUsecase.select(nextSelection);
  }

  function handleUpdateQuery(event: CustomEvent<string>): void {
    void featureTemplateUsecase.updateFilters({
      query: event.detail
    }, {
      reload: true
    });
  }
</script>

<FeatureTemplateView
  state={$featureTemplateStore}
  on:refresh={handleRefresh}
  on:retry={handleRetry}
  on:select={handleSelect}
  on:updateQuery={handleUpdateQuery}
/>

