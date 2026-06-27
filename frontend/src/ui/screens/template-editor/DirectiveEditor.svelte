<script lang="ts">
  // 指示文（directive）の編集表示。指示文ごとに本文 textarea・変数・対象 REC:FIELD（読み取り専用）を並べる。
  // 口調・文体・固有名・定型句を同じ形で出す。値の保持はしない。入力は onInstructionInput callback で上へ返す。
  import Field from "@ui/components/Field.svelte"
  import type { DirectiveSection } from "./directive-view"

  interface Props {
    // 指示文（本文・変数・対象）の一覧。
    sections: DirectiveSection[]
    onInstructionInput: (key: string, instruction: string) => void
  }

  let { sections, onInstructionInput }: Props = $props()

  function handle(key: string) {
    return (event: Event) =>
      onInstructionInput(key, (event.currentTarget as HTMLTextAreaElement).value)
  }
</script>

<div class="flex flex-col gap-7">
  {#each sections as s (s.key)}
    <div class="flex flex-col gap-2">
      <Field label={s.key} hint="Base 指示に続けて、この指示文を合成して訳します。">
        <textarea
          class="textarea w-full bg-base-100/60 u-mono text-sm leading-relaxed min-h-24"
          value={s.instruction}
          spellcheck="false"
          oninput={handle(s.key)}
        ></textarea>
      </Field>
      {#if s.variables.length > 0}
        <div class="rounded-box border border-base-300/50 bg-base-100/30 p-3">
          <span class="u-mono text-[0.65rem] uppercase tracking-widest text-base-content/40">
            差し込める変数
          </span>
          <ul class="mt-1.5 flex flex-col gap-1">
            {#each s.variables as v (v.token)}
              <li class="text-sm text-base-content/75">
                <span class="u-mono text-xs text-secondary/80">{v.token}</span>
                <span class="text-base-content/30"> — </span>
                <span>{v.description}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
      <div class="rounded-box border border-base-300/50 bg-base-100/30 px-3 py-2">
        <span class="u-mono text-[0.65rem] uppercase tracking-widest text-base-content/40">
          対象（{s.targets.length}）
        </span>
        <ul class="mt-1.5 flex flex-wrap gap-x-4 gap-y-1">
          {#each s.targets as t (t.recField)}
            <li class="text-sm text-base-content/75">
              <span class="u-mono text-xs text-base-content/85">{t.recField}</span>
              <span class="text-base-content/30"> </span>
              <span>{t.logicalName}</span>
            </li>
          {/each}
        </ul>
      </div>
    </div>
  {/each}
</div>
