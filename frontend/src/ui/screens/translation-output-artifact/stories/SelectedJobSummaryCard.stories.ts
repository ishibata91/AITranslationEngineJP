import type { Meta, StoryObj } from "@storybook/svelte-vite"

import SelectedJobSummaryCard from "../SelectedJobSummaryCard.svelte"
import { selectedJobSummaryCardFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screen Components/Translation Output Artifact/SelectedJobSummaryCard",
  component: SelectedJobSummaryCard,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof SelectedJobSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Selected: Story = {
  args: selectedJobSummaryCardFixtures.selected
}

export const Rejected: Story = {
  args: selectedJobSummaryCardFixtures.rejected
}

export const AwaitingSelection: Story = {
  args: selectedJobSummaryCardFixtures.awaitingSelection
}
