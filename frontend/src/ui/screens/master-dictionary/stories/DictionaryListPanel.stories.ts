import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryListPanel from "../DictionaryListPanel.svelte"
import { dictionaryListPanelFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screens/Master Dictionary/DictionaryListPanel",
  component: DictionaryListPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Normal: Story = {
  args: dictionaryListPanelFixtures.normal
}

export const Empty: Story = {
  args: dictionaryListPanelFixtures.empty
}

export const FilteredEmpty: Story = {
  args: dictionaryListPanelFixtures.filteredEmpty
}

export const LongText: Story = {
  args: dictionaryListPanelFixtures.longText
}
