import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ProviderSettingsSummaryPanel from "../ProviderSettingsSummaryPanel.svelte"
import { providerSettingsSummaryPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/ProviderSettingsSummaryPanel",
  component: ProviderSettingsSummaryPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof ProviderSettingsSummaryPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: providerSettingsSummaryPanelFixtures.ready
}

export const Saved: Story = {
  args: providerSettingsSummaryPanelFixtures.saved
}

export const Failed: Story = {
  args: providerSettingsSummaryPanelFixtures.failed
}
