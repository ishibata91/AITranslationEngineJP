// 翻訳実行画面の表示用の値（入力欄メタ・段階表示・訳状態トーン）。型は translation-run-view.ts。
import type { BadgeTone } from "@ui/components/status-badge"
import type { FormFieldDescriptor, RunPhase } from "./translation-run-view"

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
