import type { Meta, StoryObj } from "@storybook/svelte-vite"
import FileSelectField from "./FileSelectField.svelte"

const meta = {
  title: "UI Components/FileSelectField",
  component: FileSelectField,
  parameters: { layout: "padded" },
  args: { onSelect: () => {} }
} satisfies Meta<typeof FileSelectField>

export default meta
type Story = StoryObj<typeof meta>

export const Unselected: Story = {
  name: "未選択",
  args: {
    label: "plugin",
    buttonLabel: "plugin を選択",
    path: "",
    emptyText: "plugin ファイルが未選択です。",
    hint: "master と Strings がある Data フォルダ内の plugin を選びます。"
  }
}

export const Selected: Story = {
  name: "選択済み",
  args: {
    label: "plugin",
    buttonLabel: "plugin を選択",
    path: "/Users/me/Skyrim/Data/Dawnguard.esm",
    hint: "master と Strings がある Data フォルダ内の plugin を選びます。"
  }
}
