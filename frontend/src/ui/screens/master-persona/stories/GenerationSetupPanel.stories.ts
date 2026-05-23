import type { Meta, StoryObj } from "@storybook/svelte-vite"

import GenerationSetupPanel from "../GenerationSetupPanel.svelte"
import { generationSetupPanelFixtures } from "../__fixtures__/master-persona-panel-fixtures"

const meta = {
  title: "Screen Components/Master Persona/GenerationSetupPanel",
  component: GenerationSetupPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof GenerationSetupPanel>

export default meta

type Story = StoryObj<typeof meta>

export const NoJsonSelected: Story = {
  args: generationSetupPanelFixtures.noJsonSelected
}

export const JsonSelected: Story = {
  args: generationSetupPanelFixtures.jsonSelected
}

export const PreviewSucceeded: Story = {
  args: generationSetupPanelFixtures.previewSucceeded
}

export const MissingAISettings: Story = {
  args: generationSetupPanelFixtures.missingAISettings
}

export const ModelListRefreshing: Story = {
  args: generationSetupPanelFixtures.modelListRefreshing
}

export const LongModelName: Story = {
  args: generationSetupPanelFixtures.longModelName
}
