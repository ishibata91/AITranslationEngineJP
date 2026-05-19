import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryDetailPanel from "../DictionaryDetailPanel.svelte"
import { dictionaryDetailPanelFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screens/Master Dictionary/DictionaryDetailPanel",
  component: DictionaryDetailPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryDetailPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Selected: Story = {
  args: dictionaryDetailPanelFixtures.selected
}

export const Unselected: Story = {
  args: dictionaryDetailPanelFixtures.unselected
}

export const LongText: Story = {
  args: dictionaryDetailPanelFixtures.longText
}
