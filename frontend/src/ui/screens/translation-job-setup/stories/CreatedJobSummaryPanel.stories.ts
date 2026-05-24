import type { Meta, StoryObj } from "@storybook/svelte-vite"

import CreatedJobSummaryPanel from "../CreatedJobSummaryPanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Job Setup/CreatedJobSummaryPanel",
  component: CreatedJobSummaryPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof CreatedJobSummaryPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Created: Story = {
  args: translationJobSetupPanelFixtures.createdJobSummaryPanel.created
}
