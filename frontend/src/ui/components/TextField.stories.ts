import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TextField from "./TextField.svelte"

const meta = {
  title: "UI Components/TextField",
  component: TextField,
  parameters: { layout: "padded" },
  args: { oninput: () => {} }
} satisfies Meta<typeof TextField>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  name: "未入力",
  args: {
    label: "エンドポイント",
    value: "",
    placeholder: "https://api.openai.com/v1",
    hint: "OpenAI 互換 API の接続先 URL を入れます。"
  }
}

export const Filled: Story = {
  name: "入力済み",
  args: {
    label: "エンドポイント",
    value: "https://api.openai.com/v1",
    hint: "OpenAI 互換 API の接続先 URL を入れます。"
  }
}

export const Secret: Story = {
  name: "機密値",
  args: {
    label: "API キー",
    value: "sk-demo-key",
    secret: true,
    hint: "この画面では保存せず、実行時だけ使います。"
  }
}
