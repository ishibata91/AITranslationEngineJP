// Wails generated bindings のラッパ。generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import {
  GetModels,
  ListNarrations,
  RunExtractAndTranslate,
  SelectPluginFile
} from "../../wailsjs/go/api/App"
import { api } from "../../wailsjs/go/models"

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

export interface NarrationResult {
  edid: string
  source: string
  dest: string
  statusLabel: string
}

export interface RunOutcome {
  translatedCount: number
  narrations: NarrationResult[]
}

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
    narrations: result.narrations.map(toNarrationResult)
  }
}

// 中心 DB の現在の叙述文を取得する（起動時に前回の結果を表示するため）。
export async function listNarrations(): Promise<NarrationResult[]> {
  const rows = await ListNarrations()
  return rows.map(toNarrationResult)
}

function toNarrationResult(view: api.NarrationView): NarrationResult {
  return {
    edid: view.edid,
    source: view.source,
    dest: view.dest,
    statusLabel: view.statusLabel
  }
}
