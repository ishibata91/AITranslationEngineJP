<script lang="ts">
  import type { MasterPersonaListItem } from "@application/gateway-contract/master-persona/master-persona-gateway-contract"
  import ProcessingTargetListWrapper from "@ui/components/ProcessingTargetListWrapper.svelte"
  import type { ProcessingTargetListItem } from "@ui/components/processing-target-list-panel-types"
  import type { PersonaReviewPanelProps } from "./master-persona-panel-props"

  let {
    canMutate,
    items,
    keyword,
    page,
    pageSize,
    pluginFilter,
    pluginOptions,
    selectedEntry,
    selectedIdentityKey,
    totalCount,
    totalPages,
    selectRow,
    updateKeyword,
    updatePluginFilter,
    goToPrevPage,
    goToNextPage,
    editCurrent,
    openDelete
  }: PersonaReviewPanelProps = $props()

  function itemLabel(item: MasterPersonaListItem): string {
    return item.displayName || item.editorId || item.formId
  }

  function buildPersonaTargetItems(
    listItems: MasterPersonaListItem[]
  ): ProcessingTargetListItem[] {
    return listItems.map((item) => {
      const selectedItem =
        selectedEntry?.identityKey === item.identityKey ? selectedEntry : null
      return {
        id: item.identityKey,
        name: itemLabel(item),
        metadata: [
          { label: "FormID", value: item.formId },
          { label: "EditorID", value: item.editorId },
          { label: "対象プラグイン", value: item.targetPlugin },
          { label: "元プラグイン", value: item.sourcePlugin },
          { label: "声", value: item.voiceType },
          { label: "話し方", value: selectedItem?.speechStyle || "未入力" },
          {
            label: "ペルソナ本文",
            value:
              selectedItem?.personaBody ||
              item.personaSummary ||
              "生成後に一覧から選ぶと、本文を確認できます。"
          },
          { label: "最終更新", value: item.updatedAt }
        ],
        detail: "",
        actions: [
          {
            label: "編集",
            disabled: !canMutate,
            onAction: editCurrent
          },
          {
            label: "削除",
            variant: "danger" as const,
            disabled: !canMutate,
            onAction: openDelete
          }
        ]
      }
    })
  }

  const pageRangeText = $derived(
    totalCount === 0
      ? "0 件"
      : `${(page - 1) * pageSize + 1}-${Math.min(page * pageSize, totalCount)} / ${totalCount.toLocaleString("ja-JP")} 件`
  )
  const personaTargetItems = $derived(buildPersonaTargetItems(items))
  const visiblePluginNames = $derived(
    Array.from(new Set(items.map((item) => item.targetPlugin))).join(" / ")
  )
</script>

<section
  class="review-grid"
  aria-label="生成結果の確認"
  data-testid="master-persona-generation-result-list-panel"
>
  <ProcessingTargetListWrapper
    countText={pageRangeText}
    eyebrow="生成結果"
    filterId="masterPersonaPluginSelect"
    filterLabel="プラグイン"
    filterOptions={pluginOptions}
    filterValue={pluginFilter}
    initialExpandedItemId={selectedIdentityKey}
    items={personaTargetItems}
    pageSize={pageSize}
    searchId="masterPersonaSearchInput"
    searchLabel="検索"
    searchPlaceholder="名前またはプラグイン名で検索"
    searchValue={keyword}
    supportText={visiblePluginNames}
    title="ペルソナ一覧"
    titleId="listHeading"
    onFilterChange={updatePluginFilter}
    onSearchInput={updateKeyword}
    onSelectItem={selectRow}
  >
    {#snippet footer()}
      <nav class="pager-row" aria-label="ペルソナ一覧のページ操作">
        <span class="support-copy">{page} / {totalPages} ページ</span>
        <div class="pager-actions">
          <button
            class="button-secondary"
            disabled={page <= 1}
            id="prevPageButton"
            onclick={goToPrevPage}
            type="button"
          >
            前へ
          </button>
          <button
            class="button-secondary"
            disabled={page >= totalPages}
            id="nextPageButton"
            onclick={goToNextPage}
            type="button"
          >
            次へ
          </button>
        </div>
      </nav>
    {/snippet}
  </ProcessingTargetListWrapper>
</section>

<style>
  .review-grid {
    align-items: start;
    color: var(--text);
    display: grid;
    gap: 12px;
    grid-template-columns: minmax(0, 1fr);
  }

  .pager-row,
  .pager-actions {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    justify-content: space-between;
  }

  .support-copy {
    overflow-wrap: anywhere;
  }

  .support-copy {
    color: var(--muted);
    line-height: 1.7;
  }

  .button-secondary {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    color: var(--text);
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  @media (max-width: 560px) {
    .pager-row,
    .pager-actions {
      width: 100%;
    }

    .pager-actions > * {
      width: 100%;
    }
  }
</style>
