import type { Meta, StoryObj } from "@storybook/svelte-vite"

import LoadedInputDetail from "../LoadedInputDetail.svelte"
import { translationInputPanelFixtures } from "../__fixtures__/translation-input-panel-fixtures"

const meta = {
  title: "Screens/Translation Input/LoadedInputDetail",
  component: LoadedInputDetail,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof LoadedInputDetail>

export default meta

type Story = StoryObj<typeof meta>

export const Empty: Story = {
  args: translationInputPanelFixtures.loadedInputDetail.empty
}

export const Selected: Story = {
  args: translationInputPanelFixtures.loadedInputDetail.selected
}
