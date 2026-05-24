import type { Meta, StoryObj } from "@storybook/svelte-vite"

import JobRunTargetSummary from "../JobRunTargetSummary.svelte"
import { jobRunTargetSummaryFixtures } from "../__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Screen Components/Job Run/JobRunTargetSummary",
  component: JobRunTargetSummary,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof JobRunTargetSummary>

export default meta

type Story = StoryObj<typeof meta>

export const TermPhase: Story = {
  args: jobRunTargetSummaryFixtures.termPhase
}

export const PersonaPhase: Story = {
  args: jobRunTargetSummaryFixtures.personaPhase
}

export const LongSourcePath: Story = {
  args: jobRunTargetSummaryFixtures.longSourcePath
}
