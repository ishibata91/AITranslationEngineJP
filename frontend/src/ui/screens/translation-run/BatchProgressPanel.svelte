<script lang="ts">
  // xAI batch の進行状況表示。固有名 → 本文 → 完了 のステッパーで現在地を見せ、現段 batch の件数を出す。
  // 全体で 2 段（固有名・本文）あることが常に見えるため、取り込みが固有名取り込み＋本文送信の二重動作である
  // ことを操作前に読み取れる。表示専用。値は props で受け取り、event 購読や state・API はここに持たない。
  import {
    BATCH_COUNT_LABEL,
    BATCH_UNCHECKED_HINT,
    BATCH_WAITING_HINT,
    BATCH_APPLYABLE_HINT,
    BATCH_DONE_HINT,
    batchStepViews
  } from "./translation-run-presentation"
  import type { BatchProgressView } from "./translation-run-view"

  // 進行状況。状態確認前（未確認）は undefined。
  let { progress }: { progress?: BatchProgressView } = $props()

  // ステッパー各段の表示状態（過去段は ✓、現在段は primary、未確認は中立）。
  const steps = $derived(batchStepViews(progress))

  // パネル下部の補足。未確認・全完了・完了段あり・処理待ちで出し分ける。
  const hint = $derived.by(() => {
    if (!progress) return BATCH_UNCHECKED_HINT
    if (progress.stage === "done") return BATCH_DONE_HINT
    if (progress.canApply) return BATCH_APPLYABLE_HINT
    return BATCH_WAITING_HINT
  })
</script>

<div class="card bg-base-200/40 border border-base-300/60 u-edge-top">
  <div class="card-body gap-4 py-5">
    <div class="flex items-center justify-between gap-3">
      <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
        進行状況
      </h2>
      {#if !progress}
        <span class="u-mono text-xs text-base-content/45">未確認</span>
      {/if}
    </div>

    <ul class="steps steps-horizontal w-full text-xs">
      {#each steps as step (step.stage)}
        <li class="step {step.cls}" data-content={step.content ?? undefined}>{step.label}</li>
      {/each}
    </ul>

    {#if progress}
      <div class="flex flex-wrap gap-x-6 gap-y-1 u-mono text-xs text-base-content/70">
        <span>{BATCH_COUNT_LABEL.total} {progress.total}</span>
        <span class:text-warning={progress.pending > 0}>
          {BATCH_COUNT_LABEL.pending} {progress.pending}
        </span>
        <span>{BATCH_COUNT_LABEL.succeeded} {progress.succeeded}</span>
        <span class:text-error={progress.failed > 0}>
          {BATCH_COUNT_LABEL.failed} {progress.failed}
        </span>
      </div>
    {/if}

    <p class="text-xs text-base-content/50">{hint}</p>
  </div>
</div>
