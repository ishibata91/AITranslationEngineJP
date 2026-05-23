import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryImportPanel from "../DictionaryImportPanel.svelte"
import { dictionaryImportPanelFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screen Components/Master Dictionary/DictionaryImportPanel",
  component: DictionaryImportPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryImportPanel>

export default meta

type Story = StoryObj<typeof meta>

export const NoFileSelected: Story = {
  args: dictionaryImportPanelFixtures.noFileSelected
}

export const Running: Story = {
  args: dictionaryImportPanelFixtures.running
}
