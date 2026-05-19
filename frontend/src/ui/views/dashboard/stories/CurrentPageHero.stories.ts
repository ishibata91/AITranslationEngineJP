import type { Meta, StoryObj } from "@storybook/svelte-vite"

import CurrentPageHero from "../CurrentPageHero.svelte"
import { currentPageHeroFixtures } from "../__fixtures__/dashboard-component-fixtures"

const meta = {
  title: "Screens/Dashboard/CurrentPageHero",
  component: CurrentPageHero,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof CurrentPageHero>

export default meta

type Story = StoryObj<typeof meta>

export const Dashboard: Story = {
  args: currentPageHeroFixtures.dashboard
}

export const ProviderSettings: Story = {
  args: currentPageHeroFixtures.providerSettings
}
