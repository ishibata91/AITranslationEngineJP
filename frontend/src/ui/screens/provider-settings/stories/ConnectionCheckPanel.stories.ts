import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ConnectionCheckPanel from "../ConnectionCheckPanel.svelte"
import { connectionCheckPanelFixtures } from "../__fixtures__/provider-settings-panel-fixtures"

const meta = {
  title: "Screen Components/Provider Settings/ConnectionCheckPanel",
  component: ConnectionCheckPanel,
  parameters: { layout: "fullscreen" }
} satisfies Meta<typeof ConnectionCheckPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: connectionCheckPanelFixtures.ready
}

export const Disabled: Story = {
  args: connectionCheckPanelFixtures.disabled
}
