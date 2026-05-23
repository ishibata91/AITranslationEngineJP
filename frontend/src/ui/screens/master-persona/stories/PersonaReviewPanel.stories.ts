import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PersonaReviewPanel from "../PersonaReviewPanel.svelte"
import { personaReviewPanelFixtures } from "../__fixtures__/master-persona-panel-fixtures"

const meta = {
  title: "Screen Components/Master Persona/PersonaReviewPanel",
  component: PersonaReviewPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PersonaReviewPanel>

export default meta

type Story = StoryObj<typeof meta>

export const EmptyList: Story = {
  args: personaReviewPanelFixtures.emptyList
}

export const WithList: Story = {
  args: personaReviewPanelFixtures.withList
}

export const RowSelected: Story = {
  args: personaReviewPanelFixtures.rowSelected
}

export const FilteredEmpty: Story = {
  args: personaReviewPanelFixtures.filteredEmpty
}

export const StaleSelectionCleared: Story = {
  args: personaReviewPanelFixtures.staleSelectionCleared
}
