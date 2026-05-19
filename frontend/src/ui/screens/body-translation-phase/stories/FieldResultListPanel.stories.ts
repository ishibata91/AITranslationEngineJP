import type { Meta, StoryObj } from "@storybook/svelte-vite"
import FieldResultListPanel from "../FieldResultListPanel.svelte"
import {
  emptyFieldResultListPanelFixture,
  fieldResultListPanelFixture
} from "../__fixtures__/body-phase-card-fixture"

const meta = {
  title: "Screens/Body Translation Phase/FieldResultListPanel",
  component: FieldResultListPanel,
  args: fieldResultListPanelFixture
} satisfies Meta<typeof FieldResultListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Empty: Story = {
  args: emptyFieldResultListPanelFixture
}
