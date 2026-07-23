<script lang="ts">
  // 翻訳実行画面の container。画面状態（$state）を保持し、gateway 経由で backend を呼ぶ。
  // 表示は TranslationRunScreen に委ね、本 component は state と配線だけを持つ。
  import { onMount } from "svelte"
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
    APPLIED_BODY_NOTICE
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

  // 状態確認して進行が無かった時に出す案内（表示 copy は presentation に集約するが、
  // 「進行なし」用の定数は表示層に無いため、状態確認の結果としてここで確定文言を持つ）。
  const NO_PROGRESS_NOTICE = "この plugin に進行中の batch はありません。"

  // 配送方式ごとのエンドポイント既定値。sync はローカル OpenAI 互換サーバ、xai は xAI の接続先。
  // 方式を切り替えた時、利用者が触っていない（もう片方の既定のまま）なら新方式の既定へ寄せる。
  const SYNC_DEFAULT_ENDPOINT = "http://127.0.0.1:1234"
  const XAI_DEFAULT_ENDPOINT = "https://api.x.ai"

  // 結果一覧は keyset cursor ページングで 1 ページずつ取得する。1 ページの件数。
  const PAGE_SIZE = 50

  // 翻訳対象のプラグインは翻訳対象プラグイン画面で選び、ルーティングで渡される。ここは受け取るだけ。
  let { pluginPath = "" }: { pluginPath?: string } = $props()
  // 結果一覧を絞る対象 plugin 名（フルパスの末尾）。backend の plugin 列（filepath.Base）と一致させる。
  const pluginName = $derived(pluginPath.split(/[\\/]/).pop() ?? "")
  let endpoint = $state(SYNC_DEFAULT_ENDPOINT)
  let apiKey = $state("")
  let model = $state("")
  let models = $state<string[]>([])
  let modelsLoading = $state(false)
  let phase = $state<RunPhase>("idle")
  let results = $state<NarrationResultRow[]>([])
  let progress = $state<RunProgress | undefined>(undefined)
  let errorMessage = $state("")
  // 配送方式。sync=同期（既定）、xai=xAI の batch。取得先と起動操作（実行 / 状態確認・主アクション）を分ける。
  let provider = $state<TranslationProvider>("sync")
  // xAI の送信中・状態確認中・取り込み中フラグ。ボタンの無効化とスピナー表示に使う。
  let submitting = $state(false)
  let checking = $state(false)
  let applying = $state(false)
  // xAI batch の進行状況（状態確認で取得）。undefined は未確認。パネル表示と主アクションの活性に使う。
  let batchProgress = $state<BatchProgressView | undefined>(undefined)
  // xAI の送信直後・状態確認・取り込みの結果として出す案内。方式切替・実行・エラーでクリアする。
  let notice = $state("")
  // xTranslator 書き出し中フラグ。書き出しボタンの無効化とスピナー表示に使う。
  let exporting = $state(false)

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
  // xAI の送信可否。plugin とモデルが揃えば送れる（endpoint は空でも backend が xAI 既定を補う）。
  const canSubmit = $derived(
    provider === "xai" && pluginPath.length > 0 && model.length > 0 && !submitting
  )
  const paging: ResultsPaging = $derived({
    total,
    pageNumber: pageIndex + 1,
    canPrev: pageIndex > 0,
    canNext: hasMore
  })

  // 指定 cursor のページを取得して現在ページへ反映する。結果は選択中の plugin に絞る。
  async function loadPage(cursor: string) {
    const page = await listResultsPage(pluginName, cursor, PAGE_SIZE)
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

  async function onLoadModels() {
    modelsLoading = true
    errorMessage = ""
    try {
      // 取得先は配送方式で分ける。sync は OpenAI 互換、xai は xAI（batch 非対応モデルを除く）。
      models =
        provider === "xai"
          ? await fetchXaiModels({ endpoint, apiKey })
          : await fetchModels({ endpoint, apiKey })
      if (model.length === 0 && models.length > 0) model = models[0]
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      modelsLoading = false
    }
  }

  // 配送方式を切り替える。取得済みモデルは方式間で通用しないためクリアして取り直しを促す。
  // エンドポイントが切替前の方式の既定のままなら、新方式の既定へ寄せる（利用者が触った値は残す）。
  function onProviderChange(next: TranslationProvider) {
    if (next === provider) return
    if (provider === "sync" && endpoint === SYNC_DEFAULT_ENDPOINT) {
      endpoint = XAI_DEFAULT_ENDPOINT
    } else if (provider === "xai" && endpoint === XAI_DEFAULT_ENDPOINT) {
      endpoint = SYNC_DEFAULT_ENDPOINT
    }
    provider = next
    model = ""
    models = []
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
      lastPlugin = pluginName
      batchProgress = undefined
      notice = ""
      void loadStringsPresence(pluginPath)
    }
  })

  // xAI の batch を送信する（plugin 単位）。送信後は案内を出し、結果一覧は反映まで変えない。
  async function onSubmit() {
    submitting = true
    notice = ""
    errorMessage = ""
    try {
      await submitBatchTranslation({ pluginPath, endpoint, apiKey, model })
      notice = SUBMIT_NOTICE
      phase = "idle"
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      submitting = false
    }
  }

  // 対象 plugin の進行状況を確認する（観測・副作用なし）。dest 取り込みも送信もしない。
  // 進行が無ければ案内を出す。得た進行状況でパネルと主アクション（送信 / 取り込み）の活性が決まる。
  async function onCheckStatus() {
    checking = true
    notice = ""
    errorMessage = ""
    try {
      batchProgress = await getBatchProgress(pluginName, { endpoint, apiKey })
      if (!batchProgress) notice = NO_PROGRESS_NOTICE
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      checking = false
    }
  }

  // 対象 plugin の完了段を取り込む（前進）。固有名段なら本文 batch も送られる。
  // 取り込み後は結果一覧を読み直し、進行状況を取り直して次段（本文 / 完了）へ更新する。
  // 案内は取り込んだ段（取り込み前の段）で選ぶ。接続情報は都度使い、永続化しない。
  async function onApply() {
    const appliedStage = batchProgress?.stage
    applying = true
    notice = ""
    errorMessage = ""
    try {
      await refreshBatchTranslations(pluginName, { endpoint, apiKey })
      await resetToFirstPage()
      batchProgress = await getBatchProgress(pluginName, { endpoint, apiKey })
      notice =
        appliedStage === "proper" ? APPLIED_PROPER_NOTICE : APPLIED_BODY_NOTICE
      phase = "done"
    } catch (error) {
      errorMessage = messageOf(error)
      phase = "error"
    } finally {
      applying = false
    }
  }

  async function onRun() {
    phase = "running"
    errorMessage = ""
    notice = ""
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
  {onLoadModels}
  {onRun}
  {onPagePrev}
  {onPageNext}
  {exporting}
  {onExportXml}
  {provider}
  {onProviderChange}
  {onSubmit}
  {canSubmit}
  {submitting}
  {batchProgress}
  {onCheckStatus}
  {checking}
  {onApply}
  {applying}
  {notice}
  {stringsPresence}
/>
