import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PersonaExecutionSettingsCard from "../PersonaExecutionSettingsCard.svelte"
import { personaExecutionSettingsCardFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screens/Persona Generation Phase/PersonaExecutionSettingsCard",
  component: PersonaExecutionSettingsCard,
  args: personaExecutionSettingsCardFixture
} satisfies Meta<typeof PersonaExecutionSettingsCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
