import type { Meta, StoryObj } from "@storybook/svelte-vite"

import OutputActionPanel from "../OutputActionPanel.svelte"
import { outputActionPanelFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screen Components/Translation Output Artifact/OutputActionPanel",
  component: OutputActionPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof OutputActionPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: outputActionPanelFixtures.ready
}

export const Disabled: Story = {
  args: outputActionPanelFixtures.disabled
}

export const LongPath: Story = {
  args: outputActionPanelFixtures.longPath
}
