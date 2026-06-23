<script lang="ts">
  // 翻訳結果 1 件のコンパクト行。数万件・ページングでも一覧として読めるよう、既定は 1 行に畳む。
  // 畳んだ summary に状態・EDID・口調チップ・固有名チップ・原文/訳文の抜粋を出す。
  // 展開で原文/訳文・置換した固有名・実プロンプト（送信全文）を見せる。口調指示は実プロンプトの system で確かめる。
  // details/summary の開閉は利用者操作で、state は持たない。defaultOpen は story 表示用。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import {
    ATTITUDE_LABEL,
    DECISION_PATH_HINT,
    DECISION_PATH_LABEL,
    EMOTION_LABEL,
    statusTone
  } from "./translation-run-presentation"
  import type { NarrationResultRow } from "./translation-run-view"

  interface Props {
    row: NarrationResultRow
    defaultOpen?: boolean
  }

  let { row, defaultOpen = false }: Props = $props()

  let hasPersona = $derived(!!row.directive)
  let termCount = $derived(row.terms?.length ?? 0)
  let hasPrompt = $derived(!!row.prompt)
  // 話者の属性キャプション（性別 ・ 年齢 ・ 声質）。空の属性は区切りごと省く。
  let speakerAttrs = $derived(
    row.speaker
      ? [
          row.speaker.sex,
          row.speaker.age,
          row.speaker.voice ? `声質 ${row.speaker.voice}` : ""
        ]
          .filter(Boolean)
          .join(" ・ ")
      : ""
  )
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
    {#if row.speaker}
      <span
        class="badge badge-sm badge-outline border-accent/40 text-accent/90 shrink-0 max-w-[10rem] truncate"
        title="この台詞の話者"
      >
        {row.speaker.edid}
      </span>
    {/if}
    {#if hasPersona}
      <span
        class="badge badge-sm badge-outline border-primary/40 text-primary/80 shrink-0 max-w-[12rem] truncate"
      >
        {row.personaLabel ?? "口調あり"}
      </span>
    {:else}
      <span class="u-mono text-[0.65rem] text-base-content/30 shrink-0">口調なし</span>
    {/if}
    {#if termCount > 0}
      <span
        class="badge badge-sm badge-outline border-secondary/40 text-secondary/80 shrink-0 u-mono"
        title="本文で辞書から確定訳語へ置換した固有名の件数"
      >
        固有名 {termCount}
      </span>
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
    {#if row.speaker}
      <div class="mb-4 flex flex-wrap items-baseline gap-x-3 gap-y-1">
        <span class="u-mono text-[0.6rem] uppercase tracking-widest text-accent/60">話者</span>
        <span class="text-sm font-semibold text-base-content/90">{row.speaker.edid}</span>
        {#if speakerAttrs.length > 0}
          <span class="u-mono text-[0.65rem] text-base-content/40">{speakerAttrs}</span>
        {/if}
      </div>
    {/if}
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
    {#if row.persona}
      <div class="mt-4 rounded-box border border-primary/25 bg-primary/5 px-4 py-3">
        <span class="u-mono text-[0.6rem] uppercase tracking-widest text-primary/60">
          口調
        </span>
        <p class="mt-0.5 text-lg font-semibold text-primary">{row.persona.cell}</p>
        <p class="mt-1 text-sm leading-relaxed text-base-content/75">
          {row.persona.trait}
        </p>
        <p
          class="mt-2 u-mono text-[0.65rem] text-base-content/40"
          title={DECISION_PATH_HINT[row.persona.decisionPath]}
        >
          {DECISION_PATH_LABEL[row.persona.decisionPath]} ・ 対人 {ATTITUDE_LABEL[
            row.persona.attitudeBand
          ]} ・ 感情 {EMOTION_LABEL[row.persona.emotionBand]} ・ 印 {row.persona
            .marked}
        </p>
      </div>
    {/if}
    {#if termCount > 0}
      <div class="mt-4 rounded-box border border-secondary/25 bg-secondary/5 p-3">
        <span class="u-mono text-[0.65rem] uppercase tracking-widest text-secondary/70">
          置換した固有名
        </span>
        <ul class="mt-2 flex flex-col gap-1">
          {#each row.terms ?? [] as term (term.source)}
            <li class="text-sm text-base-content/80">
              <span class="u-mono text-xs text-base-content/85">{term.source}</span>
              <span class="text-base-content/30">→</span>
              <span class="text-base-content">{term.dest}</span>
            </li>
          {/each}
        </ul>
      </div>
    {/if}
    {#if hasPrompt}
      <div class="mt-4 rounded-box border border-info/25 bg-info/5 p-3">
        <span class="u-mono text-[0.65rem] uppercase tracking-widest text-info/70">
          実プロンプト
        </span>
        <pre class="mt-2 max-h-64 overflow-auto whitespace-pre-wrap break-words rounded-box bg-base-300/30 p-2.5 u-mono text-xs leading-relaxed text-base-content/75">{row.prompt}</pre>
      </div>
    {/if}
  </div>
</details>
