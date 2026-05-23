import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ProviderListPanel from "../ProviderListPanel.svelte"
import { providerListPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/ProviderListPanel",
  component: ProviderListPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof ProviderListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Mixed: Story = {
  args: providerListPanelFixtures.mixed
}

export const Empty: Story = {
  args: providerListPanelFixtures.empty
}
