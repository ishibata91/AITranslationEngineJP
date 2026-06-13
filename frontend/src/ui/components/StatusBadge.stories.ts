import type { Meta, StoryObj } from "@storybook/svelte-vite"
import StatusBadge from "./StatusBadge.svelte"

// Storybook 人間レビュー承認済み。通常分類（UI Components）に置く。
const meta = {
  title: "UI Components/StatusBadge",
  component: StatusBadge,
  parameters: { layout: "centered" }
} satisfies Meta<typeof StatusBadge>

export default meta
type Story = StoryObj<typeof meta>

export const Idle: Story = {
  name: "未実行",
  args: { label: "未実行", tone: "ghost" }
}

export const Running: Story = {
  name: "実行中",
  args: { label: "実行中", tone: "primary", loading: true }
}

export const Done: Story = {
  name: "完了",
  args: { label: "完了", tone: "success" }
}

export const Failed: Story = {
  name: "失敗",
  args: { label: "失敗", tone: "danger" }
}

export const Provisional: Story = {
  name: "仮訳",
  args: { label: "仮訳", tone: "accent" }
}

export const Untranslated: Story = {
  name: "未訳",
  args: { label: "未訳", tone: "outline" }
}
