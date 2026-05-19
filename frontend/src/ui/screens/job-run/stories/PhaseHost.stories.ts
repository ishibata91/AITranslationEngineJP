import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PhaseHost from "../PhaseHost.svelte"

const meta = {
  title: "Screens/Job Run/PhaseHost",
  component: PhaseHost,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof PhaseHost>

export default meta

type Story = StoryObj<typeof meta>

export const Placeholder: Story = {
  args: {}
}
