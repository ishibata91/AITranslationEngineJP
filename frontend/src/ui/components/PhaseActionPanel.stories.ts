import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PhaseActionPanel from "./PhaseActionPanel.svelte"
import { phaseActionPanelFixture } from "./__fixtures__/phase-panel-fixture"

const meta = {
  title: "UI Components/PhaseActionPanel",
  component: PhaseActionPanel,
  args: phaseActionPanelFixture
} satisfies Meta<typeof PhaseActionPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
