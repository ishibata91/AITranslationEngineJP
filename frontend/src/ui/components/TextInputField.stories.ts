import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TextInputField from "./TextInputField.svelte"
import {
  textInputDisabledFixture,
  textInputErrorFixture,
  textInputFieldFixture
} from "./__fixtures__/form-field-fixture"

const meta = {
  title: "UI Components/TextInputField",
  component: TextInputField,
  args: textInputFieldFixture
} satisfies Meta<typeof TextInputField>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Error: Story = {
  args: textInputErrorFixture
}

export const Disabled: Story = {
  args: textInputDisabledFixture
}

export const LongValue: Story = {
  args: {
    ...textInputFieldFixture,
    id: "storybook-text-input-long",
    value: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  }
}
