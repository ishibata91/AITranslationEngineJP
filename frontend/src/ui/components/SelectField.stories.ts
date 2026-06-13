import type { Meta, StoryObj } from "@storybook/svelte-vite"
import SelectField from "./SelectField.svelte"

const MODELS = ["gpt-4o-mini", "gpt-4o", "gpt-4.1-mini", "o4-mini"]

const meta = {
  title: "UI Components/SelectField",
  component: SelectField,
  parameters: { layout: "padded" },
  args: { onChange: () => {} }
} satisfies Meta<typeof SelectField>

export default meta
type Story = StoryObj<typeof meta>

// 選択肢が空のとき。emptyText を表示して無効化する。
export const Empty: Story = {
  name: "選択肢なし",
  args: {
    label: "モデル",
    value: "",
    options: [],
    emptyText: "未取得",
    hint: "エンドポイントと API キーを入れてから取得します。"
  }
}

// 選択肢があり、未選択。placeholder を表示する。
export const Unselected: Story = {
  name: "未選択",
  args: {
    label: "モデル",
    value: "",
    options: MODELS,
    placeholder: "モデルを選択"
  }
}

// 選択済み。
export const Selected: Story = {
  name: "選択済み",
  args: {
    label: "モデル",
    value: "gpt-4o-mini",
    options: MODELS
  }
}
