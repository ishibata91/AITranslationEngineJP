<script lang="ts">
  // プロンプトテンプレート画面のレコード別タブの一節。話者なし台詞（汎用・PC 発話）の口調を
  // 自由記述で書く。対人の文面は利用者が書き、感情の強弱と性別の一人称・語尾はシステムが自動で重ねる。
  // PC の性別は一人称・語尾の根拠として選ぶ。state・保存・validation は持たない。
  // 値は props、入力 event は onFieldInput callback で上へ返す（既存の入力中継口を再利用）。
  import Field from "@ui/components/Field.svelte"
  import { PC_SEX_OPTIONS, TONE_TEXT_FIELDS } from "./template-editor-presentation"
  import type { PromptTemplateField, PromptTemplateForm } from "./template-editor-view"

  interface Props {
    form: PromptTemplateForm
    onFieldInput: (field: PromptTemplateField, value: string) => void
  }

  let { form, onFieldInput }: Props = $props()

  function handleText(field: PromptTemplateField) {
    return (event: Event) =>
      onFieldInput(field, (event.currentTarget as HTMLTextAreaElement).value)
  }
</script>

<div class="flex flex-col gap-4">
  <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
    話者なし台詞の口調
  </h2>
  <p class="text-sm text-base-content/55">
    話者を特定できない台詞（汎用台詞・PC 発話）の口調を自由記述で書きます。対人の文面はここで書き、感情の強弱（台詞の本文から）と性別の一人称・語尾はシステムが自動で重ねます。
  </p>
  {#each TONE_TEXT_FIELDS as f (f.field)}
    <Field label={f.label} hint={f.hint}>
      <textarea
        class="textarea w-full bg-base-100/60 u-mono text-sm leading-relaxed min-h-24"
        value={form[f.field] ?? ""}
        spellcheck="false"
        oninput={handleText(f.field)}
      ></textarea>
    </Field>
  {/each}
  <Field
    label="PC の性別"
    hint="PC 発話の一人称・語尾の根拠です。プレイヤーキャラクターの性別を選びます。未指定なら一人称・語尾を付けません。"
  >
    <div role="radiogroup" aria-label="PC の性別" class="join">
      {#each PC_SEX_OPTIONS as opt (opt.value)}
        <button
          type="button"
          role="radio"
          aria-checked={(form.pcSex ?? "") === opt.value}
          class="btn btn-sm join-item"
          class:btn-primary={(form.pcSex ?? "") === opt.value}
          onclick={() => onFieldInput("pcSex", opt.value)}
        >
          {opt.label}
        </button>
      {/each}
    </div>
  </Field>
</div>
