import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import {
  EMPTY_STATE,
  MODELS_LOADING_STATE,
  READY_STATE,
  RUNNING_EXTRACT_STATE,
  RUNNING_TRANSLATE_STATE,
  DONE_STATE,
  DONE_PERSONA_STATE,
  ERROR_STATE
} from "./translation-run.fixtures"

// Storybook 人間レビュー承認済み。通常分類（Screens）に置く。
const meta = {
  title: "Screens/翻訳実行",
  component: TranslationRunScreen,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    onFieldInput: () => {},
    onSelectPlugin: () => {},
    onLoadModels: () => {},
    onRun: () => {}
  }
} satisfies Meta<typeof TranslationRunScreen>

export default meta
type Story = StoryObj<typeof meta>

// 初期表示。未入力・モデル未取得で実行できない。
export const Empty: Story = {
  name: "空状態",
  args: { ...EMPTY_STATE }
}

// getModels でモデル一覧を取得中。
export const ModelsLoading: Story = {
  name: "モデル取得中",
  args: { ...MODELS_LOADING_STATE }
}

// 入力が揃い、実行できる。
export const Ready: Story = {
  name: "入力済み",
  args: { ...READY_STATE }
}

// 実行中（台詞抽出）。不定の進捗バーが出る。
export const RunningExtract: Story = {
  name: "実行中（抽出）",
  args: { ...RUNNING_EXTRACT_STATE }
}

// 実行中（本文翻訳）。done/total の確定進捗バーが出る。
export const RunningTranslate: Story = {
  name: "実行中（本文翻訳）",
  args: { ...RUNNING_TRANSLATE_STATE }
}

// 完了（叙述文）。原文と訳文が並ぶ。
export const Done: Story = {
  name: "完了",
  args: { ...DONE_STATE }
}

// 完了（台詞）。訳文と口調指示が並び、話者ごとの口調差が観測できる。
export const DonePersona: Story = {
  name: "完了（口調差）",
  args: { ...DONE_PERSONA_STATE }
}

// 実行に失敗した。
export const Errored: Story = {
  name: "エラー",
  args: { ...ERROR_STATE }
}
