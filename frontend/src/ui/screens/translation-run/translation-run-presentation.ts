// 翻訳実行画面の表示用の値（入力欄メタ・段階表示・訳状態トーン）。型は translation-run-view.ts。
import type { BadgeTone } from "@ui/components/status-badge"
import type {
  AttitudeBand,
  DecisionPath,
  EmotionBand,
  FormFieldDescriptor,
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

// 対人段階の表示ラベル。
export const ATTITUDE_LABEL: Record<AttitudeBand, string> = {
  0: "尊大",
  1: "中立",
  2: "丁寧"
}

// 感情段階の表示ラベル。
export const EMOTION_LABEL: Record<EmotionBand, string> = {
  0: "抑制",
  1: "中",
  2: "激情"
}

// 決定経路の表示ラベル。voice は「声質」と読み替えて出す。
export const DECISION_PATH_LABEL: Record<DecisionPath, string> = {
  本文: "本文",
  voice: "声質",
  保留: "保留"
}

// 決定経路の補足説明。根拠を title 属性で読めるようにする。
export const DECISION_PATH_HINT: Record<DecisionPath, string> = {
  本文: "印（対人マーカーを含む台詞数）が十分で、本文 2 軸から対人段階を決めた。",
  voice: "印が不足のため、声質の気質を prior として対人段階を決めた。",
  保留: "印が不足で固有声質に気質も無いため、薄い本文値を低信頼で保持した。"
}
