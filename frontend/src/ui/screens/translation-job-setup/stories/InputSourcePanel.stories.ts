import type { Meta, StoryObj } from "@storybook/svelte-vite"

import InputSourcePanel from "../InputSourcePanel.svelte"
import { translationJobSetupPanelFixtures } from "../__fixtures__/translation-job-setup-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Job Setup/InputSourcePanel",
  component: InputSourcePanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof InputSourcePanel>

export default meta

type Story = StoryObj<typeof meta>

export const Selected: Story = {
  args: translationJobSetupPanelFixtures.inputSourcePanel.selected
}
