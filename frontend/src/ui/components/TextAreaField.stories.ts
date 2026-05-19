import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TextAreaField from "./TextAreaField.svelte"
import { textAreaFieldFixture } from "./__fixtures__/form-field-fixture"

const meta = {
  title: "UI Components/TextAreaField",
  component: TextAreaField,
  args: textAreaFieldFixture
} satisfies Meta<typeof TextAreaField>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const Error: Story = {
  args: {
    ...textAreaFieldFixture,
    id: "storybook-text-area-error",
    error: "本文を入力してください。"
  }
}

export const Disabled: Story = {
  args: {
    ...textAreaFieldFixture,
    id: "storybook-text-area-disabled",
    disabled: true
  }
}
