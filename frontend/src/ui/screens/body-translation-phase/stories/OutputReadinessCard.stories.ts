import type { Meta, StoryObj } from "@storybook/svelte-vite"
import OutputReadinessCard from "../OutputReadinessCard.svelte"
import {
  blockedOutputReadinessCardFixture,
  outputReadinessCardFixture
} from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screens/Body Translation Phase/OutputReadinessCard",
  component: OutputReadinessCard,
  args: outputReadinessCardFixture
} satisfies Meta<typeof OutputReadinessCard>

export default meta

type Story = StoryObj<typeof meta>

export const Ready: Story = {}

export const Blocked: Story = {
  args: blockedOutputReadinessCardFixture
}
