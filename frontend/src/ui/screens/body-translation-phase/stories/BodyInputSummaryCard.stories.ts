import type { Meta, StoryObj } from "@storybook/svelte-vite"
import BodyInputSummaryCard from "../BodyInputSummaryCard.svelte"
import { bodyInputSummaryCardFixture } from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screens/Body Translation Phase/BodyInputSummaryCard",
  component: BodyInputSummaryCard,
  args: bodyInputSummaryCardFixture
} satisfies Meta<typeof BodyInputSummaryCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
