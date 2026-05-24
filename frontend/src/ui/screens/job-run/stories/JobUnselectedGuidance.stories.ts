import type { Meta, StoryObj } from "@storybook/svelte-vite"

import JobUnselectedGuidance from "../JobUnselectedGuidance.svelte"
import { jobUnselectedGuidanceFixtures } from "../__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Screen Components/Job Run/JobUnselectedGuidance",
  component: JobUnselectedGuidance,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof JobUnselectedGuidance>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: jobUnselectedGuidanceFixtures.default
}
