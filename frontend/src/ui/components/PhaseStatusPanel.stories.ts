import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseStatusPanel from "./PhaseStatusPanel.svelte"
import { phaseStatusPanelFixtures } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/Phase Status Panel/PhaseStatusPanel",
  component: PhaseStatusPanel,
  args: phaseStatusPanelFixtures.personaAiSettingsReady
} satisfies Meta<typeof PhaseStatusPanel>

export default meta

type Story = StoryObj<typeof meta>

export const PersonaAiSettingsReady: Story = {}

export const TermAiSettingsReady: Story = {
  args: phaseStatusPanelFixtures.termAiSettingsReady
}

export const BodyAiSettingsReady: Story = {
  args: phaseStatusPanelFixtures.bodyAiSettingsReady
}

export const Ready: Story = {
  args: phaseStatusPanelFixtures.ready
}

export const Running: Story = {
  args: phaseStatusPanelFixtures.running
}

export const RecoverableFailed: Story = {
  args: phaseStatusPanelFixtures.recoverableFailed
}

export const Completed: Story = {
  args: phaseStatusPanelFixtures.completed
}
