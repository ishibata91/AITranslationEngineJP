<script lang="ts">
  // 翻訳実行画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TranslationRunScreen に委ね、本 component は state と配線だけを持つ。
  import { onMount } from "svelte"
  import TranslationRunScreen from "./TranslationRunScreen.svelte"
  import type {
    TranslationRunForm,
    TranslationRunFormField,
    NarrationResultRow,
    RunPhase,
    RunProgress,
    ResultsPaging
  } from "./translation-run-view"
  import {
    selectPluginFile,
    fetchModels,
    runExtractAndTranslate,
    listResultsPage,
    onRunProgress
  } from "../../../gateway/translation-gateway"

  // 結果一覧は keyset cursor ページングで 1 ページずつ取得する。1 ページの件数。
  const PAGE_SIZE = 50

  let pluginPath = $state("")
  let endpoint = $state("http://127.0.0.1:1234")
  let apiKey = $state("")
  let model = $state("")
  let models = $state<string[]>([])
  let modelsLoading = $state(false)
  let phase = $state<RunPhase>("idle")
  let results = $state<NarrationResultRow[]>([])
  let progress = $state<RunProgress | undefined>(undefined)
  let errorMessage = $state("")

  // keyset ページング state。cursorStack[i] はページ i を取得した cursor（"" 始まり）。
  // 順次送りのため、次へで nextCursor を積み、前へで履歴を 1 つ戻して再取得する。
  let total = $state(0)
  let cursorStack = $state<string[]>([""])
  let pageIndex = $state(0)
  let nextCursor = $state("")
  let hasMore = $state(false)

  const form: TranslationRunForm = $derived({ pluginPath, endpoint, apiKey, model })
  const canRun = $derived(
    pluginPath.length > 0 && endpoint.length > 0 && model.length > 0
  )

  const paging: ResultsPaging = $derived({
    total,
    pageNumber: pageIndex + 1,
    canPrev: pageIndex > 0,
    canNext: hasMore
  })

  // 指定 cursor のページを取得して現在ページへ反映する。
  async function loadPage(cursor: string) {
    const page = await listResultsPage(cursor, PAGE_SIZE)
    results = page.results
    total = page.total
    nextCursor = page.nextCursor
    hasMore = page.hasMore
  }

  // 先頭ページから読み直す（起動時・実行完了後）。
  async function resetToFirstPage() {
    cursorStack = [""]
    pageIndex = 0
    await loadPage("")
  }

  async function onPageNext() {
    if (!hasMore) return
    if (pageIndex === cursorStack.length - 1) {
      cursorStack = [...cursorStack, nextCursor]
    }
    pageIndex += 1
    try {
      await loadPage(cursorStack[pageIndex])
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  async function onPagePrev() {
    if (pageIndex === 0) return
    pageIndex -= 1
    try {
      await loadPage(cursorStack[pageIndex])
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  function onFieldInput(field: TranslationRunFormField, value: string) {
    if (field === "endpoint") endpoint = value
    else if (field === "apiKey") apiKey = value
    else if (field === "model") model = value
  }

  async function onSelectPlugin() {
    try {
      const path = await selectPluginFile()
      if (path.length > 0) pluginPath = path
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    }
  }

  async function onLoadModels() {
    modelsLoading = true
    errorMessage = ""
    try {
      models = await fetchModels({ endpoint, apiKey })
      if (model.length === 0 && models.length > 0) model = models[0]
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      modelsLoading = false
    }
  }

  async function onRun() {
    phase = "running"
    errorMessage = ""
    results = []
    progress = { stage: "extract", done: 0, total: 0 }
    try {
      await runExtractAndTranslate({ pluginPath, endpoint, apiKey, model })
      await resetToFirstPage()
      phase = "done"
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      progress = undefined
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
    return unsubscribe
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
  {onSelectPlugin}
  {onLoadModels}
  {onRun}
  {onPagePrev}
  {onPageNext}
/>
