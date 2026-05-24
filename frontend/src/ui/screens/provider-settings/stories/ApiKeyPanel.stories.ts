import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ApiKeyPanel from "../ApiKeyPanel.svelte"
import { apiKeyPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/ApiKeyPanel",
  component: ApiKeyPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof ApiKeyPanel>

export default meta

type Story = StoryObj<typeof meta>

export const MaskedState: Story = {
  args: apiKeyPanelFixtures.maskedState
}

export const DraftOpen: Story = {
  args: apiKeyPanelFixtures.draftOpen
}

export const NotRequired: Story = {
  args: apiKeyPanelFixtures.notRequired
}
