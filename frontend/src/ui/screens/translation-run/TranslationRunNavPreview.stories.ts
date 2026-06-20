import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationRunNavPreview from "./TranslationRunNavPreview.svelte"

// 翻訳実行画面の最上部に画面間ナビゲーションが乗った状態。
// 翻訳実行からプロンプトテンプレート編集へ飛ぶ導線を、画面の文脈で確認する。
// 作業中分類（Review/Changed Screens）に置く。承認後に Screens へ戻す。
const meta = {
  title: "Review/Changed Screens/翻訳実行（ナビ付き）",
  component: TranslationRunNavPreview,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof TranslationRunNavPreview>

export default meta
type Story = StoryObj<typeof meta>

// 翻訳完了後の翻訳実行画面に、テンプレート編集への導線が出る。
export const Default: Story = {
  name: "翻訳実行＋ナビ"
}
