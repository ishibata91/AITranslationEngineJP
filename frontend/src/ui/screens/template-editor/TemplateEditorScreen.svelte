<script lang="ts">
  // プロンプトテンプレート編集画面。base 翻訳指示文と口調指示テンプレートを編集する。
  // 翻訳実行画面と同じく、state・保存・validation は持たない。値は props、入力 event は callback prop で上へ返す。
  import Field from "@ui/components/Field.svelte"
  import type {
    PromptTemplateForm,
    PromptTemplateField,
    TemplatePlaceholder
  } from "./template-editor-view"

  interface Props {
    form: PromptTemplateForm
    // 口調指示テンプレートに差し込めるプレースホルダ一覧。
    placeholders: TemplatePlaceholder[]
    // 未保存の変更があるか。保存・戻すボタンの活性に使う。
    dirty?: boolean
    // 保存中か。保存ボタンに spinner を出す。
    saving?: boolean
    onFieldInput: (field: PromptTemplateField, value: string) => void
    onSave: () => void
    onReset?: () => void
  }

  let {
    form,
    placeholders,
    dirty = false,
    saving = false,
    onFieldInput,
    onSave,
    onReset = () => {}
  }: Props = $props()

  function handle(field: PromptTemplateField) {
    return (event: Event) =>
      onFieldInput(field, (event.currentTarget as HTMLTextAreaElement).value)
  }
</script>

<div class="min-h-screen w-full px-6 py-12 flex justify-center">
  <section class="w-full max-w-4xl flex flex-col gap-8" aria-labelledby="template-title">
    <header class="flex flex-col gap-3">
      <span class="u-mono text-xs tracking-[0.32em] uppercase text-accent">
        Prompt template
      </span>
      <h1 id="template-title" class="u-display text-4xl font-semibold text-base-content">
        プロンプトテンプレート
      </h1>
      <p class="max-w-2xl text-base-content/70 leading-relaxed">
        翻訳 AI へ送る指示文の雛形を編集します。base 指示は叙述文・台詞の両方に付き、口調指示は話者のいる台詞だけに付きます。
      </p>
      <div
        class="h-px w-full bg-gradient-to-r from-transparent via-primary/50 to-transparent"
      ></div>
    </header>

    <div class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top">
      <div class="card-body gap-6">
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

        <div class="flex flex-col gap-3">
          <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
            口調指示テンプレート
          </h2>
          <Field
            label="話者のいる台詞に付く口調指示"
            hint="話者の性質を差し込み口へ展開して、口調を保つ指示文を作ります。"
          >
            <textarea
              class="textarea w-full bg-base-100/60 u-mono text-sm leading-relaxed min-h-32"
              value={form.personaTemplate}
              spellcheck="false"
              oninput={handle("personaTemplate")}
            ></textarea>
          </Field>
          <div class="rounded-box border border-base-300/50 bg-base-100/30 p-3">
            <span class="u-mono text-[0.65rem] uppercase tracking-widest text-base-content/40">
              差し込めるプレースホルダ
            </span>
            <ul class="mt-2 flex flex-col gap-1.5">
              {#each placeholders as ph (ph.token)}
                <li class="text-sm text-base-content/75">
                  <span class="u-mono text-xs text-secondary/80">{ph.token}</span>
                  <span class="text-base-content/30"> — </span>
                  <span>{ph.description}</span>
                </li>
              {/each}
            </ul>
          </div>
        </div>

        <div class="card-actions items-center justify-end gap-4 pt-1">
          {#if dirty}
            <span class="u-mono text-xs text-warning/80">未保存の変更</span>
          {/if}
          <button
            class="btn btn-ghost"
            type="button"
            disabled={!dirty || saving}
            onclick={onReset}
          >
            戻す
          </button>
          <button
            class="btn btn-primary px-8"
            type="button"
            disabled={!dirty || saving}
            onclick={onSave}
          >
            {#if saving}
              <span class="loading loading-spinner loading-xs"></span>
            {/if}
            保存
          </button>
        </div>
      </div>
    </div>
  </section>
</div>
