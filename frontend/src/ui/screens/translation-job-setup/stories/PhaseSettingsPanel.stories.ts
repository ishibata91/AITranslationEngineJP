import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PhaseSettingsPanel from "../PhaseSettingsPanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Job Setup/PhaseSettingsPanel",
  component: PhaseSettingsPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PhaseSettingsPanel>

export default meta

type Story = StoryObj<typeof meta>

export const PhaseCards: Story = {
  args: translationJobSetupPanelFixtures.phaseSettingsPanel.phaseCards
}

export const LegacyRuntime: Story = {
  args: translationJobSetupPanelFixtures.phaseSettingsPanel.legacyRuntime
}
