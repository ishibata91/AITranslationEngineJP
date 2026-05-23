import type { Meta, StoryObj } from "@storybook/svelte-vite"

import DictionaryDeleteModal from "../DictionaryDeleteModal.svelte"
import { dictionaryDeleteModalFixtures } from "../__fixtures__/master-dictionary-panel-fixtures"

const meta = {
  title: "Screen Components/Master Dictionary/DictionaryDeleteModal",
  component: DictionaryDeleteModal,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof DictionaryDeleteModal>

export default meta

type Story = StoryObj<typeof meta>

export const Confirm: Story = {
  args: dictionaryDeleteModalFixtures.confirm
}

export const DeleteFailed: Story = {
  args: dictionaryDeleteModalFixtures.deleteFailed
}
