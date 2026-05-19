import type { Meta, StoryObj } from "@storybook/svelte-vite"

import JobCard from "../JobCard.svelte"
import { jobCardFixtures } from "../__fixtures__/translation-job-management-fixtures"

const meta = {
  title: "Screens/Translation Job Management/JobCard",
  component: JobCard,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof JobCard>

export default meta

type Story = StoryObj<typeof meta>

export const SelectedRunning: Story = {
  args: jobCardFixtures.selectedRunning
}

export const PausedWithDisabledReason: Story = {
  args: jobCardFixtures.pausedWithDisabledReason
}

export const RecoverableFailed: Story = {
  args: jobCardFixtures.recoverableFailed
}
