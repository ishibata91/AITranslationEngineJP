import type { Meta, StoryObj } from "@storybook/svelte-vite"
import PaginationControls from "./PaginationControls.svelte"
import {
  paginationBusyFixture,
  paginationDefaultFixture,
  paginationFirstFixture,
  paginationLastFixture,
  paginationLongFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/PaginationControls",
  component: PaginationControls,
  args: paginationDefaultFixture
} satisfies Meta<typeof PaginationControls>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const FirstPage: Story = {
  args: paginationFirstFixture
}

export const LastPage: Story = {
  args: paginationLastFixture
}

export const Busy: Story = {
  args: paginationBusyFixture
}

export const LongTotal: Story = {
  args: paginationLongFixture
}
