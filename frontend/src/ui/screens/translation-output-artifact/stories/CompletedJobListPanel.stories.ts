import type { Meta, StoryObj } from "@storybook/svelte-vite"

import CompletedJobListPanel from "../CompletedJobListPanel.svelte"
import { completedJobListPanelFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screen Components/Translation Output Artifact/CompletedJobListPanel",
  component: CompletedJobListPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof CompletedJobListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Populated: Story = {
  args: completedJobListPanelFixtures.populated
}

export const Empty: Story = {
  args: completedJobListPanelFixtures.empty
}
