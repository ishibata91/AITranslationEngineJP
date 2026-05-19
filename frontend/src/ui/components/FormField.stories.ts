import type { Meta, StoryObj } from "@storybook/svelte-vite"
import FormField from "./FormField.svelte"
import { formFieldFixture } from "./__fixtures__/form-field-fixture"

const meta = {
  title: "UI Components/FormField",
  component: FormField,
  args: formFieldFixture
} satisfies Meta<typeof FormField>

export default meta

type Story = StoryObj<typeof meta>

export const LabelHelpRequired: Story = {}

export const Error: Story = {
  args: {
    ...formFieldFixture,
    id: "storybook-form-field-error",
    error: "入力内容を確認してください。"
  }
}

export const Disabled: Story = {
  args: {
    ...formFieldFixture,
    id: "storybook-form-field-disabled",
    disabled: true
  }
}
