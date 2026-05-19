import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseStatusPanel from "./PhaseStatusPanel.svelte"
import { phaseStatusPanelFixtures } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/PhaseStatusPanel",
  component: PhaseStatusPanel,
  args: phaseStatusPanelFixtures.ready
} satisfies Meta<typeof PhaseStatusPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {}

export const Running: Story = {
  args: phaseStatusPanelFixtures.running
}

export const RecoverableFailed: Story = {
  args: phaseStatusPanelFixtures.recoverableFailed
}

export const Completed: Story = {
  args: phaseStatusPanelFixtures.completed
}
