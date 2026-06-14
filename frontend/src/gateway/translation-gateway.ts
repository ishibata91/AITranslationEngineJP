// Wails generated bindings のラッパ。generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import {
  GetModels,
  ListResults,
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

export interface RunOutcome {
  translatedCount: number
  results: ResultRow[]
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

// plugin を抽出し未訳を翻訳して、結果一覧を返す。
export async function runExtractAndTranslate(input: RunInput): Promise<RunOutcome> {
  const result = await RunExtractAndTranslate(api.RunRequest.createFrom(input))
  return {
    translatedCount: result.translatedCount,
    results: result.results.map(toResultRow)
  }
}

// 中心 DB の現在の叙述文と台詞を取得する（起動時に前回の結果を表示するため）。
export async function listResults(): Promise<ResultRow[]> {
  const rows = await ListResults()
  return rows.map(toResultRow)
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
