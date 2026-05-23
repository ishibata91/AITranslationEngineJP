import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ProviderDetailPanel from "../ProviderDetailPanel.svelte"
import { providerDetailPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/ProviderDetailPanel",
  component: ProviderDetailPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof ProviderDetailPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Selected: Story = {
  args: providerDetailPanelFixtures.selected
}

export const Warning: Story = {
  args: providerDetailPanelFixtures.warning
}
