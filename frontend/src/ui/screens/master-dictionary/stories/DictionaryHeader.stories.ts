import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryHeader from "../DictionaryHeader.svelte"
import { dictionaryHeaderFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screens/Master Dictionary/DictionaryHeader",
  component: DictionaryHeader,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryHeader>

export default meta

type Story = StoryObj<typeof meta>

export const Normal: Story = {
  args: dictionaryHeaderFixtures.normal
}

export const Error: Story = {
  args: dictionaryHeaderFixtures.error
}
