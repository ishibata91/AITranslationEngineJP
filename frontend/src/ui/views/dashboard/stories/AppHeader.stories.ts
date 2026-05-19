import type { Meta, StoryObj } from "@storybook/svelte-vite"

import AppHeader from "../AppHeader.svelte"
import { appHeaderFixtures } from "../__fixtures__/dashboard-component-fixtures"

const meta = {
  title: "Screens/Dashboard/AppHeader",
  component: AppHeader,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof AppHeader>

export default meta

type Story = StoryObj<typeof meta>

export const Desktop: Story = {
  args: appHeaderFixtures.desktop
}

export const MobileOpen: Story = {
  args: appHeaderFixtures.mobileOpen
}
