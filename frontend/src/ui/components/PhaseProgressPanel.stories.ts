import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseProgressPanel from "./PhaseProgressPanel.svelte"
import { phaseProgressPanelFixture } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/Phase Progress Panel/PhaseProgressPanel",
  component: PhaseProgressPanel,
  args: phaseProgressPanelFixture
} satisfies Meta<typeof PhaseProgressPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
