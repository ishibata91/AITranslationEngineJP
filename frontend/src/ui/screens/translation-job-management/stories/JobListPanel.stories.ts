import type { Meta, StoryObj } from "@storybook/svelte-vite"

import JobListPanel from "../JobListPanel.svelte"
import { jobListPanelFixtures } from "../__fixtures__/translation-job-management-fixtures"

const meta = {
  title: "Screens/Translation Job Management/JobListPanel",
  component: JobListPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof JobListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: jobListPanelFixtures.ready
}

export const Loading: Story = {
  args: jobListPanelFixtures.loading
}

export const Empty: Story = {
  args: jobListPanelFixtures.empty
}

export const Error: Story = {
  args: jobListPanelFixtures.error
}
