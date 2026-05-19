import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseFailureInfoCard from "./PhaseFailureInfoCard.svelte"
import { phaseFailureInfoCardFixture } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/PhaseFailureInfoCard",
  component: PhaseFailureInfoCard,
  args: phaseFailureInfoCardFixture
} satisfies Meta<typeof PhaseFailureInfoCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
