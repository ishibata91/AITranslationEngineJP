// Wails generated bindings のラッパ。generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import {
  GetModels,
  ListResultsPage,
  RunExtractAndTranslate,
  SelectPluginFile
} from "../../wailsjs/go/api/App"
import { api } from "../../wailsjs/go/models"
import { EventsOn } from "../../wailsjs/runtime"

export interface Connection {
  endpoint: string
  apiKey: string
}

export interface RunInput {
  pluginPath: string
  endpoint: string
  apiKey: string
  model: string
}

// 結果一覧の 1 行。叙述文と台詞を共通の行で表す。directive と personaLabel は話者を解決できた台詞だけ持つ。
export interface ResultRow {
  edid: string
  source: string
  dest: string
  statusLabel: string
  directive?: string
  personaLabel?: string
}

// 実行結果の要約。結果一覧は数万件になりうるため実行時には返さず、listResultsPage で取得する。
export interface RunOutcome {
  translatedCount: number
}

// 結果一覧の keyset cursor ページ。total は総件数、nextCursor は次ページ取得用、hasMore は次ページの有無。
export interface ResultPage {
  total: number
  results: ResultRow[]
  nextCursor: string
  hasMore: boolean
}

// 本文翻訳の進捗。stage は extract（台詞抽出、不定）と translate（本文翻訳、done/total）。
export interface RunProgress {
  stage: "extract" | "translate"
  done: number
  total: number
}

// 進捗 event 名。backend の runtime event 名と一致させる。
const PROGRESS_EVENT = "translation:progress"

// plugin ファイル選択ダイアログを開き、選んだフルパスを返す。
export async function selectPluginFile(): Promise<string> {
  return SelectPluginFile()
}

// 接続先の利用可能モデル一覧を取得する（getModels）。
export async function fetchModels(conn: Connection): Promise<string[]> {
  return GetModels(api.ConnRequest.createFrom(conn))
}

// plugin を抽出し未訳を翻訳して、翻訳件数の要約を返す。結果一覧は listResultsPage で取得する。
export async function runExtractAndTranslate(input: RunInput): Promise<RunOutcome> {
  const result = await RunExtractAndTranslate(api.RunRequest.createFrom(input))
  return { translatedCount: result.translatedCount }
}

// 中心 DB の叙述文と台詞を keyset cursor ページで取得する（起動時・ページ送り・実行後の取得を統一）。
// cursor は ""（先頭）/ "n:<id>" / "l:<id>"。limit はページ件数。
export async function listResultsPage(
  cursor: string,
  limit: number
): Promise<ResultPage> {
  const page = await ListResultsPage(cursor, limit)
  return {
    total: page.total,
    results: page.results.map(toResultRow),
    nextCursor: page.nextCursor,
    hasMore: page.hasMore
  }
}

// 本文翻訳の進捗 event を購読する。返り値の関数を呼ぶと購読を解除する。
export function onRunProgress(handler: (progress: RunProgress) => void): () => void {
  return EventsOn(PROGRESS_EVENT, (ev: RunProgress) => handler(ev))
}

function toResultRow(view: api.ResultView): ResultRow {
  return {
    edid: view.edid,
    source: view.source,
    dest: view.dest,
    statusLabel: view.statusLabel,
    directive: view.directive,
    personaLabel: view.personaLabel
  }
}
