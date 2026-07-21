// Wails generated bindings のラッパ。generated wailsjs の import はこの gateway 境界にだけ閉じ込める。
// View・Container からは本 gateway 経由で backend を呼ぶ。
import {
  DeleteTargetPlugin,
  ExportXTranslatorXML,
  GetModels,
  GetXAIModels,
  ListResultsPage,
  ListTargetPlugins,
  RefreshBatchTranslations,
  RunExtractAndTranslate,
  SelectPluginFile,
  SubmitBatchTranslation
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

// 結果行の機械置換内訳 1 件（原語 → 確定訳語）。取得時に原文へ辞書を当て直して再構成した値。
export interface TermRow {
  source: string
  dest: string
}

// 結果行の口調メタデータ。判定された基底口調（cell・trait）と根拠（段階・印・経路・性別）。口調が付いた台詞だけ持つ。
// 名指し話者は cell・attitudeBand・marked を持つ。汎用台詞・PC 発話はそれらが空（cell="" / 0）で、emotionBand・sex を持つ。
export interface PersonaRow {
  cell: string
  trait: string
  attitudeBand: 0 | 1 | 2
  emotionBand: 0 | 1 | 2
  marked: number
  decisionPath: "本文" | "voice" | "保留" | "汎用" | "PC"
  sex: string
}

// 結果行の話者識別と属性。誰の台詞かと、口調指示の根拠（性別・年齢・声型）。話者を解決できた台詞だけ持つ。
export interface SpeakerRow {
  edid: string
  sex: string
  age: string
  voice: string
}

// 結果行の元レコード種別バッジ。box は概念モデルの箱（叙述文/固有名/定型句/台詞）、recField は "REC:FIELD"。
export interface RecordTypeRow {
  box: string
  recField: string
}

// 結果一覧の 1 行。叙述文・定型句・台詞・固有名を共通の行で表す。speaker・directive・personaLabel・persona は話者を解決できた台詞だけ持つ。
// terms は本文で辞書から確定訳語へ置換した固有名の内訳。置換が無い行は省略する。
// prompt は実際に翻訳 AI へ投げた完成プロンプトを取得時に再構成した全文。再構成できない行は省略する。
// recordType は元レコード種別バッジ。master に無い種別の行は省略する。
export interface ResultRow {
  edid: string
  source: string
  dest: string
  statusLabel: string
  recordType?: RecordTypeRow
  speaker?: SpeakerRow
  directive?: string
  personaLabel?: string
  persona?: PersonaRow
  terms?: TermRow[]
  prompt?: string
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

// 翻訳対象プラグイン一覧の 1 行。plugin=ファイル名（識別子）、sourcePath=選んだフルパス（再実行に使う）、
// createdLabel=作成時刻の表示ラベル、total=束ねる翻訳対象の総数、translated=訳済み数。
export interface TargetPluginSummary {
  plugin: string
  sourcePath: string
  createdLabel: string
  total: number
  translated: number
}

// plugin ファイル選択ダイアログを開き、選んだフルパスを返す（選択なしは空文字）。
export async function selectPluginFile(): Promise<string> {
  return SelectPluginFile()
}

// 翻訳した対象 plugin を新しい順で取得する（一覧・進捗・削除・再実行の入口）。
// 作成時刻（backend の "YYYY-MM-DD HH:MM:SS"）は分までの表示ラベルへ整える。
export async function listTargetPlugins(): Promise<TargetPluginSummary[]> {
  const rows = await ListTargetPlugins()
  return rows.map((r) => ({
    plugin: r.plugin,
    sourcePath: r.sourcePath,
    createdLabel: r.createdAt.slice(0, 16),
    total: r.total,
    translated: r.translated
  }))
}

// 対象 plugin の翻訳成果を消す（対象スコープの行と連関を消し、共有辞書・共有 entity は残す）。
export async function deleteTargetPlugin(plugin: string): Promise<void> {
  await DeleteTargetPlugin(plugin)
}

// 訳出済みの翻訳結果を xTranslator 互換 XML として書き出す。
// 出力先フォルダ選択と結果通知（成功の出力先・失敗）は backend のネイティブダイアログで行うため、
// この関数は完了、または失敗時の throw だけを返す。
export async function exportXTranslatorXml(): Promise<void> {
  await ExportXTranslatorXML()
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

// xAI（batch）の接続先で利用可能なモデル一覧を取得する（getXAIModels）。
// backend 側で batch 非対応モデルを除く。endpoint が空なら backend が xAI 既定接続先を補う。
export async function fetchXaiModels(conn: Connection): Promise<string[]> {
  return GetXAIModels(api.ConnRequest.createFrom(conn))
}

// xAI の batch 翻訳を送信する（plugin 単位）。送信だけを行い、結果は後の反映で取り込む。
// 外部 batch ID の永続まで backend が行うため、ここは送信の完了、または失敗時の throw だけを返す。
export async function submitBatchTranslation(input: RunInput): Promise<void> {
  await SubmitBatchTranslation(api.RunRequest.createFrom(input))
}

// 進行中の全 batch をまとめて反映する。完了した batch の結果を dest へ取り込む。
// plugin 引数を取らない global 操作で、接続情報は反映のたびに渡す（永続化しない）。
export async function refreshBatchTranslations(conn: Connection): Promise<void> {
  await RefreshBatchTranslations(api.ConnRequest.createFrom(conn))
}

// 中心 DB の叙述文と台詞を keyset cursor ページで取得する（起動時・ページ送り・実行後の取得を統一）。
// plugin が空でなければその対象 plugin（plugin ファイル名）の結果だけに絞る。空なら全 plugin。
// cursor は ""（先頭）/ "n:<id>" / "l:<id>"。limit はページ件数。
export async function listResultsPage(
  plugin: string,
  cursor: string,
  limit: number
): Promise<ResultPage> {
  const page = await ListResultsPage(plugin, cursor, limit)
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
    recordType: view.recordType
      ? { box: view.recordType.box, recField: view.recordType.recField }
      : undefined,
    speaker: view.speaker ? toSpeakerRow(view.speaker) : undefined,
    directive: view.directive,
    personaLabel: view.personaLabel,
    persona: view.persona ? toPersonaRow(view.persona) : undefined,
    terms: view.terms?.map((t) => ({ source: t.source, dest: t.dest })),
    prompt: view.prompt
  }
}

// 話者 DTO を結果行へ写す。表示用の識別・属性ラベルをそのまま渡す。
function toSpeakerRow(s: api.SpeakerView): SpeakerRow {
  return { edid: s.edid, sex: s.sex, age: s.age, voice: s.voice }
}

// 口調メタの DTO を結果行へ写す。段階（0..2）と決定経路（本文・voice・保留・汎用・PC）は backend が妥当な値を保証する境界のため、union へ確定する。
function toPersonaRow(p: api.PersonaView): PersonaRow {
  return {
    cell: p.cell,
    trait: p.trait,
    attitudeBand: p.attitudeBand as 0 | 1 | 2,
    emotionBand: p.emotionBand as 0 | 1 | 2,
    marked: p.marked,
    decisionPath: p.decisionPath as "本文" | "voice" | "保留" | "汎用" | "PC",
    sex: p.sex
  }
}
