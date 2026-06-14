import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationProgress from "./TranslationProgress.svelte"

// Storybook 人間レビュー承認済み。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/TranslationProgress",
  component: TranslationProgress,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationProgress>

export default meta
type Story = StoryObj<typeof meta>

// 台詞を抽出している段階。件数が出ないため不定バー。
export const Extracting: Story = {
  name: "抽出中（不定）",
  args: { stage: "extract", done: 0, total: 0 }
}

// 本文翻訳の途中。
export const TranslatingMid: Story = {
  name: "翻訳中（途中）",
  args: { stage: "translate", done: 64, total: 121 }
}

// 本文翻訳がほぼ完了。
export const TranslatingNearDone: Story = {
  name: "翻訳中（ほぼ完了）",
  args: { stage: "translate", done: 118, total: 121 }
}
