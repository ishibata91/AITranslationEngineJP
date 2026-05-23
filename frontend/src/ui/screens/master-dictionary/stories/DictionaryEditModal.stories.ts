import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryEditModal from "../DictionaryEditModal.svelte"
import { dictionaryEditModalFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screen Components/Master Dictionary/DictionaryEditModal",
  component: DictionaryEditModal,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryEditModal>

export default meta

type Story = StoryObj<typeof meta>

export const Create: Story = {
  args: dictionaryEditModalFixtures.create
}

export const Edit: Story = {
  args: dictionaryEditModalFixtures.edit
}

export const SaveFailed: Story = {
  args: dictionaryEditModalFixtures.saveFailed
}
