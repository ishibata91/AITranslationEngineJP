import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BodyReadinessInputCard from "../BodyReadinessInputCard.svelte"
import { bodyReadinessInputCardFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screens/Persona Generation Phase/BodyReadinessInputCard",
  component: BodyReadinessInputCard,
  args: bodyReadinessInputCardFixture
} satisfies Meta<typeof BodyReadinessInputCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
