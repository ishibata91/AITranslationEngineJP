<script lang="ts">
  import type { BootstrapStatusScreenInput } from "@application/ports/input/bootstrap-status";
  import type { JobListScreenInput, JobListScreenStore } from "@application/usecases/job-list";
  import type { JobCreateScreenInput, JobCreateScreenStore } from "@application/usecases/job-create";
  import type {
    DictionaryObserveScreenInput,
    DictionaryObserveScreenStore,
  } from "@application/usecases/dictionary-observe";
  import type {
    PersonaObserveScreenInput,
    PersonaObserveScreenStore,
  } from "@application/usecases/persona-observe";
  import type {
    ExecutionControlScreenInput,
    ExecutionControlScreenStore,
  } from "@application/usecases/execution-control";
  import type {
    ExecutionObserveScreenInput,
    ExecutionObserveScreenStore,
  } from "@application/usecases/execution-observe";
  import type {
    TranslationPreviewScreenInput,
    TranslationPreviewScreenStore,
  } from "@application/usecases/translation-preview";
  import type { BootstrapStatusScreenStore } from "@ui/stores/bootstrap-status";
  import BootstrapStatusScreen from "@ui/screens/bootstrap-status/BootstrapStatusScreen.svelte";
  import DictionaryObserveScreen from "@ui/screens/dictionary-observe/DictionaryObserveScreen.svelte";
  import JobCreateScreen from "@ui/screens/job-create/JobCreateScreen.svelte";
  import JobListScreen from "@ui/screens/job-list/JobListScreen.svelte";
  import PersonaObserveScreen from "@ui/screens/persona-observe/PersonaObserveScreen.svelte";
  import ExecutionControlScreen from "@ui/screens/execution-control/ExecutionControlScreen.svelte";
  import ExecutionObserveScreen from "@ui/screens/execution-observe/ExecutionObserveScreen.svelte";
  import TranslationPreviewScreen from "@ui/screens/translation-preview/TranslationPreviewScreen.svelte";

  export let bootstrapStatusStore = undefined as unknown as BootstrapStatusScreenStore;
  export let bootstrapStatusUsecase = undefined as unknown as BootstrapStatusScreenInput;
  export let dictionaryObserveStore =
    undefined as unknown as DictionaryObserveScreenStore | undefined;
  export let dictionaryObserveUsecase =
    undefined as unknown as DictionaryObserveScreenInput | undefined;
  export let executionControlStore =
    undefined as unknown as ExecutionControlScreenStore | undefined;
  export let executionControlUsecase =
    undefined as unknown as ExecutionControlScreenInput | undefined;
  export let executionObserveStore =
    undefined as unknown as ExecutionObserveScreenStore | undefined;
  export let executionObserveUsecase =
    undefined as unknown as ExecutionObserveScreenInput | undefined;
  export let jobCreateStore = undefined as unknown as JobCreateScreenStore;
  export let jobCreateUsecase = undefined as unknown as JobCreateScreenInput;
  export let jobListStore = undefined as unknown as JobListScreenStore;
  export let jobListUsecase = undefined as unknown as JobListScreenInput;
  export let personaObserveStore =
    undefined as unknown as PersonaObserveScreenStore | undefined;
  export let personaObserveUsecase =
    undefined as unknown as PersonaObserveScreenInput | undefined;
  export let translationPreviewStore =
    undefined as unknown as TranslationPreviewScreenStore | undefined;
  export let translationPreviewUsecase =
    undefined as unknown as TranslationPreviewScreenInput | undefined;
</script>

<svelte:head>
  <title>AITranslationEngineJp Phase 1 Screens</title>
</svelte:head>

<main class="shell">
  {#if dictionaryObserveStore !== undefined && dictionaryObserveUsecase !== undefined}
    <DictionaryObserveScreen {dictionaryObserveStore} {dictionaryObserveUsecase} />
  {/if}
  {#if personaObserveStore !== undefined && personaObserveUsecase !== undefined}
    <PersonaObserveScreen {personaObserveStore} {personaObserveUsecase} />
  {/if}
  {#if executionControlStore !== undefined && executionControlUsecase !== undefined}
    <ExecutionControlScreen {executionControlStore} {executionControlUsecase} />
  {/if}
  {#if executionObserveStore !== undefined && executionObserveUsecase !== undefined}
    <ExecutionObserveScreen {executionObserveStore} {executionObserveUsecase} />
  {/if}
  {#if translationPreviewStore !== undefined && translationPreviewUsecase !== undefined}
    <TranslationPreviewScreen {translationPreviewStore} {translationPreviewUsecase} />
  {/if}
  <JobListScreen {jobListStore} {jobListUsecase} />
  <JobCreateScreen {jobCreateStore} {jobCreateUsecase} />
  <BootstrapStatusScreen {bootstrapStatusStore} {bootstrapStatusUsecase} />
</main>

<style>
  .shell {
    min-height: 100vh;
    display: grid;
    place-items: center;
    gap: 1.5rem;
    padding: 2rem;
  }
</style>
