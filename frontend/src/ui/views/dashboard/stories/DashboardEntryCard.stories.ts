import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DashboardEntryCard from "../DashboardEntryCard.svelte"
import { dashboardEntryCardFixtures } from "../__fixtures__/dashboard-component-fixtures"

const meta = {
  title: "Screens/Dashboard/DashboardEntryCard",
  component: DashboardEntryCard,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof DashboardEntryCard>

export default meta

type Story = StoryObj<typeof meta>

export const ProviderSettings: Story = {
  args: dashboardEntryCardFixtures.providerSettings
}

export const LongLabel: Story = {
  args: dashboardEntryCardFixtures.longLabel
}
