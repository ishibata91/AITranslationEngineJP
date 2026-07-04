<script lang="ts">
  // 結果一覧のパネル。総件数バッジ、空状態（未実行／実行中）、現在ページの結果行、ページャを出す。
  // results はページ分のみ。総件数・現在ページ番号・前後可否はページング props で受ける。
  // ページング props は state を持たず、操作は callback で上へ返す。既定は単一ページ（前後無効）。
  import TranslationResultRow from "./TranslationResultRow.svelte"
  import ResultsPager from "./ResultsPager.svelte"
  import type { NarrationResultRow, RunPhase } from "./translation-run-view"

  interface Props {
    phase: RunPhase
    // 現在ページの結果行。ページングするため全件ではない。
    results: NarrationResultRow[]
    // 結果全体の総件数。未指定なら現在ページの件数を総件数とみなす（単一ページ）。
    total?: number
    // 現在ページ番号（1 始まり）。
    pageNumber?: number
    canPrev?: boolean
    canNext?: boolean
    onPrev?: () => void
    onNext?: () => void
    // xTranslator XML への書き出し操作。結果があるときだけ押せる。実処理は callback で上へ返す。
    onExport?: () => void
    // 書き出し中フラグ。true の間はボタンを無効化しスピナーを出す。
    exporting?: boolean
  }

  let {
    phase,
    results,
    total = results.length,
    pageNumber = 1,
    canPrev = false,
    canNext = false,
    onPrev = () => {},
    onNext = () => {},
    onExport = () => {},
    exporting = false
  }: Props = $props()
</script>

<div class="card bg-base-200/40 border border-base-300/60">
  <div class="card-body gap-5">
    <div class="flex items-center justify-between gap-3">
      <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
        結果一覧
      </h2>
      <div class="flex items-center gap-3">
        {#if total > 0}
          <span class="badge badge-outline badge-sm u-mono">{total} 件</span>
        {/if}
        {#if results.length > 0}
          <button
            type="button"
            class="btn btn-sm btn-outline btn-primary"
            disabled={exporting}
            onclick={onExport}
          >
            {#if exporting}
              <span class="loading loading-spinner loading-xs"></span>
            {/if}
            {exporting ? "書き出し中…" : "xTranslator へ書き出し"}
          </button>
        {/if}
      </div>
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
      <ul class="flex flex-col gap-2">
        <!-- 同一原文の台詞が重複しうるため、key は一意な並び順 index にする。一覧は毎回まるごと差し替えで順序は安定。 -->
        {#each results as row, i (i)}
          <li>
            <TranslationResultRow {row} />
          </li>
        {/each}
      </ul>
      <ResultsPager {pageNumber} {canPrev} {canNext} {onPrev} {onNext} />
    {/if}
  </div>
</div>
