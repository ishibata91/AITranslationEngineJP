import type { Meta, StoryObj } from "@storybook/svelte-vite"
import ConfirmDangerModal from "./ConfirmDangerModal.svelte"
import {
  confirmDangerBusyFixture,
  confirmDangerDefaultFixture,
  confirmDangerFailureFixture,
  confirmDangerLongFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/ConfirmDangerModal",
  component: ConfirmDangerModal,
  args: confirmDangerDefaultFixture
} satisfies Meta<typeof ConfirmDangerModal>

export default meta

type Story = StoryObj<typeof meta>

export const DangerOperation: Story = {}

export const Busy: Story = {
  args: confirmDangerBusyFixture
}

export const Failure: Story = {
  args: confirmDangerFailureFixture
}

export const LongTarget: Story = {
  args: confirmDangerLongFixture
}
