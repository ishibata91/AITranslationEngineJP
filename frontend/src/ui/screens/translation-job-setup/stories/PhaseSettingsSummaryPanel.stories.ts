import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PhaseSettingsSummaryPanel from "../PhaseSettingsSummaryPanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Job Setup/PhaseSettingsSummaryPanel",
  component: PhaseSettingsSummaryPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PhaseSettingsSummaryPanel>

export default meta

type Story = StoryObj<typeof meta>

export const PhaseCards: Story = {
  args: translationJobSetupPanelFixtures.phaseSettingsSummaryPanel.phaseCards
}

export const Legacy: Story = {
  args: translationJobSetupPanelFixtures.phaseSettingsSummaryPanel.legacy
}
