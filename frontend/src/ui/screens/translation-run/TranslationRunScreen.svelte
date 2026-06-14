<script lang="ts">
  // 翻訳実行画面。汎用の表示部品を組み立てるだけにし、state・API・Wails・validation は持たない。
  // すべて props で受け、入力 event は callback prop で上へ返す。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import TextField from "@ui/components/TextField.svelte"
  import SelectField from "@ui/components/SelectField.svelte"
  import FileSelectField from "@ui/components/FileSelectField.svelte"
  import ResultsPanel from "./ResultsPanel.svelte"
  import TranslationProgress from "./TranslationProgress.svelte"
  import { PROVIDER_FIELDS, PHASE_PRESENTATION } from "./translation-run-presentation"
  import type {
    TranslationRunForm,
    TranslationRunFormField,
    NarrationResultRow,
    RunPhase,
    RunProgress
  } from "./translation-run-view"

  interface Props {
    form: TranslationRunForm
    phase: RunPhase
    canRun: boolean
    models: string[]
    modelsLoading: boolean
    results: NarrationResultRow[]
    errorMessage: string
    // 実行中の進捗。phase==="running" のときだけ表示する。未指定なら進捗バーを出さない。
    progress?: RunProgress
    onFieldInput: (field: TranslationRunFormField, value: string) => void
    onSelectPlugin: () => void
    onLoadModels: () => void
    onRun: () => void
  }

  let {
    form,
    phase,
    canRun,
    models,
    modelsLoading,
    results,
    errorMessage,
    progress,
    onFieldInput,
    onSelectPlugin,
    onLoadModels,
    onRun
  }: Props = $props()
</script>

<div class="min-h-screen w-full px-6 py-12 flex justify-center">
  <section class="w-full max-w-4xl flex flex-col gap-8" aria-labelledby="screen-title">
    <header class="flex flex-col gap-3">
      <span class="u-mono text-xs tracking-[0.32em] uppercase text-accent">
        Skyrim mod translation
      </span>
      <h1 id="screen-title" class="u-display text-4xl font-semibold text-base-content">
        翻訳実行
      </h1>
      <p class="max-w-2xl text-base-content/70 leading-relaxed">
        plugin を選んで叙述文を抽出し、AI 翻訳を実行して、原文と訳文を 1 つの画面で確かめます。
      </p>
      <div
        class="h-px w-full bg-gradient-to-r from-transparent via-primary/50 to-transparent"
      ></div>
    </header>

    <div class="card bg-base-200/55 border border-base-300/70 shadow-xl u-edge-top">
      <div class="card-body gap-6">
        <div class="flex flex-col gap-3">
          <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
            翻訳対象
          </h2>
          <FileSelectField
            label="plugin"
            buttonLabel="plugin を選択"
            path={form.pluginPath}
            emptyText="plugin ファイルが未選択です。"
            hint="master と Strings がある Data フォルダ内の plugin を選びます。Data フォルダはこのファイルの場所から判断します。"
            onSelect={onSelectPlugin}
          />
        </div>

        <div class="flex flex-col gap-4">
          <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
            AI サービス
          </h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-1">
            {#each PROVIDER_FIELDS as descriptor (descriptor.field)}
              <TextField
                label={descriptor.label}
                value={form[descriptor.field]}
                placeholder={descriptor.placeholder}
                hint={descriptor.hint}
                secret={descriptor.secret}
                oninput={(value: string) => onFieldInput(descriptor.field, value)}
              />
            {/each}
            <SelectField
              label="モデル"
              value={form.model}
              options={models}
              placeholder="モデルを選択"
              emptyText="未取得"
              hint="エンドポイントと API キーを入れてから取得します。"
              onChange={(value: string) => onFieldInput("model", value)}
            >
              {#snippet action()}
                <button
                  type="button"
                  class="btn btn-outline btn-primary"
                  disabled={modelsLoading}
                  onclick={onLoadModels}
                >
                  {#if modelsLoading}
                    <span class="loading loading-spinner loading-xs"></span>
                  {/if}
                  取得
                </button>
              {/snippet}
            </SelectField>
          </div>
        </div>

        {#if phase === "error" && errorMessage.length > 0}
          <div role="alert" class="alert alert-error alert-soft">
            <span>{errorMessage}</span>
          </div>
        {/if}

        <div class="card-actions items-center justify-end gap-4 pt-1">
          <StatusBadge
            label={PHASE_PRESENTATION[phase].label}
            tone={PHASE_PRESENTATION[phase].tone}
            loading={phase === "running"}
          />
          <button
            class="btn btn-primary px-8"
            type="button"
            disabled={!canRun || phase === "running"}
            onclick={onRun}
          >
            {phase === "running" ? "実行中…" : "実行"}
          </button>
        </div>
      </div>
    </div>

    {#if phase === "running" && progress}
      <TranslationProgress {...progress} />
    {/if}

    <ResultsPanel {phase} {results} />
  </section>
</div>
