import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PersonaResultSummaryCard from "../PersonaResultSummaryCard.svelte"
import { personaResultSummaryCardFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screens/Persona Generation Phase/PersonaResultSummaryCard",
  component: PersonaResultSummaryCard,
  args: personaResultSummaryCardFixture
} satisfies Meta<typeof PersonaResultSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
