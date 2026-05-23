import type { Meta, StoryObj } from "@storybook/svelte-vite"

import JobSetupPurposeHeader from "../JobSetupPurposeHeader.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Job Setup/JobSetupPurposeHeader",
  component: JobSetupPurposeHeader,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof JobSetupPurposeHeader>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: translationJobSetupPanelFixtures.purposeHeader.ready
}

export const Failed: Story = {
  args: translationJobSetupPanelFixtures.purposeHeader.failed
}
