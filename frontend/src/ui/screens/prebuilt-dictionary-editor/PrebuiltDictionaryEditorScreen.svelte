<script lang="ts">
  import type {
    PrebuiltDictionaryFilters,
    PrebuiltDictionaryRow
  } from "./prebuilt-dictionary-editor-view"
  import DataTable from "../../components/DataTable.svelte"

  interface Props {
    rows: PrebuiltDictionaryRow[]
    filters: PrebuiltDictionaryFilters
    pageNumber: number
    totalCount: number
    canPrev: boolean
    canNext: boolean
    hasPendingChanges?: boolean
    expandedRowIds?: string[]
    onFilterInput: (field: keyof PrebuiltDictionaryFilters, value: string) => void
    onCreate: () => void
    onDelete: (id: string) => void
    onToggleCategories: (id: string) => void
    onConfirmChanges: () => void
    onCancelChanges: () => void
    onPrev: () => void
    onNext: () => void
  }

  let {
    rows,
    filters,
    pageNumber,
    totalCount,
    canPrev,
    canNext,
    hasPendingChanges = false,
    expandedRowIds = [],
    onFilterInput,
    onCreate,
    onDelete,
    onToggleCategories,
    onConfirmChanges,
    onCancelChanges,
    onPrev,
    onNext
  }: Props = $props()

  const columns: Array<{ key: keyof PrebuiltDictionaryFilters; label: string }> = [
    { key: "source", label: "原語" },
    { key: "destination", label: "訳語" },
    { key: "partOfSpeech", label: "品詞" },
    { key: "category", label: "Skyrimカテゴリ" }
  ]

  function filterInput(field: keyof PrebuiltDictionaryFilters, event: Event): void {
    onFilterInput(field, (event.currentTarget as HTMLInputElement).value)
  }
</script>

<div class="min-h-screen w-full px-6 py-12 flex justify-center">
  <section class="w-full max-w-7xl flex flex-col gap-8" aria-labelledby="prebuilt-dictionary-title">
    <header class="flex flex-col gap-3">
      <span class="u-mono text-xs tracking-[0.32em] uppercase text-accent">Term dictionary</span>
      <div class="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 id="prebuilt-dictionary-title" class="u-display text-4xl font-semibold text-base-content">用語辞書</h1>
          <p class="mt-2 text-base-content/70">辞書を検索し，内容を追加，編集，削除します。</p>
        </div>
        <button class="btn btn-primary" type="button" onclick={onCreate}>新規作成</button>
      </div>
      <div class="h-px w-full bg-gradient-to-r from-transparent via-primary/50 to-transparent"></div>
    </header>

    <section class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top overflow-hidden" aria-label="用語辞書の一覧">
      <div class="flex items-center justify-between gap-4 px-6 pt-5">
        <p class="text-sm text-base-content/65"><span class="u-mono text-base-content">{totalCount.toLocaleString()}</span> 件</p>
        <p class="text-sm text-base-content/55">1ページ 50件</p>
      </div>
      <DataTable ariaLabel="用語辞書の一覧">
          <thead>
            <tr>
              {#each columns as column (column.key)}
                <th class="min-w-48">
                  <label class="flex flex-col gap-2">
                    <span>{column.label}</span>
                    <input
                      class="input input-sm input-bordered w-full bg-base-100 font-normal"
                      value={filters[column.key]}
                      placeholder={`${column.label}で絞り込む`}
                      oninput={(event) => filterInput(column.key, event)}
                    />
                  </label>
                </th>
              {/each}
              <th class="min-w-44"><span class="sr-only">操作</span></th>
            </tr>
          </thead>
          <tbody>
            {#if rows.length === 0}
              <tr><td class="py-12 text-center text-base-content/60" colspan="5">一致する辞書はありません。</td></tr>
            {:else}
              {#each rows as row (row.id)}
                <tr
                  class={row.pending === "edited"
                    ? "bg-warning/15"
                    : row.pending === "deleted"
                      ? "bg-base-300/50 opacity-50"
                      : ""}
                >
                  <td><input class="input input-sm input-ghost u-mono min-w-64 max-w-96 w-full" value={row.source} title={row.source} disabled={row.pending === "deleted"} /></td>
                  <td><input class="input input-sm input-ghost min-w-64 max-w-96 w-full" value={row.destination} title={row.destination} disabled={row.pending === "deleted"} /></td>
                  <td><input class="input input-sm input-ghost w-full" value={row.partOfSpeech} disabled={row.pending === "deleted"} /></td>
                  <td>
                    {#if row.categories.length > 1}
                      <button class="btn btn-sm btn-ghost" type="button" onclick={() => onToggleCategories(row.id)}>{expandedRowIds.includes(row.id) ? "カテゴリを閉じる" : `カテゴリ ${row.categories.length}件`}</button>
                    {:else}
                      <span class="u-mono text-sm">{row.categories[0]}</span>
                    {/if}
                  </td>
                  <td>
                    <div class="flex flex-nowrap items-center justify-end gap-2 whitespace-nowrap">
                      <button class="btn btn-sm btn-ghost text-error" type="button" onclick={() => onDelete(row.id)}>削除</button>
                    </div>
                  </td>
                </tr>
                {#if expandedRowIds.includes(row.id)}
                  {#each row.categories as category (category)}
                    <tr class="bg-base-200/60">
                      <td colspan="3" class="text-right text-sm text-base-content/55">Skyrimカテゴリ</td>
                      <td><span class="u-mono text-sm">{category}</span></td>
                      <td></td>
                    </tr>
                  {/each}
                {/if}
              {/each}
            {/if}
          </tbody>
      </DataTable>
      {#if hasPendingChanges}
        <div class="flex items-center justify-end gap-3 border-t border-base-300/70 p-5">
          <span class="mr-auto text-sm text-warning">変更を確定または取消するまでページを移動できません。</span>
          <button class="btn btn-ghost" type="button" onclick={onCancelChanges}>取消</button>
          <button class="btn btn-primary" type="button" onclick={onConfirmChanges}>変更を確定</button>
        </div>
      {:else}
      <nav class="flex items-center justify-center gap-3 p-5" aria-label="用語辞書のページ送り">
        <button class="btn btn-sm btn-ghost" type="button" disabled={!canPrev} onclick={onPrev}>← 前へ</button>
        <span class="u-mono text-xs text-base-content/60 min-w-[5rem] text-center">ページ {pageNumber}</span>
        <button class="btn btn-sm btn-ghost" type="button" disabled={!canNext} onclick={onNext}>次へ →</button>
      </nav>
      {/if}
    </section>
  </section>
</div>
