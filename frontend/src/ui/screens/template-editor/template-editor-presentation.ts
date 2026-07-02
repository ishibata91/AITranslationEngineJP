// テンプレート編集画面の表示用の値（自由記述の口調欄メタと PC 性別の選択肢）。型は template-editor-view.ts。
import type { PcSex, PromptTemplateField } from "./template-editor-view"

// 自由記述の口調指示文の編集 field のメタ。汎用台詞と PC 発話の 2 つを同じ形で出す。
export const TONE_TEXT_FIELDS: ReadonlyArray<{
  field: Extract<PromptTemplateField, "genericToneText" | "pcToneText">
  label: string
  hint: string
}> = [
  {
    field: "genericToneText",
    label: "汎用台詞の口調（話者なし）",
    hint: "衛兵などの汎用台詞は話者を特定できないため、対人の口調をここで自由に書きます。感情の強弱と性別の一人称・語尾はシステムが自動で重ねます。"
  },
  {
    field: "pcToneText",
    label: "PC 発話の口調（選択肢）",
    hint: "プレイヤーの選択肢の台詞へ当てる口調を自由に書きます。感情の強弱と、下で選んだ PC の性別の一人称・語尾を自動で重ねます。"
  }
]

// PC の性別の選択肢。一人称・語尾の根拠。未指定は付けない。
export const PC_SEX_OPTIONS: ReadonlyArray<{ value: PcSex; label: string }> = [
  { value: "", label: "未指定" },
  { value: "Male", label: "男性" },
  { value: "Female", label: "女性" }
]
