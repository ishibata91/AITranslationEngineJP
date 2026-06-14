<script lang="ts">
  // 翻訳結果 1 件のコンパクト行。数万件・ページングでも一覧として読めるよう、既定は 1 行に畳む。
  // 畳んだ summary に状態・EDID・口調チップ・原文/訳文の抜粋を出し、展開で全文と口調指示を見せる。
  // details/summary の開閉は利用者操作で、state は持たない。defaultOpen は story 表示用。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import { statusTone } from "./translation-run-presentation"
  import type { NarrationResultRow } from "./translation-run-view"

  interface Props {
    row: NarrationResultRow
    defaultOpen?: boolean
  }

  let { row, defaultOpen = false }: Props = $props()

  let hasPersona = $derived(!!row.directive && row.directive.length > 0)
</script>

<details
  class="group rounded-box border border-base-300/50 bg-base-100/30 open:bg-base-100/50"
  open={defaultOpen}
>
  <summary
    class="flex cursor-pointer list-none items-center gap-3 px-4 py-2.5 [&::-webkit-details-marker]:hidden"
  >
    <StatusBadge label={row.statusLabel} tone={statusTone(row.statusLabel)} />
    <span class="u-mono text-xs text-base-content/55 shrink-0 max-w-[10rem] truncate">
      {row.edid}
    </span>
    {#if hasPersona}
      <span
        class="badge badge-sm badge-outline border-primary/40 text-primary/80 shrink-0 max-w-[12rem] truncate"
      >
        {row.personaLabel ?? "口調あり"}
      </span>
    {:else}
      <span class="u-mono text-[0.65rem] text-base-content/30 shrink-0">口調なし</span>
    {/if}
    <span class="min-w-0 flex-1 truncate text-sm text-base-content/65">
      {row.source}
      <span class="text-base-content/30">→</span>
      {#if row.dest.length > 0}
        <span class="text-base-content/80">{row.dest}</span>
      {:else}
        <span class="italic text-base-content/30">（訳文なし）</span>
      {/if}
    </span>
    <span
      class="text-base-content/30 shrink-0 transition-transform group-open:rotate-90"
      aria-hidden="true">▸</span
    >
  </summary>

  <div class="border-t border-base-300/40 px-4 py-4">
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
    {#if hasPersona}
      <div class="mt-4 rounded-box border border-primary/25 bg-primary/5 p-3">
        <span class="u-mono text-[0.65rem] uppercase tracking-widest text-primary/70">
          口調指示
        </span>
        <p class="mt-1 whitespace-pre-line text-xs leading-relaxed text-base-content/70">
          {row.directive}
        </p>
      </div>
    {/if}
  </div>
</details>
