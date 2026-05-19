import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PhaseNavigationFooter from "../PhaseNavigationFooter.svelte"
import { phaseNavigationFooterFixtures } from "../__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Screens/Job Run/PhaseNavigationFooter",
  component: PhaseNavigationFooter,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PhaseNavigationFooter>

export default meta

type Story = StoryObj<typeof meta>

export const TermReady: Story = {
  args: phaseNavigationFooterFixtures.termReady
}

export const TermBlocked: Story = {
  args: phaseNavigationFooterFixtures.termBlocked
}

export const Complete: Story = {
  args: phaseNavigationFooterFixtures.complete
}
