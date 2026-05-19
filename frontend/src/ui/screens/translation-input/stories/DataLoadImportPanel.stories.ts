import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DataLoadImportPanel from "../DataLoadImportPanel.svelte"
import { translationInputPanelFixtures } from "../__fixtures__/translation-input-panel-fixtures"

const meta = {
  title: "Screens/Translation Input/DataLoadImportPanel",
  component: DataLoadImportPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DataLoadImportPanel>

export default meta

type Story = StoryObj<typeof meta>

export const NoFile: Story = {
  args: translationInputPanelFixtures.dataLoadImportPanel.noFile
}

export const SelectedFile: Story = {
  args: translationInputPanelFixtures.dataLoadImportPanel.selectedFile
}

export const Importing: Story = {
  args: translationInputPanelFixtures.dataLoadImportPanel.importing
}
