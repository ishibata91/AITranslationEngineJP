import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DashboardEntryGrid from "../DashboardEntryGrid.svelte"
import { dashboardEntryGridFixtures } from "../__fixtures__/dashboard-component-fixtures"

const meta = {
  title: "Screens/Dashboard/DashboardEntryGrid",
  component: DashboardEntryGrid,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof DashboardEntryGrid>

export default meta

type Story = StoryObj<typeof meta>

export const Standard: Story = {
  args: dashboardEntryGridFixtures.standard
}
