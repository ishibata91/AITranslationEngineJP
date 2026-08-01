<script lang="ts">
  // 翻訳実行画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TranslationRunScreen に委ね、本 component は state と配線だけを持つ。
  import { onMount, tick } from "svelte"
  import TranslationRunScreen from "./TranslationRunScreen.svelte"
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
  import {
    SUBMIT_NOTICE,
    APPLIED_PROPER_NOTICE,
    APPLIED_BODY_NOTICE,
    batchUntranslatedNotice,
    untranslatedNotice
  } from "./translation-run-presentation"
  import {
    fetchModels,
    fetchStringsPresence,
    fetchXaiModels,
    runExtractAndTranslate,
    submitBatchTranslation,
    refreshBatchTranslations,
    getBatchProgress,
    listResultsPage,
    onRunProgress,
    exportXTranslatorXml
  } from "../../../gateway/translation-gateway"
  import {
    canOperateBatch as batchOperationAllowed,
    orderProviderModels,
    providerDefaults
  } from "./translation-run-provider"

  // 結果一覧は keyset cursor ページングで 1 ページずつ取得する。1 ページの件数。
  const PAGE_SIZE = 50
  // 前の状態確認が完了してから次の状態確認を始めるまでの間隔。
  const BATCH_POLL_INTERVAL_MS = 10_000

  // 翻訳対象のプラグインは翻訳対象プラグイン画面で選び、ルーティングで渡される。ここは受け取るだけ。
  let { pluginPath = "" }: { pluginPath?: string } = $props()
  // 結果一覧を絞る対象 plugin 名（フルパスの末尾）。backend の plugin 列（filepath.Base）と一致させる。
  const pluginName = $derived(pluginPath.split(/[\\/]/).pop() ?? "")
  let endpoint = $state(providerDefaults("sync").endpoint)
  let apiKey = $state("")
  let model = $state("")
  let models = $state<string[]>([])
  let modelsLoading = $state(false)
  let phase = $state<RunPhase>("idle")
  let results = $state<NarrationResultRow[]>([])
  let progress = $state<RunProgress | undefined>(undefined)
  let errorMessage = $state("")
  // 配送方式。sync=同期（既定）、openai=OpenAI batch、xai=xAI batch。
  let provider = $state<TranslationProvider>("sync")
  // batch の開始処理と自動状態確認を一つの実行として扱い、多重実行を防ぐ。
  let batchRunning = $state(false)
  // batch の進行状況（状態確認で取得）。undefined は未確認。パネル表示と主アクションの活性に使う。
  let batchProgress = $state<BatchProgressView | undefined>(undefined)
  // batch の送信直後・状態確認・取り込みの結果として出す案内。方式切替・実行・エラーでクリアする。
  let notice = $state("")
  // xTranslator 書き出し中フラグ。書き出しボタンの無効化とスピナー表示に使う。
  let exporting = $state(false)

  // keyset ページング state。cursorStack[i] はページ i を取得した cursor（"" 始まり）。
  // 順次送りのため、次へで nextCursor を積み、前へで履歴を 1 つ戻して再取得する。
  let total = $state(0)
  let unfilteredTotal = $state(0)
  let untranslatedOnly = $state(false)
  let cursorStack = $state<string[]>([""])
  let pageIndex = $state(0)
  let nextCursor = $state("")
  let hasMore = $state(false)

  const form: TranslationRunForm = $derived({
    pluginPath,
    endpoint,
    apiKey,
    model
  })
  const canRun = $derived(
    pluginPath.length > 0 && endpoint.length > 0 && model.length > 0
  )
  // OpenAI は API キーを必須にし、送信・状態確認・取り込みの全操作を同じ条件で止める。
  const canOperateBatch = $derived(batchOperationAllowed(provider, apiKey))
  const canSubmit = $derived(
    provider !== "sync" &&
      canOperateBatch &&
      pluginPath.length > 0 &&
      model.length > 0 &&
      !batchRunning
  )
  const paging: ResultsPaging = $derived({
    total,
    pageNumber: pageIndex + 1,
    canPrev: pageIndex > 0,
    canNext: hasMore
  })
  const hasUnfilteredResults = $derived(unfilteredTotal > 0)

  type LoadedPage = Awaited<ReturnType<typeof listResultsPage>>
  type BatchProvider = Exclude<TranslationProvider, "sync">
  interface BatchLoopContext {
    generation: number
    pluginPath: string
    pluginName: string
    provider: BatchProvider
    endpoint: string
    apiKey: string
    model: string
  }

  let batchTimer: ReturnType<typeof setTimeout> | undefined
  let batchGeneration = 0

  // 指定 cursor と条件のページを一時値へ取得する。呼び出し側が必要な取得を終えてから画面へ反映する。
  async function fetchPage(cursor: string, onlyUntranslated: boolean) {
    return listResultsPage(pluginName, cursor, PAGE_SIZE, onlyUntranslated)
  }

  async function fetchPageForPlugin(
    targetPlugin: string,
    cursor: string,
    onlyUntranslated: boolean
  ) {
    return listResultsPage(targetPlugin, cursor, PAGE_SIZE, onlyUntranslated)
  }

  function applyPage(page: LoadedPage) {
    results = page.results
    total = page.total
    unfilteredTotal = page.unfilteredTotal
    nextCursor = page.nextCursor
    hasMore = page.hasMore
  }

  // 先頭ページから読み直す（起動時・実行完了後）。
  async function resetToFirstPage(onlyUntranslated = untranslatedOnly) {
    const page = await fetchPage("", onlyUntranslated)
    applyPage(page)
    cursorStack = [""]
    pageIndex = 0
  }

  async function onPageNext() {
    if (!hasMore) return
    const targetIndex = pageIndex + 1
    const targetCursor = cursorStack[targetIndex] ?? nextCursor
    try {
      const page = await fetchPage(targetCursor, untranslatedOnly)
      applyPage(page)
      if (targetIndex === cursorStack.length) {
        cursorStack = [...cursorStack, targetCursor]
      }
      pageIndex = targetIndex
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  async function onPagePrev() {
    if (pageIndex === 0) return
    const targetIndex = pageIndex - 1
    try {
      const page = await fetchPage(cursorStack[targetIndex], untranslatedOnly)
      applyPage(page)
      pageIndex = targetIndex
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  // 未訳条件の先頭ページを取得できた場合だけ、条件・一覧・cursor・ページ番号をまとめて切り替える。
  async function onUntranslatedOnlyChange(checked: boolean) {
    if (checked === untranslatedOnly) return
    errorMessage = ""
    try {
      const page = await fetchPage("", checked)
      applyPage(page)
      untranslatedOnly = checked
      cursorStack = [""]
      pageIndex = 0
    } catch (error) {
      // checkbox 自身が先に切り替えた DOM を、保持した state の値へ戻すために 1 度差分を作る。
      untranslatedOnly = checked
      await tick()
      untranslatedOnly = !checked
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  function onFieldInput(field: TranslationRunFormField, value: string) {
    if (field === "endpoint") endpoint = value
    else if (field === "apiKey") apiKey = value
    else if (field === "model") model = value
  }

  async function onLoadModels() {
    modelsLoading = true
    errorMessage = ""
    try {
      // OpenAI は Luna が取得結果にあれば先頭へ置き、他モデルは残す。
      if (provider === "xai") {
        models = await fetchXaiModels({ endpoint, apiKey })
      } else {
        models = orderProviderModels(
          provider,
          await fetchModels({ endpoint, apiKey })
        )
      }
      if (models.length > 0) model = models[0]
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      modelsLoading = false
    }
  }

  function invalidateBatchLoop() {
    batchGeneration += 1
    if (batchTimer !== undefined) {
      clearTimeout(batchTimer)
      batchTimer = undefined
    }
    batchRunning = false
  }

  function isCurrentBatchLoop(generation: number) {
    return generation === batchGeneration
  }

  // 配送方式を切り替える。切替前の自動状態確認を止め、選択先の既定値へ替える。
  function onProviderChange(next: TranslationProvider) {
    if (next === provider) return
    invalidateBatchLoop()
    provider = next
    const defaults = providerDefaults(next)
    endpoint = defaults.endpoint
    model = defaults.model
    models = defaults.models
    notice = ""
    batchProgress = undefined
    errorMessage = ""
    phase = "idle"
  }

  // 対象 plugin の Data フォルダにある Strings の言語別有無。undefined は未判定（警告を出さない）。
  // 片側欠けだと既存訳（参照訳・固有名の確定訳語）の対を作れないため、画面警告の判定材料として取得する。
  let stringsPresence = $state<StringsPresence | undefined>(undefined)

  async function loadStringsPresence(path: string) {
    if (path.length === 0) {
      stringsPresence = undefined
      return
    }
    try {
      stringsPresence = await fetchStringsPresence(path)
    } catch {
      // 判定に失敗しても翻訳自体は進められるため、未判定（警告なし）のまま続行する。
      stringsPresence = undefined
    }
  }

  // 対象 plugin が変わったら進行状況と案内をクリアし（別 plugin の状態を持ち越さない）、Strings の有無を判定し直す。
  let lastPlugin = $state("")
  $effect(() => {
    if (pluginName !== lastPlugin) {
      invalidateBatchLoop()
      lastPlugin = pluginName
      batchProgress = undefined
      notice = ""
      void loadStringsPresence(pluginPath)
    }
  })

  function finishBatchLoop(
    context: BatchLoopContext,
    progress: BatchProgressView,
    completionNotice = ""
  ) {
    if (!isCurrentBatchLoop(context.generation)) return
    if (batchTimer !== undefined) {
      clearTimeout(batchTimer)
      batchTimer = undefined
    }
    batchProgress = progress
    batchRunning = false
    phase = "done"
    notice =
      completionNotice || batchUntranslatedNotice(progress.untranslatedCount)
  }

  function failBatchLoop(context: BatchLoopContext, error: unknown) {
    if (!isCurrentBatchLoop(context.generation)) return
    if (batchTimer !== undefined) {
      clearTimeout(batchTimer)
      batchTimer = undefined
    }
    batchRunning = false
    errorMessage = messageOf(error)
    phase = "error"
  }

  function scheduleBatchPoll(context: BatchLoopContext) {
    if (!isCurrentBatchLoop(context.generation)) return
    if (batchTimer !== undefined) clearTimeout(batchTimer)
    batchTimer = setTimeout(() => {
      batchTimer = undefined
      void pollBatch(context)
    }, BATCH_POLL_INTERVAL_MS)
  }

  async function applyCompletedBatch(
    context: BatchLoopContext,
    progressBeforeApply: BatchProgressView
  ) {
    await refreshBatchTranslations(context.pluginName, context.provider, {
      endpoint: context.endpoint,
      apiKey: context.apiKey
    })
    if (!isCurrentBatchLoop(context.generation)) return

    const [page, nextProgress] = await Promise.all([
      fetchPageForPlugin(context.pluginName, "", untranslatedOnly),
      getBatchProgress(context.pluginName, context.provider, {
        endpoint: context.endpoint,
        apiKey: context.apiKey
      })
    ])
    if (!isCurrentBatchLoop(context.generation)) return
    if (!nextProgress) {
      throw new Error("取り込み後の batch 状態を取得できませんでした。")
    }

    applyPage(page)
    cursorStack = [""]
    pageIndex = 0
    batchProgress = nextProgress
    if (
      progressBeforeApply.stage === "proper" &&
      nextProgress.stage === "body"
    ) {
      notice = APPLIED_PROPER_NOTICE
    } else if (nextProgress.stage === "done") {
      notice =
        batchUntranslatedNotice(nextProgress.untranslatedCount) ||
        APPLIED_BODY_NOTICE
    } else {
      notice = ""
    }

    if (nextProgress.stage === "done") {
      finishBatchLoop(context, nextProgress, notice)
      return
    }
    scheduleBatchPoll(context)
  }

  async function continueBatchLoop(
    context: BatchLoopContext,
    progress: BatchProgressView
  ) {
    if (!isCurrentBatchLoop(context.generation)) return
    batchProgress = progress
    if (progress.stage === "done") {
      finishBatchLoop(context, progress)
      return
    }
    if (progress.canApply) {
      await applyCompletedBatch(context, progress)
      return
    }
    scheduleBatchPoll(context)
  }

  async function pollBatch(context: BatchLoopContext) {
    if (!isCurrentBatchLoop(context.generation)) return
    try {
      const progress = await getBatchProgress(
        context.pluginName,
        context.provider,
        { endpoint: context.endpoint, apiKey: context.apiKey }
      )
      if (!isCurrentBatchLoop(context.generation)) return
      if (!progress) {
        throw new Error("自動確認中の batch 状態を取得できませんでした。")
      }
      await continueBatchLoop(context, progress)
    } catch (error) {
      failBatchLoop(context, error)
    }
  }

  // 一つの入口で、新規送信、保存済み進行の再開、完了後の未訳再送信を開始する。
  async function onSubmit() {
    if (provider === "sync" || !canOperateBatch || batchRunning) return

    invalidateBatchLoop()
    const context: BatchLoopContext = {
      generation: batchGeneration,
      pluginPath,
      pluginName,
      provider,
      endpoint,
      apiKey,
      model
    }
    batchRunning = true
    notice = ""
    errorMessage = ""
    phase = "idle"
    let submitted = false

    try {
      const savedProgress = await getBatchProgress(
        context.pluginName,
        context.provider,
        { endpoint: context.endpoint, apiKey: context.apiKey }
      )
      if (!isCurrentBatchLoop(context.generation)) return
      batchProgress = savedProgress

      if (savedProgress && savedProgress.stage !== "done") {
        await continueBatchLoop(context, savedProgress)
        return
      }
      if (
        savedProgress?.stage === "done" &&
        savedProgress.untranslatedCount === 0
      ) {
        finishBatchLoop(context, savedProgress)
        return
      }

      const outcome = await submitBatchTranslation({
        pluginPath: context.pluginPath,
        endpoint: context.endpoint,
        apiKey: context.apiKey,
        model: context.model,
        provider: context.provider
      })
      submitted = true
      if (!isCurrentBatchLoop(context.generation)) return

      const nextProgress = await getBatchProgress(
        context.pluginName,
        context.provider,
        { endpoint: context.endpoint, apiKey: context.apiKey }
      )
      if (!isCurrentBatchLoop(context.generation)) return
      if (!nextProgress) {
        throw new Error("送信後の batch 状態を取得できませんでした。")
      }

      if (outcome.completedWithoutExternalBatch) {
        const page = await fetchPageForPlugin(
          context.pluginName,
          "",
          untranslatedOnly
        )
        if (!isCurrentBatchLoop(context.generation)) return
        applyPage(page)
        cursorStack = [""]
        pageIndex = 0
      }

      batchProgress = nextProgress
      const reusedNotice = outcome.reusedPreparation
        ? "保存済みの準備を使って未訳だけを処理しました。"
        : ""
      notice = [reusedNotice, SUBMIT_NOTICE].filter(Boolean).join(" ")
      if (nextProgress.stage === "done") {
        const completionNotice =
          batchUntranslatedNotice(nextProgress.untranslatedCount) ||
          APPLIED_BODY_NOTICE
        finishBatchLoop(
          context,
          nextProgress,
          [reusedNotice, completionNotice].filter(Boolean).join(" ")
        )
        return
      }
      scheduleBatchPoll(context)
    } catch (error) {
      if (!isCurrentBatchLoop(context.generation)) return
      const failure = submitted
        ? new Error(
            `batch の送信は完了しましたが、画面の更新に失敗しました: ${messageOf(error)}`
          )
        : error
      failBatchLoop(context, failure)
    }
  }

  async function onRun() {
    phase = "running"
    errorMessage = ""
    notice = ""
    results = []
    progress = undefined
    try {
      const outcome = await runExtractAndTranslate({
        pluginPath,
        endpoint,
        apiKey,
        model,
        provider
      })
      await resetToFirstPage()
      // 未訳のまま残った件数を案内として出す。0 件なら untranslatedNotice が空文字を返し、案内は出ない。
      notice = untranslatedNotice(outcome.untranslatedCount)
      phase = "done"
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      progress = undefined
    }
  }

  // 訳出済みの翻訳結果を xTranslator XML へ書き出す。出力先フォルダ選択と結果通知（成功の出力先・失敗）は
  // backend のネイティブダイアログで行うため、ここでは書き出し中フラグの制御だけを持つ。
  async function onExportXml() {
    exporting = true
    try {
      await exportXTranslatorXml()
    } catch {
      // 失敗は backend のネイティブダイアログで通知済みのため、二重表示を避けてここでは握る。
    } finally {
      exporting = false
    }
  }

  function messageOf(error: unknown): string {
    if (error instanceof Error) return error.message
    if (typeof error === "string") return error
    return "予期しないエラーが発生しました。"
  }

  // 起動時に前回の結果を先頭ページから読み込み、本文翻訳の進捗 event を購読する。
  async function loadPrevious() {
    try {
      await resetToFirstPage()
    } catch {
      // 起動時の読み込み失敗は致命的でないため、空のまま続行する。
    }
  }

  onMount(() => {
    const unsubscribe = onRunProgress((p) => {
      progress = p
    })
    void loadPrevious()
    return () => {
      invalidateBatchLoop()
      unsubscribe()
    }
  })
</script>

<TranslationRunScreen
  {form}
  {phase}
  {canRun}
  {models}
  {modelsLoading}
  {results}
  {progress}
  {paging}
  {errorMessage}
  {onFieldInput}
  {onLoadModels}
  {onRun}
  {onPagePrev}
  {onPageNext}
  {untranslatedOnly}
  {onUntranslatedOnlyChange}
  {hasUnfilteredResults}
  {exporting}
  {onExportXml}
  {provider}
  {onProviderChange}
  {onSubmit}
  {canSubmit}
  {canOperateBatch}
  {batchProgress}
  {batchRunning}
  {notice}
  {stringsPresence}
/>
