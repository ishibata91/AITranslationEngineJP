import type { Meta, StoryObj } from "@storybook/svelte-vite"
import EmptyStatePanel from "./EmptyStatePanel.svelte"
import {
  emptyStateBusyFixture,
  emptyStateDefaultFixture,
  emptyStateFailureFixture,
  emptyStateLongFixture,
  emptyStateWarningFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/EmptyStatePanel",
  component: EmptyStatePanel,
  args: emptyStateDefaultFixture
} satisfies Meta<typeof EmptyStatePanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Warning: Story = {
  args: emptyStateWarningFixture
}

export const Failure: Story = {
  args: emptyStateFailureFixture
}

export const Busy: Story = {
  args: emptyStateBusyFixture
}

export const LongMessage: Story = {
  args: emptyStateLongFixture
}
