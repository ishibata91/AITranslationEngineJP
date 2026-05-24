import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PersonaActionModal from "../PersonaActionModal.svelte"
import { personaActionModalFixtures } from "../__fixtures__/master-persona-panel-fixtures"

const meta = {
  title: "Screen Components/Master Persona/PersonaActionModal",
  component: PersonaActionModal,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PersonaActionModal>

export default meta

type Story = StoryObj<typeof meta>

export const Editing: Story = {
  args: personaActionModalFixtures.editing
}

export const Deleting: Story = {
  args: personaActionModalFixtures.deleting
}

export const SaveFailed: Story = {
  args: personaActionModalFixtures.saveFailed
}

export const DeleteFailed: Story = {
  args: personaActionModalFixtures.deleteFailed
}
