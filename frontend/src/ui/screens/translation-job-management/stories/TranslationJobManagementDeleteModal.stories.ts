import type { Meta, StoryObj } from "@storybook/svelte-vite"

import TranslationJobManagementDeleteModal from "../TranslationJobManagementDeleteModal.svelte"
import { deleteModalFixtures } from "../__fixtures__/translation-job-management-fixtures"

const meta = {
  title: "Screen Components/Translation Job Management/DeleteModal",
  component: TranslationJobManagementDeleteModal,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof TranslationJobManagementDeleteModal>

export default meta

type Story = StoryObj<typeof meta>

export const Closed: Story = {
  args: deleteModalFixtures.closed
}

export const Open: Story = {
  args: deleteModalFixtures.open
}

export const Deleting: Story = {
  args: deleteModalFixtures.deleting
}

export const DeleteFailed: Story = {
  args: deleteModalFixtures.deleteFailed
}
