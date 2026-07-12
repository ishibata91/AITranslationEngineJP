import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TargetPluginsNavPreview from "./TargetPluginsNavPreview.svelte"

// 翻訳対象プラグイン画面の最上部に画面間ナビゲーションが乗った状態。
// ナビに「翻訳対象プラグイン」タブがあり、他画面へ飛べることを画面の文脈で確認する。
// Storybook 人間レビュー承認済み（通常分類 Screens）。
const meta = {
  title: "Screens/翻訳対象プラグイン（ナビ付き）",
  component: TargetPluginsNavPreview,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof TargetPluginsNavPreview>

export default meta
type Story = StoryObj<typeof meta>

// ナビ付きの翻訳対象プラグイン画面。
export const Default: Story = {
  name: "翻訳対象プラグイン＋ナビ"
}
