<script lang="ts">
  import { onMount } from "svelte";
  import type {
    PersonaObserveScreenInput,
    PersonaObserveScreenStore,
  } from "@application/usecases/persona-observe";
  import { PersonaObserveView } from "@ui/views/persona-observe";

  export let personaObserveStore =
    undefined as unknown as PersonaObserveScreenStore;
  export let personaObserveUsecase =
    undefined as unknown as PersonaObserveScreenInput;

  onMount(() => {
    void personaObserveUsecase.initialize();
  });
</script>

<PersonaObserveView
  state={$personaObserveStore}
  on:observe={() => void personaObserveUsecase.observe()}
  on:refresh={() => void personaObserveUsecase.refresh()}
  on:retry={() => void personaObserveUsecase.retry()}
  on:select={(event) => personaObserveUsecase.select(event.detail)}
  on:updatePersonaName={(event) =>
    void personaObserveUsecase.updateFilters({
      ...$personaObserveStore.filters,
      personaName: event.detail,
    })}
/>
