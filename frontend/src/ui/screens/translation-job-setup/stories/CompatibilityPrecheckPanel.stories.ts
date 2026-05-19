import type { Meta, StoryObj } from "@storybook/svelte-vite"

import CompatibilityPrecheckPanel from "../CompatibilityPrecheckPanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screens/Translation Job Setup/CompatibilityPrecheckPanel",
  component: CompatibilityPrecheckPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof CompatibilityPrecheckPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: translationJobSetupPanelFixtures.compatibilityPrecheckPanel.ready
}

export const Dirty: Story = {
  args: translationJobSetupPanelFixtures.compatibilityPrecheckPanel.dirty
}
