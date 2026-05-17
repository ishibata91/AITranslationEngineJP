import type { Meta, StoryObj } from "@storybook/svelte-vite"

import RunStatusPanel from "../RunStatusPanel.svelte"
import { runStatusPanelFixtures } from "../__fixtures__/master-persona-panel-fixtures"

const meta = {
  title: "Screens/Master Persona/RunStatusPanel",
  component: RunStatusPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof RunStatusPanel>

export default meta

type Story = StoryObj<typeof meta>

export const BeforeGeneration: Story = {
  args: runStatusPanelFixtures.beforeGeneration
}

export const Generating: Story = {
  args: runStatusPanelFixtures.generating
}

export const Failed: Story = {
  args: runStatusPanelFixtures.failed
}

export const Completed: Story = {
  args: runStatusPanelFixtures.completed
}
