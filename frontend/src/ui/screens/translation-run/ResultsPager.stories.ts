import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ResultsPager from "./ResultsPager.svelte"

// Storybook 人間レビュー承認済み。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/ResultsPager",
  component: ResultsPager,
  parameters: { layout: "padded" },
  args: { onPrev: () => {}, onNext: () => {} }
} satisfies Meta<typeof ResultsPager>

export default meta
type Story = StoryObj<typeof meta>

// 先頭ページ。前へは無効、次へは有効。
export const FirstPage: Story = {
  name: "先頭ページ（前へ無効）",
  args: { pageNumber: 1, canPrev: false, canNext: true }
}

// 中間ページ。前へ・次へとも有効。
export const MiddlePage: Story = {
  name: "中間ページ（両方有効）",
  args: { pageNumber: 2, canPrev: true, canNext: true }
}

// 末尾ページ。次へは無効、前へは有効。
export const LastPage: Story = {
  name: "末尾ページ（次へ無効）",
  args: { pageNumber: 3, canPrev: true, canNext: false }
}

// 単一ページ。前へ・次へとも無効。
export const SinglePage: Story = {
  name: "単一ページ（両方無効）",
  args: { pageNumber: 1, canPrev: false, canNext: false }
}
