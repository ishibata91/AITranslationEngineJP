import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BodyExecutionSummaryCard from "../BodyExecutionSummaryCard.svelte"
import { bodyExecutionSummaryCardFixture } from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screens/Body Translation Phase/BodyExecutionSummaryCard",
  component: BodyExecutionSummaryCard,
  args: bodyExecutionSummaryCardFixture
} satisfies Meta<typeof BodyExecutionSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
