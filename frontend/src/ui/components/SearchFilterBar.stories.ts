import type { Meta, StoryObj } from "@storybook/svelte-vite"
import SearchFilterBar from "./SearchFilterBar.svelte"
import {
  searchFilterDefaultFixture,
  searchFilterEmptyFixture,
  searchFilterLongFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/SearchFilterBar",
  component: SearchFilterBar,
  args: searchFilterDefaultFixture
} satisfies Meta<typeof SearchFilterBar>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Empty: Story = {
  args: searchFilterEmptyFixture
}

export const LongText: Story = {
  args: searchFilterLongFixture
}
