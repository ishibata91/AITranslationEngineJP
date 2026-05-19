import type { Meta, StoryObj } from "@storybook/svelte-vite"

import FoundationDataPanel from "../FoundationDataPanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screens/Translation Job Setup/FoundationDataPanel",
  component: FoundationDataPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof FoundationDataPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Populated: Story = {
  args: translationJobSetupPanelFixtures.foundationDataPanel.populated
}

export const Empty: Story = {
  args: translationJobSetupPanelFixtures.foundationDataPanel.empty
}
