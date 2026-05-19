import type { Meta, StoryObj } from "@storybook/svelte-vite"
import IconActionButton from "./IconActionButton.svelte"
import {
  iconActionButtonBusyFixture,
  iconActionButtonDisabledFixture,
  iconActionButtonFixture
} from "./__fixtures__/action-button-fixture"

const meta = {
  title: "UI Components/IconActionButton",
  component: IconActionButton,
  args: iconActionButtonFixture
} satisfies Meta<typeof IconActionButton>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Busy: Story = {
  args: iconActionButtonBusyFixture
}

export const Disabled: Story = {
  args: iconActionButtonDisabledFixture
}
