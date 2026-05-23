import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TranslationCompleteSummaryPanel from "../TranslationCompleteSummaryPanel.svelte"
import { translationCompleteSummaryPanelFixture } from "../__fixtures__/translation-complete-fixture"

const meta = {
  title: "Screen Components/Job Run/TranslationCompleteSummaryPanel",
  component: TranslationCompleteSummaryPanel,
  args: translationCompleteSummaryPanelFixture
} satisfies Meta<typeof TranslationCompleteSummaryPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
