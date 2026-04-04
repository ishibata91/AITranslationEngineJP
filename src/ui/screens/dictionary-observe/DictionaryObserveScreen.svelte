<script lang="ts">
  import { onMount } from "svelte";
  import type {
    DictionaryObserveScreenInput,
    DictionaryObserveScreenStore,
  } from "@application/usecases/dictionary-observe";
  import { DictionaryObserveView } from "@ui/views/dictionary-observe";

  export let dictionaryObserveStore =
    undefined as unknown as DictionaryObserveScreenStore;
  export let dictionaryObserveUsecase =
    undefined as unknown as DictionaryObserveScreenInput;

  onMount(() => {
    void dictionaryObserveUsecase.initialize();
  });
</script>

<DictionaryObserveView
  state={$dictionaryObserveStore}
  on:observe={() => void dictionaryObserveUsecase.observe()}
  on:refresh={() => void dictionaryObserveUsecase.refresh()}
  on:retry={() => void dictionaryObserveUsecase.retry()}
  on:select={(event) => dictionaryObserveUsecase.select(event.detail)}
  on:updateSourceTexts={(event) =>
    void dictionaryObserveUsecase.updateFilters({
      ...$dictionaryObserveStore.filters,
      sourceTexts: event.detail,
    })}
/>
