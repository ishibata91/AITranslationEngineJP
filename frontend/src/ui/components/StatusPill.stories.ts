import type { Meta, StoryObj } from "@storybook/svelte-vite"
import StatusPill from "./StatusPill.svelte"
import {
  statusPillBusyFixture,
  statusPillDefaultFixture,
  statusPillFailureFixture,
  statusPillLongFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/StatusPill",
  component: StatusPill,
  args: statusPillDefaultFixture
} satisfies Meta<typeof StatusPill>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Busy: Story = {
  args: statusPillBusyFixture
}

export const Failure: Story = {
  args: statusPillFailureFixture
}

export const LongLabel: Story = {
  args: statusPillLongFixture
}
