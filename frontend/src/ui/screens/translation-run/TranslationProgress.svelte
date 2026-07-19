<script lang="ts">
  // 実行中の進捗表示。台詞抽出は不定バー、本文翻訳は done/total の確定バーを出す。
  // 表示専用。進捗値は props で受け取り、event 購読や state はここに持たない。
  import {
    EXTRACT_STEP_LABEL,
    PROGRESS_STAGE_LABEL
  } from "./translation-run-presentation"
  import type { RunProgress } from "./translation-run-view"

  let { stage, done, total, step }: RunProgress = $props()

  // 本文翻訳の進捗率（0〜100）。総件数 0 のときは 0 にして 0 除算を避ける。
  let percent = $derived(
    stage === "translate" && total > 0 ? Math.round((done / total) * 100) : 0
  )

  // 見出しラベル。extract 段でサブ段が指定されていればサブ段ラベル、
  // それ以外（未指定・translate 段）は従来の段ラベルにフォールバックする。
  let heading = $derived(
    stage === "extract" && step
      ? EXTRACT_STEP_LABEL[step]
      : PROGRESS_STAGE_LABEL[stage]
  )
</script>

<div class="card bg-base-200/40 border border-base-300/60 u-edge-top">
  <div class="card-body gap-3 py-5">
    <div class="flex items-center justify-between gap-3">
      <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
        {heading}
      </h2>
      {#if stage === "translate"}
        <span class="u-mono text-xs text-base-content/55">
          {done} / {total}（{percent}%）
        </span>
      {/if}
    </div>

    {#if stage === "translate"}
      <progress
        class="progress progress-primary w-full"
        value={done}
        max={total}
        aria-label="本文翻訳の進捗"
      ></progress>
    {:else}
      <progress
        class="progress progress-primary w-full"
        aria-label="台詞抽出の進捗（処理中）"
      ></progress>
    {/if}
  </div>
</div>
