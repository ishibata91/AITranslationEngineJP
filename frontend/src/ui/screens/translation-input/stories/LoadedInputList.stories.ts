import type { Meta, StoryObj } from "@storybook/svelte-vite"

import LoadedInputList from "../LoadedInputList.svelte"
import { translationInputPanelFixtures } from "../__fixtures__/translation-input-panel-fixtures"

const meta = {
  title: "Screen Components/Translation Input/LoadedInputList",
  component: LoadedInputList,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof LoadedInputList>

export default meta

type Story = StoryObj<typeof meta>

export const Empty: Story = {
  args: translationInputPanelFixtures.loadedInputList.empty
}

export const Selected: Story = {
  args: translationInputPanelFixtures.loadedInputList.selected
}
