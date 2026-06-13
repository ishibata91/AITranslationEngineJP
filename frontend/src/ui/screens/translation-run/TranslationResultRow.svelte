<script lang="ts">
  // 翻訳結果 1 件の見開き表示。左に原文、右に訳文を並べ、状態をバッジで示す。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import { statusTone } from "./translation-run-presentation"
  import type { NarrationResultRow } from "./translation-run-view"

  interface Props {
    row: NarrationResultRow
  }

  let { row }: Props = $props()
</script>

<article class="rounded-box border border-base-300/60 bg-base-100/40 p-5 u-edge-top">
  <div class="mb-3 flex items-center justify-between gap-3">
    <span class="u-mono text-xs text-base-content/55">{row.edid}</span>
    <StatusBadge label={row.statusLabel} tone={statusTone(row.statusLabel)} />
  </div>
  <div class="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6">
    <div class="flex flex-col gap-1.5">
      <span class="u-mono text-[0.65rem] uppercase tracking-widest text-base-content/40">
        原文
      </span>
      <p class="leading-relaxed text-base-content/85">{row.source}</p>
    </div>
    <div class="flex flex-col gap-1.5 md:border-l md:border-base-300/50 md:pl-6">
      <span class="u-mono text-[0.65rem] uppercase tracking-widest text-accent/70">
        訳文
      </span>
      {#if row.dest.length > 0}
        <p class="leading-relaxed text-base-content">{row.dest}</p>
      {:else}
        <p class="italic leading-relaxed text-base-content/35">（訳文なし）</p>
      {/if}
    </div>
  </div>
</article>
