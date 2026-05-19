import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ActionButton from "./ActionButton.svelte"
import {
  busyActionButtonFixture,
  dangerActionButtonFixture,
  disabledActionButtonFixture,
  longActionButtonFixture,
  primaryActionButtonFixture,
  secondaryActionButtonFixture
} from "./__fixtures__/action-button-fixture"

const meta = {
  title: "UI Components/ActionButton",
  component: ActionButton,
  args: primaryActionButtonFixture
} satisfies Meta<typeof ActionButton>

export default meta

type Story = StoryObj<typeof meta>

export const Primary: Story = {}

export const Secondary: Story = {
  args: secondaryActionButtonFixture
}

export const Danger: Story = {
  args: dangerActionButtonFixture
}

export const Busy: Story = {
  args: busyActionButtonFixture
}

export const Disabled: Story = {
  args: disabledActionButtonFixture
}

export const LongLabel: Story = {
  args: longActionButtonFixture
}
