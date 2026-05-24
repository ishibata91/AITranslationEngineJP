import type { Meta, StoryObj } from "@storybook/svelte-vite"

import SettingsActionPanel from "../SettingsActionPanel.svelte"
import { settingsActionPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/SettingsActionPanel",
  component: SettingsActionPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof SettingsActionPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Enabled: Story = {
  args: settingsActionPanelFixtures.enabled
}

export const Disabled: Story = {
  args: settingsActionPanelFixtures.disabled
}
