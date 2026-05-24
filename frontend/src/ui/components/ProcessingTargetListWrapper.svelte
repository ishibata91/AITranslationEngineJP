<script lang="ts">
  import type { Snippet } from "svelte"

  import ProcessingTargetListPanel from "./ProcessingTargetListPanel.svelte"
  import type { ProcessingTargetListItem } from "./processing-target-list-panel-types"

  export interface ProcessingTargetListFilterOption {
    value: string
    label: string
  }

  interface Props {
    title: string
    titleId: string
    items: ProcessingTargetListItem[]
    searchId?: string
    searchLabel?: string
    searchValue?: string
    countText?: string
    eyebrow?: string
    filterId?: string
    filterLabel?: string
    filterOptions?: ProcessingTargetListFilterOption[]
    filterValue?: string
    initialExpandedItemId?: string | null
    pageSize?: number
    currentPage?: number
    totalCount?: number
    busy?: boolean
    searchPlaceholder?: string
    supportText?: string
    actions?: Snippet
    footer?: Snippet
    summary?: Snippet
    onFilterChange?: (event: Event) => void
    onSearchInput?: (event: Event) => void
    onPreviousPage?: () => void
    onNextPage?: () => void
    onPageChange?: (page: number) => void
    onSelectItem?: (itemId: string) => void
  }

  let {
    title,
    titleId,
    items,
    searchId = "",
    searchLabel = "",
    searchValue = "",
    countText = "",
    eyebrow = "",
    filterId = "",
    filterLabel = "",
    filterOptions = [],
    filterValue = "",
    initialExpandedItemId = null,
    pageSize = 50,
    currentPage = undefined,
    totalCount = undefined,
    busy = false,
    searchPlaceholder = "",
    supportText = "",
    actions,
    footer,
    summary,
    onFilterChange,
    onSearchInput,
    onPreviousPage,
    onNextPage,
    onPageChange,
    onSelectItem
  }: Props = $props()

  const hasSearch = $derived(Boolean(searchId && searchLabel && onSearchInput))
  const hasFilter = $derived(
    Boolean(filterId && filterLabel && filterOptions.length > 0 && onFilterChange)
  )
</script>

<section class="target-list-wrapper" aria-labelledby={titleId}>
  <div class="wrapper-head">
    <div class="heading-copy">
      {#if eyebrow}
        <p class="eyebrow">{eyebrow}</p>
      {/if}
      <h3 id={titleId}>{title}</h3>
      {#if supportText}
        <p class="support-copy">{supportText}</p>
      {/if}
    </div>

    <div class="head-actions">
      {#if countText}
        <span class="status-pill">{countText}</span>
      {/if}
      {@render actions?.()}
    </div>
  </div>

  {#if hasSearch || hasFilter}
    <div class:has-filter={hasFilter} class="filter-grid">
      {#if hasSearch}
        <label class="field-group" for={searchId}>
          <span class="field-label">{searchLabel}</span>
          <input
            class="text-field"
            id={searchId}
            oninput={onSearchInput}
            placeholder={searchPlaceholder}
            type="search"
            value={searchValue}
          />
        </label>
      {/if}

      {#if hasFilter}
        <label class="field-group" for={filterId}>
          <span class="field-label">{filterLabel}</span>
          <select
            class="text-field"
            id={filterId}
            onchange={onFilterChange}
            value={filterValue}
          >
            {#each filterOptions as option (option.value)}
              <option value={option.value}>{option.label}</option>
            {/each}
          </select>
        </label>
      {/if}
    </div>
  {/if}

  {@render summary?.()}

  <ProcessingTargetListPanel
    {items}
    {currentPage}
    {initialExpandedItemId}
    {pageSize}
    {totalCount}
    {busy}
    showHeadingTitle={false}
    {onPreviousPage}
    {onNextPage}
    {onPageChange}
    {onSelectItem}
  />

  {@render footer?.()}
</section>

<style>
  .target-list-wrapper {
    background: rgba(17, 13, 12, 0.42);
    border: 0.5px solid var(--line);
    border-radius: 20px;
    box-shadow: var(--shadow);
    color: var(--text);
    display: grid;
    gap: 12px;
    min-width: 0;
    padding: 16px;
  }

  .wrapper-head,
  .head-actions {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: space-between;
    min-width: 0;
  }

  .heading-copy {
    display: grid;
    gap: 4px;
    min-width: 0;
  }

  .head-actions {
    justify-content: flex-end;
  }

  .eyebrow,
  .field-label {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  h3,
  p,
  .support-copy {
    margin: 0;
    overflow-wrap: anywhere;
  }

  .support-copy {
    color: var(--muted);
    line-height: 1.7;
  }

  .status-pill {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    display: inline-flex;
    min-height: 28px;
    padding: 0 10px;
  }

  .filter-grid {
    align-items: end;
    display: grid;
    gap: 8px;
    grid-template-columns: minmax(0, 1fr);
  }

  .filter-grid.has-filter {
    grid-template-columns: minmax(0, 1.5fr) minmax(220px, 0.7fr);
  }

  .field-group {
    display: grid;
    gap: 4px;
    min-width: 0;
  }

  .text-field {
    appearance: none;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    border-radius: 10px;
    color: var(--text);
    min-height: 38px;
    min-width: 0;
    padding: 0 10px;
    width: 100%;
  }

  @media (max-width: 760px) {
    .filter-grid.has-filter {
      grid-template-columns: 1fr;
    }

    .head-actions {
      justify-content: flex-start;
      width: 100%;
    }
  }
</style>
