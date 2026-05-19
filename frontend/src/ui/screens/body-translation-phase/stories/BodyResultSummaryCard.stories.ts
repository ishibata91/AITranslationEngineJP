import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BodyResultSummaryCard from "../BodyResultSummaryCard.svelte"
import { bodyResultSummaryCardFixture } from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screens/Body Translation Phase/BodyResultSummaryCard",
  component: BodyResultSummaryCard,
  args: bodyResultSummaryCardFixture
} satisfies Meta<typeof BodyResultSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
