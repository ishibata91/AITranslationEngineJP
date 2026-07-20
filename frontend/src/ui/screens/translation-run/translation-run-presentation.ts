// 翻訳実行画面の表示用の値（入力欄メタ・段階表示・訳状態トーン）。型は translation-run-view.ts。
import type { BadgeTone } from "@ui/components/status-badge"
import type {
  AttitudeBand,
  DecisionPath,
  EmotionBand,
  ExtractStep,
  FormFieldDescriptor,
  PersonaMeta,
  ProgressStage,
  RunPhase
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

// トークン予算入力欄のメタ。接続情報ではなく実行パラメータのため、PROVIDER_FIELDS とは別に持つ。
// 単位は k（千トークン）。1 リクエストにまとめる原文の最大トークン数の目安を千トークン単位で受ける。
// 長い台詞は少なく、短い台詞は多くまとまる。弱いローカルモデルは小さく、クラウドモデルは大きく指定する想定。
// 表示は文字列で受ける。k からトークン数への換算は配線層が行う。
export const TOKEN_BUDGET_FIELD: FormFieldDescriptor = {
  field: "tokenBudget",
  label: "最大トークン",
  hint: "1 リクエストにまとめる原文の最大トークン数の目安を千トークン（k）単位で入れます。長い台詞は少なく、短い台詞は多くまとまります。",
  placeholder: "2",
  secret: false,
  suffix: "k"
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
