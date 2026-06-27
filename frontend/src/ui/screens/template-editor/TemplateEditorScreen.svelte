<script lang="ts">
  // プロンプトテンプレート編集画面。プロンプト = Base 指示 + その種別の指示文（変数を実行時に埋めたもの）。
  // ベースタブで Base 指示、レコード別タブで種別ごとの指示文（口調・文体・固有名・定型句を同じ形）を編集する。
  // state・保存・validation は持たない。値は props、入力 event は callback prop で上へ返す。
  // レコード別の props は任意。配線前（Container 未更新）でも壊れないようにする。
  import TemplateBasePane from "./TemplateBasePane.svelte"
  import DirectiveEditor from "./DirectiveEditor.svelte"
  import { buildDirectiveSections } from "./directive-presentation"
  import type {
    PromptTemplateForm,
    PromptTemplateField,
    TemplateTab
  } from "./template-editor-view"
  import type { Directive, RecordAssignment } from "./directive-view"

  interface Props {
    form: PromptTemplateForm
    // 表示中のサブタブ。
    activeTab?: TemplateTab
    // 種別ごとの指示文（口調・文体・固有名・定型句）。レコード別タブで使う。
    directives?: Directive[]
    // REC:FIELD と指示文キーの対応。指示文ごとの対象一覧を組む。
    assignments?: RecordAssignment[]
    // 未保存の変更があるか。保存・戻すボタンの活性に使う。
    dirty?: boolean
    // 保存中か。保存ボタンに spinner を出す。
    saving?: boolean
    onFieldInput: (field: PromptTemplateField, value: string) => void
    onInstructionInput?: (key: string, instruction: string) => void
    onTabChange?: (tab: TemplateTab) => void
    onSave: () => void
    onReset?: () => void
  }

  let {
    form,
    activeTab = "base",
    directives = [],
    assignments = [],
    dirty = false,
    saving = false,
    onFieldInput,
    onInstructionInput = () => {},
    onTabChange = () => {},
    onSave,
    onReset = () => {}
  }: Props = $props()

  let sections = $derived(buildDirectiveSections(directives, assignments))
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
        翻訳 AI へ送るプロンプトを編集します。プロンプトは Base 指示と、その種別に割り当てた指示文（変数を実行時に埋めたもの）で組み立てます。
      </p>
      <div
        class="h-px w-full bg-gradient-to-r from-transparent via-primary/50 to-transparent"
      ></div>
    </header>

    <div role="tablist" class="tabs tabs-border">
      <button
        role="tab"
        type="button"
        class="tab"
        class:tab-active={activeTab === "base"}
        onclick={() => onTabChange("base")}
      >
        ベース
      </button>
      <button
        role="tab"
        type="button"
        class="tab"
        class:tab-active={activeTab === "record"}
        onclick={() => onTabChange("record")}
      >
        レコード別
      </button>
    </div>

    {#if activeTab === "base"}
      <div class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top">
        <div class="card-body gap-6">
          <TemplateBasePane {form} {onFieldInput} />
        </div>
      </div>
    {:else}
      <div class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top">
        <div class="card-body gap-4">
          <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
            種別ごとの指示文
          </h2>
          <p class="text-sm text-base-content/55">
            REC:FIELD ごとに割り当てた指示文を編集します。各指示文の「対象」は、その指示文で訳す種別です（読み取り専用）。複数の種別が同じ指示文を共有するため、直すと対象すべてに効きます。
          </p>
          <DirectiveEditor {sections} {onInstructionInput} />
        </div>
      </div>
    {/if}

    <div class="flex items-center justify-end gap-4">
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
  </section>
</div>
