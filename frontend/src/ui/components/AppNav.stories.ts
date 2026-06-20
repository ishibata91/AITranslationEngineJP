import type { Meta, StoryObj } from "@storybook/svelte-vite"
import AppNav from "./AppNav.svelte"
import type { NavItem } from "./app-nav-view"

// 画面間ナビゲーションバー。翻訳実行からテンプレート編集へ飛ぶ導線を確認する。
// 作業中分類（Review/Changed Components）に置く。承認後に UI Components へ戻す。
const meta = {
  title: "Review/Changed Components/ナビゲーション",
  component: AppNav,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof AppNav>

export default meta
type Story = StoryObj<typeof meta>

const ITEMS: NavItem[] = [
  { route: "translation", label: "翻訳実行" },
  { route: "template", label: "プロンプトテンプレート" }
]

const noop = () => {}

// 翻訳実行画面を見ている状態。テンプレート編集への導線が出る。
export const TranslationActive: Story = {
  name: "翻訳実行中",
  args: { active: "translation", items: ITEMS, onNavigate: noop }
}

// テンプレート編集画面を見ている状態。翻訳実行へ戻る導線が出る。
export const TemplateActive: Story = {
  name: "テンプレート編集中",
  args: { active: "template", items: ITEMS, onNavigate: noop }
}
