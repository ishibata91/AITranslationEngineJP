import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ButtonGroupStoryDemo from "./ButtonGroupStoryDemo.svelte"
import {
  buttonGroupFixture,
  buttonGroupStretchFixture
} from "./__fixtures__/button-group-fixture"

const meta = {
  title: "UI Components/ButtonGroup",
  component: ButtonGroupStoryDemo,
  args: buttonGroupFixture
} satisfies Meta<typeof ButtonGroupStoryDemo>

export default meta

type Story = StoryObj<typeof meta>

export const EndAligned: Story = {}

export const StretchAligned: Story = {
  args: buttonGroupStretchFixture
}
