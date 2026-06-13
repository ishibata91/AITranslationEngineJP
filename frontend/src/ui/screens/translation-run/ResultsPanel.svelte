<script lang="ts">
  // 結果一覧のパネル。件数バッジ、空状態（未実行／実行中）、結果行の並びを出す。
  import TranslationResultRow from "./TranslationResultRow.svelte"
  import type { NarrationResultRow, RunPhase } from "./translation-run-view"

  interface Props {
    phase: RunPhase
    results: NarrationResultRow[]
  }

  let { phase, results }: Props = $props()
</script>

<div class="card bg-base-200/40 border border-base-300/60">
  <div class="card-body gap-5">
    <div class="flex items-center justify-between">
      <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
        結果一覧
      </h2>
      {#if results.length > 0}
        <span class="badge badge-outline badge-sm u-mono">{results.length} 件</span>
      {/if}
    </div>

    {#if results.length === 0}
      <div class="flex flex-col items-center gap-3 py-12 text-center text-base-content/50">
        {#if phase === "running"}
          <span class="loading loading-ring loading-lg text-primary"></span>
          <p>抽出と翻訳を実行しています。完了すると原文と訳文が並びます。</p>
        {:else}
          <span class="u-display text-3xl text-base-content/25">ᚱ</span>
          <p>まだ結果はありません。plugin を選び、AI サービスを入力して実行してください。</p>
        {/if}
      </div>
    {:else}
      <ul class="flex flex-col gap-4">
        {#each results as row (row.edid + row.source)}
          <li>
            <TranslationResultRow {row} />
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>
