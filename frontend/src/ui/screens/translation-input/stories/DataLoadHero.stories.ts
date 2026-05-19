import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DataLoadHero from "../DataLoadHero.svelte"
import { translationInputPanelFixtures } from "../__fixtures__/translation-input-panel-fixtures"

const meta = {
  title: "Screens/Translation Input/DataLoadHero",
  component: DataLoadHero,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DataLoadHero>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {
  args: translationInputPanelFixtures.dataLoadHero.ready
}

export const Failed: Story = {
  args: translationInputPanelFixtures.dataLoadHero.failed
}
