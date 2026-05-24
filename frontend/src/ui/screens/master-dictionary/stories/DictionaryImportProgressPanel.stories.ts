import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryImportProgressPanel from "../DictionaryImportProgressPanel.svelte"
import { dictionaryImportProgressPanelFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screen Components/Master Dictionary/DictionaryImportProgressPanel",
  component: DictionaryImportProgressPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryImportProgressPanel>

export default meta

type Story = StoryObj<typeof meta>

export const NoFileSelected: Story = {
  args: dictionaryImportProgressPanelFixtures.noFileSelected
}

export const Running: Story = {
  args: dictionaryImportProgressPanelFixtures.running
}

export const Completed: Story = {
  args: dictionaryImportProgressPanelFixtures.completed
}
