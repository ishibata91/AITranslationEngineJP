import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ProgressBar from "./ProgressBar.svelte"
import {
  progressDoneFixture,
  progressFailureFixture,
  progressHalfFixture,
  progressLongFixture,
  progressZeroFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/ProgressBar",
  component: ProgressBar,
  args: progressHalfFixture
} satisfies Meta<typeof ProgressBar>

export default meta

type Story = StoryObj<typeof meta>

export const Zero: Story = {
  args: progressZeroFixture
}

export const Half: Story = {}

export const Done: Story = {
  args: progressDoneFixture
}

export const Failure: Story = {
  args: progressFailureFixture
}

export const LongLabel: Story = {
  args: progressLongFixture
}
