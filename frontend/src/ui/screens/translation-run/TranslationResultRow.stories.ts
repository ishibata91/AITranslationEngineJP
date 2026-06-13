import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationResultRow from "./TranslationResultRow.svelte"

const meta = {
  title: "UI Components/TranslationResultRow",
  component: TranslationResultRow,
  parameters: { layout: "padded" }
} satisfies Meta<typeof TranslationResultRow>

export default meta
type Story = StoryObj<typeof meta>

export const Provisional: Story = {
  name: "仮訳",
  args: {
    row: {
      edid: "DLC1BookSerana",
      source:
        "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
      dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
      statusLabel: "仮訳"
    }
  }
}

export const Untranslated: Story = {
  name: "未訳",
  args: {
    row: {
      edid: "DLC1BookAetherium",
      source:
        "The Aetherium Forge lies beneath Blackreach, its fires long since dimmed by the dwarves who fled.",
      dest: "",
      statusLabel: "未訳"
    }
  }
}
