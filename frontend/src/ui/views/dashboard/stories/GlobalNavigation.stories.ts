import type { Meta, StoryObj } from "@storybook/svelte-vite"

import GlobalNavigation from "../GlobalNavigation.svelte"
import { globalNavigationFixtures } from "../__fixtures__/dashboard-component-fixtures"

const meta = {
  title: "Screens/Dashboard/GlobalNavigation",
  component: GlobalNavigation,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof GlobalNavigation>

export default meta

type Story = StoryObj<typeof meta>

export const Dashboard: Story = {
  args: globalNavigationFixtures.dashboard
}
