import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PersonaTargetSummaryCard from "../PersonaTargetSummaryCard.svelte"
import { personaTargetSummaryCardFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screens/Persona Generation Phase/PersonaTargetSummaryCard",
  component: PersonaTargetSummaryCard,
  args: personaTargetSummaryCardFixture
} satisfies Meta<typeof PersonaTargetSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
