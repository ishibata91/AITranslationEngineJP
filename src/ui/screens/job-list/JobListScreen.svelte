<script lang="ts">
  import { onMount } from "svelte";
  import type { JobListScreenInput, JobListScreenStore } from "@application/usecases/job-list";
  import { JobListView } from "@ui/views/job-list";

  export let jobListStore = undefined as unknown as JobListScreenStore;
  export let jobListUsecase = undefined as unknown as JobListScreenInput;

  onMount(() => {
    void jobListUsecase.initialize();
  });
</script>

<JobListView
  state={$jobListStore}
  on:refresh={() => void jobListUsecase.refresh()}
  on:retry={() => void jobListUsecase.retry()}
  on:select={(event) => jobListUsecase.select(event.detail)}
/>
