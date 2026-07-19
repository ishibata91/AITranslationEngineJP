import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationProgress from "./TranslationProgress.svelte"

// 進捗表示の変更中。人間レビュー中は作業中分類（Review/Changed Components）に置く。
// 承認後に通常分類（UI Components）へ戻す。
const meta = {
  title: "UI Components/TranslationProgress",
  component: TranslationProgress,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationProgress>

export default meta
type Story = StoryObj<typeof meta>

// extract 段のサブ段。C# 抽出子で台詞を抽出中。件数が出ないため不定バー。
export const ExtractingLines: Story = {
  name: "抽出中・台詞抽出",
  args: { stage: "extract", done: 0, total: 0, step: "extract" }
}

// extract 段のサブ段。固有名の部分形を派生中。
export const ExtractingDerive: Story = {
  name: "抽出中・固有名派生",
  args: { stage: "extract", done: 0, total: 0, step: "derive" }
}

// extract 段のサブ段。xTranslator の既存訳を取り込み中。
export const ExtractingReference: Story = {
  name: "抽出中・既存訳取込",
  args: { stage: "extract", done: 0, total: 0, step: "reference" }
}

// extract 段のサブ段。抽出結果を箱へ仕分ける取込段。
export const ExtractingIngest: Story = {
  name: "抽出中・取込段",
  args: { stage: "extract", done: 0, total: 0, step: "ingest" }
}

// サブ段が未指定の extract 段。従来の段ラベルにフォールバックする。
export const ExtractingFallback: Story = {
  name: "抽出中・サブ段未指定",
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
