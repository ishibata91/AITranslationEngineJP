// 翻訳実行画面の表示用の値（入力欄メタ・段階表示・訳状態トーン）。型は translation-run-view.ts。
import type { BadgeTone } from "@ui/components/status-badge"
import type {
  AttitudeBand,
  BatchProgressView,
  BatchStage,
  DecisionPath,
  EmotionBand,
  ExtractStep,
  FormFieldDescriptor,
  PersonaMeta,
  ProgressStage,
  RunPhase,
  TranslationProvider
} from "./translation-run-view"

// plugin はファイル選択、model は取得した一覧からの選択で扱うため、テキスト入力欄は接続情報の 2 つだけ。
export const PROVIDER_FIELDS: ReadonlyArray<FormFieldDescriptor> = [
  {
    field: "endpoint",
    label: "エンドポイント",
    hint: "OpenAI 互換 API の接続先 URL を入れます。",
    placeholder: "https://api.openai.com/v1",
    secret: false
  },
  {
    field: "apiKey",
    label: "API キー",
    hint: "この画面では保存せず、実行時だけ使います。",
    placeholder: "sk-...",
    secret: true
  }
]

// OpenAI（batch）選択時の接続情報欄。公式 endpoint と OpenAI API キーを使うことを明示する。
const OPENAI_PROVIDER_FIELDS: ReadonlyArray<FormFieldDescriptor> = [
  {
    field: "endpoint",
    label: "エンドポイント",
    hint: "OpenAI の公式 API の接続先 URL を使います。",
    placeholder: "https://api.openai.com/v1",
    secret: false
  },
  {
    field: "apiKey",
    label: "OpenAI API キー",
    hint: "この画面では保存せず、送信・状態確認・取り込みのたびに使います。",
    placeholder: "sk-...",
    secret: true
  }
]

// xAI（batch）選択時の接続情報欄。エンドポイントは xAI 用に読み替える。API キー欄は同期と同じ。
const XAI_PROVIDER_FIELDS: ReadonlyArray<FormFieldDescriptor> = [
  {
    field: "endpoint",
    label: "エンドポイント",
    hint: "xAI の接続先 URL。空なら既定の https://api.x.ai を使います。",
    placeholder: "https://api.x.ai",
    secret: false
  },
  {
    field: "apiKey",
    label: "API キー",
    hint: "この画面では保存せず、送信・反映のたびに使います。",
    placeholder: "xai-...",
    secret: true
  }
]

// 配送方式ごとの接続情報欄。batch は提供元に対応する接続情報欄を返す。
export function providerFields(
  provider: TranslationProvider
): ReadonlyArray<FormFieldDescriptor> {
  if (provider === "openai") return OPENAI_PROVIDER_FIELDS
  if (provider === "xai") return XAI_PROVIDER_FIELDS
  return PROVIDER_FIELDS
}

// 配送方式の選択肢と、その表示ラベル。翻訳実行画面の先頭で切り替える。
export const PROVIDER_OPTIONS: ReadonlyArray<TranslationProvider> = [
  "sync",
  "openai",
  "xai"
]

export const PROVIDER_LABEL: Record<TranslationProvider, string> = {
  sync: "OpenAI 互換 API",
  openai: "OpenAI（batch）",
  xai: "xAI（batch）"
}

// モデル取得ボタンの補足。配送方式で取得先が変わるため、文言を出し分ける。
export const MODEL_HINT: Record<TranslationProvider, string> = {
  sync: "エンドポイントと API キーを入れてから取得します。",
  openai:
    "OpenAI の batch 対応モデルを取得します。API キーを入れてから取得します。",
  xai: "xAI の batch 対応モデルを取得します。API キーを入れてから取得します。"
}

// batch の送信直後に出す案内。反映で結果を取りにいく運用を伝える。
export const SUBMIT_NOTICE =
  "batch を送信しました。しばらく後に「反映」で結果を取得します（最大約 24 時間）。"

// batch 操作の補足。状態確認で進行状況を最新化し、完了段があれば主アクションで取り込んで次へ進む。
export const BATCH_ACTION_HINT =
  "「状態確認」で進行状況を最新化します。完了した段があれば、右のボタンで取り込んで次へ進みます。"

// batch の進行段の表示ラベル。進行状況ステッパーの各段に使う。
export const BATCH_STAGE_LABEL: Record<BatchStage, string> = {
  proper: "固有名",
  body: "本文",
  done: "完了"
}

// 進行状況ステッパーの段の並び。固有名 → 本文 → 完了 の 2 段構成（＋完了）を常に見せる。
export const BATCH_STAGE_STEPS: ReadonlyArray<BatchStage> = [
  "proper",
  "body",
  "done"
]

// ステッパー 1 段の表示値。cls は daisyUI steps の色クラス、content は marker の上書き（完了段は ✓）。
export interface BatchStepView {
  stage: BatchStage
  label: string
  cls: string
  content?: string
}

// 進行状況からステッパー各段の表示状態を組む。
// 未確認（progress 無し）は色を付けず中立。全完了は全段を success。
// 進行中は、過去段を success（✓）、現在段を primary、以降を中立にする。
export function batchStepViews(progress?: BatchProgressView): BatchStepView[] {
  const current = progress ? BATCH_STAGE_STEPS.indexOf(progress.stage) : -1
  return BATCH_STAGE_STEPS.map((stage, i) => {
    const label = BATCH_STAGE_LABEL[stage]
    if (!progress) return { stage, label, cls: "" }
    if (progress.stage === "done")
      return { stage, label, cls: "step-success", content: "✓" }
    if (i < current) return { stage, label, cls: "step-success", content: "✓" }
    if (i === current) return { stage, label, cls: "step-primary" }
    return { stage, label, cls: "" }
  })
}

// 主アクションの表示文言。送信（新規 / 再送信）と、固有名・本文それぞれの取り込みを分ける。
export const BATCH_SEND_LABEL = "送信して開始"
export const BATCH_RETRY_UNTRANSLATED_LABEL = "未訳だけを再送信"
export const BATCH_APPLY_PROPER_LABEL = "取り込んで本文を送信"
export const BATCH_APPLY_BODY_LABEL = "取り込んで完了"

// 主アクションの種別。send=送信（onSubmit）、apply=取り込み（onApply）。進行状況で排他に切り替わる。
export type BatchMainActionKind = "send" | "apply"

// 主アクションの表示値。kind で押下時の動作、label で文言、enabled で活性（処理中・busy は呼び出し側で重ねる）。
export interface BatchMainAction {
  kind: BatchMainActionKind
  label: string
  enabled: boolean
}

// 進行状況から主アクションの種別・ラベル・活性を導く純関数。
// 未確認・全完了は「送信して開始」（新規 / 未訳の残りを再送信）。
// 固有名段が完了なら「取り込んで本文を送信」、本文段が完了なら「取り込んで完了」。処理待ちが残る間は非活性。
export function batchMainAction(
  progress: BatchProgressView | undefined,
  canSubmit: boolean
): BatchMainAction {
  if (progress?.stage === "done" && progress.untranslatedCount > 0) {
    return {
      kind: "send",
      label: BATCH_RETRY_UNTRANSLATED_LABEL,
      enabled: canSubmit
    }
  }
  if (!progress || progress.stage === "done") {
    return { kind: "send", label: BATCH_SEND_LABEL, enabled: canSubmit }
  }
  if (progress.stage === "proper") {
    return {
      kind: "apply",
      label: BATCH_APPLY_PROPER_LABEL,
      enabled: progress.canApply
    }
  }
  return {
    kind: "apply",
    label: BATCH_APPLY_BODY_LABEL,
    enabled: progress.canApply
  }
}

// 現段 batch の件数ラベル。進行状況パネルで内訳を出す。
export const BATCH_COUNT_LABEL = {
  total: "総数",
  pending: "処理待ち",
  succeeded: "成功",
  failed: "失敗"
} as const

// 進行状況パネル下部の補足。状態未確認・処理中・完了段あり・全完了で出し分ける。
export const BATCH_UNCHECKED_HINT = "「状態確認」で進行状況を取得します。"
export const BATCH_WAITING_HINT =
  "処理待ちが残っています。完了までお待ちください。"
export const BATCH_APPLYABLE_HINT =
  "現段が完了しました。「取り込み」で結果を取り込めます。"
export const BATCH_DONE_HINT = "すべての翻訳が完了しました。"

// 取り込みの結果として出す案内（Container が完了段に応じて選ぶ）。
export const APPLIED_PROPER_NOTICE =
  "固有名を取り込み、本文の翻訳を送信しました。"
export const APPLIED_BODY_NOTICE = "本文を取り込みました。翻訳が完了しました。"
export const APPLY_NOTHING_NOTICE = "取り込める完了段はまだありません。"

// batch の取り込み後に未訳が残った場合の案内。主操作と同じ「未訳だけを再送信」を使う。
export function batchUntranslatedNotice(untranslatedCount: number): string {
  if (untranslatedCount <= 0) return ""
  return `${untranslatedCount} 件が未訳のまま残りました。未訳だけを再送信できます。`
}

// 実行の完了時に、未訳のまま残った件数を伝える案内文を組む。
// 0 件なら空文字を返し、案内を出さない（残りが無い完了は状態表示だけで足りる）。
// 1 件以上なら残った件数と、再実行でその件数を訳し直せることを伝える。
// 未訳が残るのは、応答が空だった行やサーバの一時失敗で飛ばした行があるため。内訳は実行ログが持つ。
export function untranslatedNotice(untranslatedCount: number): string {
  if (untranslatedCount <= 0) return ""
  return `${untranslatedCount} 件が未訳のまま残りました。もう一度実行すると、残った ${untranslatedCount} 件だけを翻訳します。`
}

// 段階ごとの表示ラベルと、StatusBadge の意味トーン。
export const PHASE_PRESENTATION: Record<
  RunPhase,
  { label: string; tone: BadgeTone }
> = {
  idle: { label: "未実行", tone: "ghost" },
  running: { label: "実行中", tone: "primary" },
  done: { label: "完了", tone: "success" },
  error: { label: "失敗", tone: "danger" }
}

// 実行中の段階ごとの表示ラベル。進捗バーの見出しに使う。
export const PROGRESS_STAGE_LABEL: Record<ProgressStage, string> = {
  extract: "台詞を抽出しています",
  translate: "本文を翻訳しています"
}

// extract 段の内訳サブ段ごとの表示ラベル。無音区間で今どのサブ段かを見出しに出す。
// PROGRESS_STAGE_LABEL と同じ「〜しています」体にそろえる。
export const EXTRACT_STEP_LABEL: Record<ExtractStep, string> = {
  extract: "台詞を抽出しています",
  derive: "辞書を準備しています",
  reference: "既存訳を取り込んでいます",
  ingest: "翻訳対象を仕分けています"
}

// 訳状態ラベルごとの意味トーン。未知ラベルは控えめな outline にする。
export function statusTone(statusLabel: string): BadgeTone {
  switch (statusLabel) {
    case "承認":
      return "success"
    case "訳済":
      return "primary"
    case "仮訳":
      return "accent"
    default:
      return "outline"
  }
}

// 対人段階の表示ラベル。口調メタの根拠行の組み立て（personaMetaParts）でだけ使う。
const ATTITUDE_LABEL: Record<AttitudeBand, string> = {
  0: "尊大",
  1: "中立",
  2: "丁寧"
}

// 感情段階の表示ラベル。口調メタの根拠行の組み立て（personaMetaParts）でだけ使う。
const EMOTION_LABEL: Record<EmotionBand, string> = {
  0: "抑制",
  1: "中",
  2: "激情"
}

// 決定経路の表示ラベル。voice は「声質」と読み替えて出す。
// 汎用 / PC は話者を解決できない台詞で、利用者の自由記述の口調指示文を当てた経路。
export const DECISION_PATH_LABEL: Record<DecisionPath, string> = {
  本文: "本文",
  voice: "声質",
  保留: "保留",
  汎用: "汎用台詞",
  PC: "PC発話"
}

// 決定経路の補足説明。根拠を title 属性で読めるようにする。
export const DECISION_PATH_HINT: Record<DecisionPath, string> = {
  本文: "印（対人マーカーを含む台詞数）が十分で、本文 2 軸から対人段階を決めた。",
  voice: "印が不足のため、声質の気質を prior として対人段階を決めた。",
  保留: "印が不足で固有声質に気質も無いため、薄い本文値を低信頼で保持した。",
  汎用: "話者を特定できない汎用台詞のため、利用者が書いた汎用の口調指示文を当て、感情段階を本文 1 行から、性別を INFO の条件から重ねた。",
  PC: "プレイヤーの選択肢の台詞のため、利用者が書いた PC の口調指示文を当て、感情段階を本文 1 行から、性別を利用者選択の PC 性別から重ねた。"
}

// personaMetaParts は口調メタの根拠行を区切り表示用の文字列列にする。
// 名指し話者は 決定経路・対人・感情・(性別)・印 を出す。汎用・PC は 感情・(性別) だけを出し、
// 決定経路（汎用台詞 / PC発話）は呼び出し側が見出しに出すため重複させない。
export function personaMetaParts(p: PersonaMeta): string[] {
  if (p.cell) {
    return [
      DECISION_PATH_LABEL[p.decisionPath],
      `対人 ${ATTITUDE_LABEL[p.attitudeBand ?? 1]}`,
      `感情 ${EMOTION_LABEL[p.emotionBand]}`,
      p.sex ? `性別 ${p.sex}` : "",
      `印 ${p.marked ?? 0}`
    ].filter(Boolean)
  }
  return [
    `感情 ${EMOTION_LABEL[p.emotionBand]}`,
    p.sex ? `性別 ${p.sex}` : ""
  ].filter(Boolean)
}
