<script lang="ts">
  // プロンプトテンプレート画面のベースタブ。base 翻訳指示文（全文に付く system 指示）を編集する。
  // state・保存・validation は持たない。値は props、入力 event は callback prop で上へ返す。
  import Field from "@ui/components/Field.svelte"
  import type { PromptTemplateForm, PromptTemplateField } from "./template-editor-view"

  interface Props {
    form: PromptTemplateForm
    onFieldInput: (field: PromptTemplateField, value: string) => void
  }

  let { form, onFieldInput }: Props = $props()

  function handle(field: PromptTemplateField) {
    return (event: Event) =>
      onFieldInput(field, (event.currentTarget as HTMLTextAreaElement).value)
  }
</script>

<div class="flex flex-col gap-3">
  <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
    base 翻訳指示
  </h2>
  <Field
    label="全文に付く system 指示"
    hint="叙述文・台詞のどちらにも付く冒頭の指示文です。訳文だけを返させるなどの基本方針を書きます。"
  >
    <textarea
      class="textarea w-full bg-base-100/60 u-mono text-sm leading-relaxed min-h-28"
      value={form.baseDirective}
      spellcheck="false"
      oninput={handle("baseDirective")}
    ></textarea>
  </Field>
</div>
