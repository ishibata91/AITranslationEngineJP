import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseMetricCounterGrid from "./PhaseMetricCounterGrid.svelte"
import { phaseMetricCounterGridFixture } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/PhaseMetricCounterGrid",
  component: PhaseMetricCounterGrid,
  args: phaseMetricCounterGridFixture
} satisfies Meta<typeof PhaseMetricCounterGrid>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
