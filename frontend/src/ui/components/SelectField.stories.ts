import type { Meta, StoryObj } from "@storybook/svelte-vite"
import SelectField from "./SelectField.svelte"
import { selectFieldFixture } from "./__fixtures__/form-field-fixture"

const meta = {
  title: "UI Components/SelectField",
  component: SelectField,
  args: selectFieldFixture
} satisfies Meta<typeof SelectField>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Error: Story = {
  args: {
    ...selectFieldFixture,
    id: "storybook-select-error",
    value: "",
    error: "分類を選択してください。",
    placeholder: "選んでください"
  }
}

export const Disabled: Story = {
  args: {
    ...selectFieldFixture,
    id: "storybook-select-disabled",
    disabled: true
  }
}
