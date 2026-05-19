<script lang="ts">
  import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"
  import TranslationCompleteSummaryPanel from "./TranslationCompleteSummaryPanel.svelte"
  import TranslationResultListPanel from "./TranslationResultListPanel.svelte"

  interface Props {
    jobId: number
    rows: BodyTranslationFieldResultItem[]
  }

  let { jobId, rows }: Props = $props()
  let pageIndex = $state(0)

  const pageSize = 5
  const pageCount = $derived(Math.max(1, Math.ceil(rows.length / pageSize)))
  const currentPageRows = $derived(
    rows.slice(pageIndex * pageSize, pageIndex * pageSize + pageSize)
  )
  const pageLabel = $derived(`${pageIndex + 1} / ${pageCount}`)

  $effect(() => {
    if (pageIndex > pageCount - 1) {
      pageIndex = pageCount - 1
    }
  })

  function previousPage(): void {
    pageIndex = Math.max(0, pageIndex - 1)
  }

  function nextPage(): void {
    pageIndex = Math.min(pageCount - 1, pageIndex + 1)
  }
</script>

<section
  class="complete-shell"
  data-testid="translation-complete-translation-complete-screen"
  id="translationCompleteView"
>
  <TranslationCompleteSummaryPanel {jobId} />
  <TranslationResultListPanel
    rows={currentPageRows}
    {pageIndex}
    {pageCount}
    {pageLabel}
    onPreviousPage={previousPage}
    onNextPage={nextPage}
  />
</section>

<style>
  .complete-shell {
    display: grid;
    gap: 1rem;
  }
</style>
