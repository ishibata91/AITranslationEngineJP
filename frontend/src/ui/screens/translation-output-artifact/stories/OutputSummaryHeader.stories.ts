import type { Meta, StoryObj } from "@storybook/svelte-vite"

import OutputSummaryHeader from "../OutputSummaryHeader.svelte"
import { outputSummaryHeaderFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screen Components/Translation Output Artifact/OutputSummaryHeader",
  component: OutputSummaryHeader,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof OutputSummaryHeader>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: outputSummaryHeaderFixtures.ready
}
