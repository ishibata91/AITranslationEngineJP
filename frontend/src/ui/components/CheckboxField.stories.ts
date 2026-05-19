import type { Meta, StoryObj } from "@storybook/svelte-vite"
import CheckboxField from "./CheckboxField.svelte"
import {
  checkboxFieldErrorFixture,
  checkboxFieldFixture
} from "./__fixtures__/form-field-fixture"

const meta = {
  title: "UI Components/CheckboxField",
  component: CheckboxField,
  args: checkboxFieldFixture
} satisfies Meta<typeof CheckboxField>

export default meta

type Story = StoryObj<typeof meta>

export const Checked: Story = {}

export const Error: Story = {
  args: checkboxFieldErrorFixture
}

export const Disabled: Story = {
  args: {
    ...checkboxFieldFixture,
    id: "storybook-checkbox-disabled",
    disabled: true
  }
}
