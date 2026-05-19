import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TermResultSummaryCard from "../TermResultSummaryCard.svelte"
import { termResultSummaryCardFixture } from "../__fixtures__/term-phase-card-fixture"

const meta = {
  title: "Screens/Term Translation Phase/TermResultSummaryCard",
  component: TermResultSummaryCard,
  args: termResultSummaryCardFixture
} satisfies Meta<typeof TermResultSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
