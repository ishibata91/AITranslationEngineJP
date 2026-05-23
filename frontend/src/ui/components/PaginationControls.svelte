<script lang="ts">
  import ActionButton from "./ActionButton.svelte"

  interface Props {
    page: number
    pageCount: number
    totalLabel?: string
    disabled?: boolean
    busy?: boolean
    onPrevious: () => void
    onNext: () => void
  }

  let {
    page,
    pageCount,
    totalLabel = "",
    disabled = false,
    busy = false,
    onPrevious,
    onNext
  }: Props = $props()

  const currentPage = $derived(
    Math.min(Math.max(page, 1), Math.max(pageCount, 1))
  )
  const safePageCount = $derived(Math.max(pageCount, 1))
  const previousDisabled = $derived(disabled || busy || currentPage <= 1)
  const nextDisabled = $derived(
    disabled || busy || currentPage >= safePageCount
  )
</script>

<nav class="pagination-controls" aria-label="ページ移動">
  <ActionButton
    label="前へ"
    variant="secondary"
    disabled={previousDisabled}
    onClick={onPrevious}
  />
  <p aria-live="polite">
    <span>{currentPage} / {safePageCount}</span>
    {#if totalLabel}
      <span>{totalLabel}</span>
    {/if}
  </p>
  <ActionButton
    label="次へ"
    variant="secondary"
    disabled={nextDisabled}
    {busy}
    onClick={onNext}
  />
</nav>

<style>
  .pagination-controls {
    align-items: center;
    display: flex;
    gap: 0.8rem;
    justify-content: flex-end;
    min-width: 0;
  }

  p {
    color: #334155;
    display: grid;
    font-weight: 700;
    gap: 0.15rem;
    margin: 0;
    min-width: 6rem;
    text-align: center;
  }

  span {
    overflow-wrap: anywhere;
  }

  @media (max-width: 520px) {
    .pagination-controls {
      align-items: stretch;
      flex-direction: column;
    }

    p {
      text-align: left;
    }
  }
</style>
