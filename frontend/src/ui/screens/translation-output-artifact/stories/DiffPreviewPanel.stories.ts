import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DiffPreviewPanel from "../DiffPreviewPanel.svelte"
import { diffPreviewPanelFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screens/Translation Output Artifact/DiffPreviewPanel",
  component: DiffPreviewPanel,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DiffPreviewPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Populated: Story = {
  args: diffPreviewPanelFixtures.populated
}

export const Empty: Story = {
  args: diffPreviewPanelFixtures.empty
}
