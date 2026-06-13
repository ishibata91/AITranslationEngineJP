import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ResultsPanel from "./ResultsPanel.svelte"
import type { NarrationResultRow } from "./translation-run-view"

const ROWS: NarrationResultRow[] = [
  {
    edid: "DLC1BookSerana",
    source:
      "I have walked these halls for centuries, and still the cold of Castle Volkihar finds me.",
    dest: "私は何世紀もこの広間を歩いてきたが、それでもヴォルキハル城の冷気は私を捉えて離さない。",
    statusLabel: "仮訳"
  },
  {
    edid: "DLC1BookAetherium",
    source:
      "The Aetherium Forge lies beneath Blackreach, its fires long since dimmed by the dwarves who fled.",
    dest: "",
    statusLabel: "未訳"
  }
]

const meta = {
  title: "UI Components/ResultsPanel",
  component: ResultsPanel,
  parameters: { layout: "padded" }
} satisfies Meta<typeof ResultsPanel>

export default meta
type Story = StoryObj<typeof meta>

export const EmptyIdle: Story = {
  name: "空（未実行）",
  args: { phase: "idle", results: [] }
}

export const EmptyRunning: Story = {
  name: "空（実行中）",
  args: { phase: "running", results: [] }
}

export const WithResults: Story = {
  name: "結果あり",
  args: { phase: "done", results: ROWS }
}
