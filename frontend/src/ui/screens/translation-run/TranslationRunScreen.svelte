<script lang="ts">
  // 翻訳実行画面。汎用の表示部品を組み立てるだけにし、state・API・Wails・validation は持たない。
  // すべて props で受け、入力 event は callback prop で上へ返す。
  import StatusBadge from "@ui/components/StatusBadge.svelte"
  import TextField from "@ui/components/TextField.svelte"
  import SelectField from "@ui/components/SelectField.svelte"
  import ResultsPanel from "./ResultsPanel.svelte"
  import MissingStringsWarning from "./MissingStringsWarning.svelte"
  import TranslationProgress from "./TranslationProgress.svelte"
  import BatchProgressPanel from "./BatchProgressPanel.svelte"
  import {
    providerFields,
    PROVIDER_OPTIONS,
    PROVIDER_LABEL,
    MODEL_HINT,
    BATCH_ACTION_HINT,
    batchMainAction,
    PHASE_PRESENTATION
  } from "./translation-run-presentation"
  import type {
    TranslationRunForm,
    TranslationRunFormField,
    TranslationProvider,
    BatchProgressView,
    NarrationResultRow,
    RunPhase,
    RunProgress,
    ResultsPaging,
    StringsPresence
  } from "./translation-run-view"

  interface Props {
    form: TranslationRunForm
    phase: RunPhase
    canRun: boolean
    models: string[]
    modelsLoading: boolean
    // 現在ページの結果行。ページングするため全件ではない。
    results: NarrationResultRow[]
    errorMessage: string
    // 実行中の進捗。phase==="running" のときだけ表示する。未指定なら進捗バーを出さない。
    progress?: RunProgress
    // 結果一覧のページング表示値。未指定なら単一ページ（前後無効）として扱う。
    paging?: ResultsPaging
    // 翻訳対象 plugin の Strings の言語別有無。未指定（未判定）なら警告を出さない。
    stringsPresence?: StringsPresence
    onFieldInput: (field: TranslationRunFormField, value: string) => void
    onLoadModels: () => void
    onRun: () => void
    // ページ送り操作。state は container が持つ。
    onPagePrev?: () => void
    onPageNext?: () => void
    // xTranslator XML への書き出し操作。結果一覧パネルの書き出しボタンから起動する。
    onExportXml?: () => void
    // 書き出し中フラグ。true の間は書き出しボタンを無効化する。
    exporting?: boolean
    // 翻訳の配送方式。sync=同期（既定）、xai=xAI の batch。未指定なら同期で従来表示。
    provider?: TranslationProvider
    // 配送方式の切り替え操作。
    onProviderChange?: (provider: TranslationProvider) => void
    // xAI の batch 送信操作（plugin 単位）。
    onSubmit?: () => void
    // 送信の可否。接続情報とモデルが揃ったら true。
    canSubmit?: boolean
    // 送信中フラグ。true の間は送信ボタンを無効化しスピナーを出す。
    submitting?: boolean
    // xAI の送信・取り込みの結果として出す案内。空なら出さない。
    notice?: string
    // xAI batch の進行状況（状態確認で取得）。未指定なら未確認としてパネルを控えめに出す。
    batchProgress?: BatchProgressView
    // xAI の状態確認操作（副作用なし。進行状況を最新化する）。
    onCheckStatus?: () => void
    // 状態確認中フラグ。true の間は状態確認ボタンにスピナーを出す。
    checking?: boolean
    // xAI の取り込み操作（完了段を dest へ取り込み、次段 batch を送る）。
    onApply?: () => void
    // 取り込み中フラグ。true の間は取り込みボタンにスピナーを出す。
    applying?: boolean
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
    paging,
    stringsPresence = undefined,
    onFieldInput,
    onLoadModels,
    onRun,
    onPagePrev = () => {},
    onPageNext = () => {},
    onExportXml = () => {},
    exporting = false,
    provider = "sync",
    onProviderChange = () => {},
    onSubmit = () => {},
    canSubmit = false,
    submitting = false,
    notice = "",
    batchProgress = undefined,
    onCheckStatus = () => {},
    checking = false,
    onApply = () => {},
    applying = false
  }: Props = $props()

  // 配送方式に応じた接続情報欄（sync は既定、xai は xAI 用の読み替え欄）。
  const fields = $derived(providerFields(provider))

  // xAI の主アクション（送信と取り込みを 1 ボタンに束ねた表示）。進行状況で種別・ラベル・活性が変わる。
  const mainAction = $derived(batchMainAction(batchProgress, canSubmit))
  // 主アクションの実行中フラグ（送信中 / 取り込み中）。busy の間はスピナーを出し無効化する。
  const mainActionBusy = $derived(mainAction.kind === "send" ? submitting : applying)
  const mainActionDisabled = $derived(!mainAction.enabled || mainActionBusy)

  // 翻訳対象のプラグイン名（フルパスの末尾）。選択は翻訳対象プラグイン画面で行い、ここは表示専用。
  const pluginName = $derived.by(() => {
    const segments = form.pluginPath.split(/[\\/]/)
    return segments[segments.length - 1] ?? ""
  })
</script>

<div class="min-h-screen w-full px-6 py-12 flex justify-center">
  <section class="w-full max-w-4xl flex flex-col gap-8" aria-labelledby="screen-title">
    <header class="flex flex-col gap-3">
      <span class="u-mono text-xs tracking-[0.32em] text-accent">
        翻訳の実行
      </span>
      <h1 id="screen-title" class="u-display text-4xl font-semibold text-base-content">
        翻訳実行
      </h1>
      <p class="max-w-2xl text-base-content/70 leading-relaxed">
        選んだプラグインを AI 翻訳し、原文と訳文を 1 つの画面で確かめます。翻訳対象は翻訳対象プラグイン画面で選びます。
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
          {#if form.pluginPath.length > 0}
            <div class="flex flex-col gap-1">
              <span class="u-mono text-sm text-base-content">{pluginName}</span>
              <span class="u-mono text-xs text-base-content/45 truncate">{form.pluginPath}</span>
            </div>
          {:else}
            <span class="text-base-content/45">翻訳対象プラグインが選ばれていません。</span>
          {/if}
          <MissingStringsWarning presence={stringsPresence} />
        </div>

        <div class="flex flex-col gap-4">
          <h2 class="u-display text-sm tracking-widest uppercase text-base-content/60">
            AI サービス
          </h2>
          <div class="join self-start" role="group" aria-label="配送方式">
            {#each PROVIDER_OPTIONS as option (option)}
              <button
                type="button"
                class="btn join-item {provider === option
                  ? 'btn-primary'
                  : 'btn-outline btn-primary'}"
                aria-pressed={provider === option}
                onclick={() => onProviderChange(option)}
              >
                {PROVIDER_LABEL[option]}
              </button>
            {/each}
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-1">
            {#each fields as descriptor (descriptor.field)}
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
              hint={MODEL_HINT[provider]}
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

        {#if provider === "xai" && notice.length > 0}
          <div role="status" class="alert alert-info alert-soft">
            <span>{notice}</span>
          </div>
        {/if}

        {#if phase === "error" && errorMessage.length > 0}
          <div role="alert" class="alert alert-error alert-soft">
            <span>{errorMessage}</span>
          </div>
        {/if}

        {#if provider === "xai"}
          <BatchProgressPanel progress={batchProgress} />
          <p class="text-xs text-base-content/50">{BATCH_ACTION_HINT}</p>
        {/if}

        <div class="card-actions items-center justify-end gap-4 pt-1">
          <StatusBadge
            label={PHASE_PRESENTATION[phase].label}
            tone={PHASE_PRESENTATION[phase].tone}
            loading={phase === "running"}
          />
          {#if provider === "sync"}
            <button
              class="btn btn-primary px-8"
              type="button"
              disabled={!canRun || phase === "running"}
              onclick={onRun}
            >
              {phase === "running" ? "実行中…" : "実行"}
            </button>
          {:else}
            <button
              class="btn btn-outline btn-primary px-6"
              type="button"
              disabled={checking}
              onclick={onCheckStatus}
            >
              {#if checking}
                <span class="loading loading-spinner loading-xs"></span>
              {/if}
              状態確認
            </button>
            <button
              class="btn btn-primary px-8"
              type="button"
              disabled={mainActionDisabled}
              onclick={mainAction.kind === "send" ? onSubmit : onApply}
            >
              {#if mainActionBusy}
                <span class="loading loading-spinner loading-xs"></span>
              {/if}
              {mainAction.label}
            </button>
          {/if}
        </div>
      </div>
    </div>

    {#if phase === "running" && progress}
      <TranslationProgress {...progress} />
    {/if}

    <ResultsPanel
      {phase}
      {results}
      total={paging?.total}
      pageNumber={paging?.pageNumber}
      canPrev={paging?.canPrev}
      canNext={paging?.canNext}
      onPrev={onPagePrev}
      onNext={onPageNext}
      onExport={onExportXml}
      {exporting}
    />
  </section>
</div>
