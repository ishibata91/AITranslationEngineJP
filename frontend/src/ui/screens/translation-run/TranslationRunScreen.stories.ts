import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationRunScreen from "./TranslationRunScreen.svelte"
import {
  EMPTY_STATE,
  MODELS_LOADING_STATE,
  READY_STATE,
  RUNNING_STATE,
  DONE_STATE,
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

// 抽出と翻訳を実行中。
export const Running: Story = {
  name: "実行中",
  args: { ...RUNNING_STATE }
}

// 完了し、原文と訳文が並ぶ。
export const Done: Story = {
  name: "完了",
  args: { ...DONE_STATE }
}

// 実行に失敗した。
export const Errored: Story = {
  name: "エラー",
  args: { ...ERROR_STATE }
}
