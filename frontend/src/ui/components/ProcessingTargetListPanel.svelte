<script lang="ts">
  import PaginationControls from "@ui/components/PaginationControls.svelte"

  import ProcessingTargetTitleTooltip from "./ProcessingTargetTitleTooltip.svelte"
  import type { ProcessingTargetListItem } from "./processing-target-list-panel-types"

  const DEFAULT_PAGE_SIZE = 50

  interface Props {
    items: ProcessingTargetListItem[]
    pageSize?: number
    initialPage?: number
    currentPage?: number
    totalCount?: number
    busy?: boolean
    initialExpandedItemId?: string | null
    showHeadingTitle?: boolean
    titleColumnLabels?: string[]
    onPreviousPage?: () => void
    onNextPage?: () => void
    onPageChange?: (page: number) => void
    onSelectItem?: (itemId: string) => void
    rowTestId?: string
    totalCountTestId?: string
    emptyStateTestId?: string
  }

  let {
    items,
    pageSize = DEFAULT_PAGE_SIZE,
    initialPage: providedInitialPage = 1,
    currentPage: controlledPage = undefined,
    totalCount: controlledTotalCount = undefined,
    busy = false,
    initialExpandedItemId = null,
    showHeadingTitle = true,
    titleColumnLabels: providedTitleColumnLabels = undefined,
    onPreviousPage,
    onNextPage,
    onPageChange,
    onSelectItem,
    rowTestId = undefined,
    totalCountTestId = undefined,
    emptyStateTestId = undefined
  }: Props = $props()

  let currentPage = $state(1)
  let appliedInitialPage = $state<number | null>(null)
  let appliedInitialExpandedItemId = $state<string | null>(null)
  let expandedItemId = $state<string | null>(null)

  const safePageSize = $derived(Math.max(1, pageSize))
  const isControlledPage = $derived(
    controlledPage !== undefined || controlledTotalCount !== undefined
  )
  const itemCount = $derived(
    controlledTotalCount !== undefined
      ? Math.max(0, controlledTotalCount)
      : items.length
  )
  const pageCount = $derived(Math.max(1, Math.ceil(itemCount / safePageSize)))
  const requestedPage = $derived(
    controlledPage !== undefined ? controlledPage : currentPage
  )
  const safeCurrentPage = $derived(
    Math.min(Math.max(requestedPage, 1), pageCount)
  )
  const startIndex = $derived((safeCurrentPage - 1) * safePageSize)
  const endIndex = $derived(Math.min(startIndex + safePageSize, itemCount))
  const visibleItems = $derived(
    isControlledPage ? items : items.slice(startIndex, endIndex)
  )
  const expandedItem = $derived(
    visibleItems.find((item) => item.id === expandedItemId) ?? null
  )
  const pageRangeLabel = $derived(
    itemCount === 0 ? "0 件" : `${startIndex + 1}-${endIndex} / ${itemCount} 件`
  )

  $effect(() => {
    if (currentPage > pageCount) {
      currentPage = pageCount
    }
  })

  $effect(() => {
    if (appliedInitialPage !== providedInitialPage) {
      currentPage = Math.min(Math.max(providedInitialPage, 1), pageCount)
      appliedInitialPage = providedInitialPage
    }
  })

  $effect(() => {
    if (appliedInitialExpandedItemId !== initialExpandedItemId) {
      expandedItemId = initialExpandedItemId
      appliedInitialExpandedItemId = initialExpandedItemId
    }
  })

  $effect(() => {
    if (expandedItemId && !expandedItem) {
      expandedItemId = null
    }
  })

  function toggleItem(itemId: string): void {
    expandedItemId = expandedItemId === itemId ? null : itemId
    onSelectItem?.(itemId)
  }

  function handleRowKeydown(event: KeyboardEvent, itemId: string): void {
    if (event.key !== "Enter" && event.key !== " ") {
      return
    }

    event.preventDefault()
    toggleItem(itemId)
  }

  function showPreviousPage(): void {
    const previousPage = Math.max(1, safeCurrentPage - 1)
    onPreviousPage?.()
    onPageChange?.(previousPage)
    if (!isControlledPage) {
      currentPage = previousPage
    }
  }

  function showNextPage(): void {
    const nextPage = Math.min(pageCount, safeCurrentPage + 1)
    onNextPage?.()
    onPageChange?.(nextPage)
    if (!isControlledPage) {
      currentPage = nextPage
    }
  }

  function getItemTitleParts(item: ProcessingTargetListItem): string[] {
    return item.titleParts?.map((part) => part.text) ?? [item.name]
  }

  function getItemTitle(item: ProcessingTargetListItem): string {
    const nonBlankTitleParts = getItemTitleParts(item).filter(
      (titlePart) => titlePart.trim() !== ""
    )
    return nonBlankTitleParts.length > 0
      ? nonBlankTitleParts.join(" / ")
      : item.name
  }

  function parseTitlePart(
    titlePart: string,
    fallbackLabel?: string
  ): { label: string; value: string } {
    const match = titlePart.match(/^(.+?)(?:\s+\d+)?[:：]\s*(.+)$/)
    if (!match) {
      return { label: fallbackLabel ?? "対象", value: titlePart }
    }

    return {
      label: (match[1] ?? fallbackLabel ?? "対象").trim(),
      value: (match[2] ?? titlePart).trim()
    }
  }

  const visibleTitleColumnCount = $derived(
    Math.max(1, ...visibleItems.map((item) => getItemTitleParts(item).length))
  )
  const titleColumnCount = $derived(
    providedTitleColumnLabels?.length
      ? providedTitleColumnLabels.length
      : visibleTitleColumnCount
  )
  const titleColumnIndexes = $derived(
    Array.from({ length: titleColumnCount }, (_, index) => index)
  )
  const COLUMN_FALLBACK_LABELS = ["対象 1", "対象 2", "対象 3"] as const
  const derivedTitleColumnLabels = $derived(
    titleColumnCount === 1
      ? ["対象"]
      : (
          visibleItems
            .map((item) => getItemTitleParts(item))
            .find((titleParts) => titleParts.length > 1) ?? []
        )
          .slice(0, titleColumnCount)
          .map((titlePart, i) =>
            parseTitlePart(
              titlePart,
              COLUMN_FALLBACK_LABELS[i] ?? `対象 ${i + 1}`
            ).label
          )
  )
  const resolvedTitleColumnLabels = $derived(
    providedTitleColumnLabels ?? derivedTitleColumnLabels
  )
</script>

<section class="processing-target-list-panel" aria-label="処理対象一覧">
  <div class="panel-heading">
    {#if showHeadingTitle}
      <h3>処理対象一覧</h3>
    {/if}
    <p class="target-count" data-testid={totalCountTestId}>{pageRangeLabel}</p>
  </div>

  <div class="target-table-frame">
    <table
      class="target-table"
      class:has-title-columns={titleColumnCount > 1}
      aria-label="処理対象ページ"
    >
      <thead>
        <tr>
          {#each resolvedTitleColumnLabels as titleColumnLabel, index (`${index}-${titleColumnLabel}`)}
            <th scope="col">{titleColumnLabel}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#if visibleItems.length === 0}
          <tr class="target-empty-row">
            <td colspan={titleColumnCount} data-testid={emptyStateTestId}>
              処理対象がありません
            </td>
          </tr>
        {:else}
          {#each visibleItems as item (item.id)}
            {@const isExpanded = expandedItemId === item.id}
            {@const itemTitleParts = getItemTitleParts(item)}
            {@const itemTitle = getItemTitle(item)}
            <tr
              class:expanded={isExpanded}
              tabindex="0"
              aria-expanded={isExpanded}
              aria-controls={`processing-target-detail-${item.id}`}
              aria-label={itemTitle}
              data-item-id={item.id}
              data-testid={rowTestId}
              onclick={() => toggleItem(item.id)}
              onkeydown={(event) => handleRowKeydown(event, item.id)}
            >
              {#if titleColumnCount > 1}
                {#each titleColumnIndexes as index (`${item.id}-${index}`)}
                  {@const titlePart = itemTitleParts[index] ?? ""}
                  {@const parsedTitlePart = parseTitlePart(titlePart)}
                  <td>
                    {#if titlePart}
                      <span class="target-title">
                        <ProcessingTargetTitleTooltip
                          text={titlePart}
                          displayText={parsedTitlePart.value}
                          tooltipId={`processing-target-title-${item.id}-${index}`}
                        />
                      </span>
                    {/if}
                  </td>
                {/each}
              {:else}
                <td>
                  <span class="target-title">
                    {#each itemTitleParts as titlePart, index (`${item.id}-${index}`)}
                      <ProcessingTargetTitleTooltip
                        text={titlePart}
                        tooltipId={`processing-target-title-${item.id}-${index}`}
                      />
                    {/each}
                  </span>
                </td>
              {/if}
            </tr>
            {#if isExpanded}
              <tr
                class="target-detail-row"
                id={`processing-target-detail-${item.id}`}
              >
                <td colspan={titleColumnCount}>
                  <div class="target-detail-shell">
                    <dl
                      class="target-metadata-list"
                      aria-label={`${itemTitle} のメタデータ`}
                    >
                      {#each item.metadata as metadata (`${metadata.label}-${metadata.value}`)}
                        <div class="target-metadata-item">
                          <dt>{metadata.label}</dt>
                          <dd>{metadata.value}</dd>
                        </div>
                      {/each}
                    </dl>
                    {#if item.actions?.length}
                      <div class="target-detail-extra">
                        <div
                          class="target-detail-actions"
                          aria-label={`${itemTitle} の操作`}
                        >
                          {#each item.actions as action (action.label)}
                            <button
                              class:danger-action={action.variant === "danger"}
                              class="target-detail-action"
                              data-testid={action.testId}
                              disabled={action.disabled}
                              onclick={action.onAction}
                              type="button"
                            >
                              {action.label}
                            </button>
                          {/each}
                        </div>
                      </div>
                    {/if}
                  </div>
                </td>
              </tr>
            {/if}
          {/each}
        {/if}
      </tbody>
    </table>
  </div>

  <PaginationControls
    page={safeCurrentPage}
    {pageCount}
    totalLabel={pageRangeLabel}
    {busy}
    onPrevious={showPreviousPage}
    onNext={showNextPage}
  />
</section>

<style>
  .processing-target-list-panel {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    display: grid;
    gap: 1rem;
    min-width: 0;
    padding: 1.15rem 1.3rem;
  }

  .panel-heading {
    align-items: start;
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    min-width: 0;
  }

  .target-count {
    color: rgba(236, 223, 205, 0.78);
  }

  h3,
  p {
    margin: 0;
  }

  h3 {
    color: #fff6ea;
    overflow-wrap: anywhere;
  }

  .target-count {
    font-weight: 700;
    text-align: right;
  }

  .target-table-frame {
    border: 1px solid rgba(255, 212, 165, 0.14);
    border-radius: 8px;
    min-width: 0;
    overflow-x: clip;
    overflow-y: visible;
  }

  .target-table {
    border-collapse: collapse;
    inline-size: 100%;
    min-width: 100%;
    table-layout: fixed;
  }

  .target-table th,
  .target-table td {
    padding: 0.72rem 0.9rem;
    text-align: left;
  }

  .target-table td {
    min-width: 0;
  }

  .target-table th {
    background: rgba(18, 16, 15, 0.84);
    color: rgba(236, 223, 205, 0.78);
    font-weight: 700;
  }

  .target-table tbody tr {
    background: rgba(18, 16, 15, 0.54);
    border-top: 1px solid rgba(255, 212, 165, 0.12);
    color: #fff6ea;
  }

  .target-table tbody tr:not(.target-detail-row, .target-empty-row) {
    cursor: pointer;
  }

  .target-table tbody tr:not(.target-detail-row, .target-empty-row):hover,
  .target-table tbody tr:not(.target-detail-row, .target-empty-row):focus {
    background: rgba(89, 65, 45, 0.74);
    outline: 2px solid rgba(255, 212, 165, 0.36);
    outline-offset: -2px;
  }

  .target-table tbody tr.expanded {
    background: rgba(122, 80, 44, 0.78);
  }

  .target-empty-row td {
    color: rgba(236, 223, 205, 0.78);
  }

  .target-title {
    align-items: center;
    display: flex;
    inline-size: 100%;
    min-width: 0;
    overflow: hidden;
  }

  .target-title :global(.target-title-part-shell) {
    flex: 1 1 0;
    min-width: 0;
  }

  .target-detail-row td {
    background: rgba(18, 16, 15, 0.72);
    padding: 0;
  }

  .target-detail-shell {
    animation: target-detail-enter 180ms ease-out;
    overflow: hidden;
    transform-origin: top;
  }

  .target-metadata-list {
    animation: target-detail-content-enter 180ms ease-out;
    display: grid;
    gap: 0;
    margin: 0;
    min-width: 0;
    padding: 0.35rem 0.9rem 0.8rem;
  }

  .target-detail-extra {
    animation: target-detail-content-enter 180ms ease-out;
    border-top: 1px solid rgba(255, 212, 165, 0.1);
    display: grid;
    gap: 0.8rem;
    padding: 0 0.9rem 0.9rem;
  }

  .target-detail-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.55rem;
    justify-content: flex-end;
  }

  .target-detail-action {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 212, 165, 0.18);
    border-radius: 999px;
    color: #fff6ea;
    cursor: pointer;
    min-height: 2.5rem;
    min-width: 4.6rem;
    padding: 0 1rem;
  }

  .target-detail-action.danger-action {
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9c7c 100%);
    border-color: transparent;
    color: #35150d;
  }

  .target-detail-action:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .target-metadata-item {
    border-top: 1px solid rgba(255, 212, 165, 0.1);
    display: grid;
    gap: 0.6rem;
    grid-template-columns: minmax(9rem, 0.32fr) minmax(0, 1fr);
    padding: 0.62rem 0;
  }

  .target-metadata-item:first-child {
    border-top: 0;
  }

  .target-metadata-item dt,
  .target-metadata-item dd {
    margin: 0;
    overflow-wrap: anywhere;
  }

  .target-metadata-item dt {
    color: rgba(236, 223, 205, 0.78);
    font-weight: 700;
  }

  .target-metadata-item dd {
    color: #fff6ea;
  }

  @keyframes target-detail-enter {
    from {
      max-height: 0;
    }

    to {
      max-height: 28rem;
    }
  }

  @keyframes target-detail-content-enter {
    from {
      opacity: 0;
      transform: translateY(-0.25rem);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @media (max-width: 720px) {
    .panel-heading {
      align-items: stretch;
      flex-direction: column;
    }

    .target-count {
      text-align: left;
    }

    .target-metadata-item {
      grid-template-columns: 1fr;
    }
  }
</style>
