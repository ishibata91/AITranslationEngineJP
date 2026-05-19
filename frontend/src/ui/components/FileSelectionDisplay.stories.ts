import type { Meta, StoryObj } from "@storybook/svelte-vite"
import FileSelectionDisplay from "./FileSelectionDisplay.svelte"
import {
  fileSelectionBusyFixture,
  fileSelectionDefaultFixture,
  fileSelectionEmptyFixture,
  fileSelectionFailureFixture,
  fileSelectionLongFixture
} from "./__fixtures__/status-list-primitives-fixture"

const meta = {
  title: "UI Components/FileSelectionDisplay",
  component: FileSelectionDisplay,
  args: fileSelectionDefaultFixture
} satisfies Meta<typeof FileSelectionDisplay>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Empty: Story = {
  args: fileSelectionEmptyFixture
}

export const Failure: Story = {
  args: fileSelectionFailureFixture
}

export const Busy: Story = {
  args: fileSelectionBusyFixture
}

export const LongPath: Story = {
  args: fileSelectionLongFixture
}
