import type { Meta, StoryObj } from "@storybook/svelte-vite"
import MissingStringsWarning from "./MissingStringsWarning.svelte"

// 片側 Strings 欠け時の警告。既存訳（参照訳・固有名の確定訳語）を再利用できない状態を知らせる。
const meta = {
  title: "UI Components/MissingStringsWarning",
  component: MissingStringsWarning
} satisfies Meta<typeof MissingStringsWarning>

export default meta
type Story = StoryObj<typeof meta>

// 日本語 Strings だけ無い。実運用で最も起きやすい欠け方。
export const MissingJapanese: Story = {
  name: "日本語欠け",
  args: { presence: { english: true, japanese: false } }
}

// 英語 Strings だけ無い。
export const MissingEnglish: Story = {
  name: "英語欠け",
  args: { presence: { english: false, japanese: true } }
}

// 両方無い。
export const MissingBoth: Story = {
  name: "両方欠け",
  args: { presence: { english: false, japanese: false } }
}

// 両方ある。警告は出ない（非表示の確認用）。
export const BothPresent: Story = {
  name: "揃い（非表示）",
  args: { presence: { english: true, japanese: true } }
}

// 未判定（presence 未指定）。警告は出ない。
export const Undetermined: Story = {
  name: "未判定（非表示）",
  args: {}
}
