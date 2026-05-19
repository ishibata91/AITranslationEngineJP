import type { Meta, StoryObj } from "@storybook/svelte-vite"
import InlineFeedback from "./InlineFeedback.svelte"
import {
  errorInlineFeedbackFixture,
  longInlineFeedbackFixture,
  neutralInlineFeedbackFixture,
  successInlineFeedbackFixture,
  warningInlineFeedbackFixture
} from "./__fixtures__/inline-feedback-fixture"

const meta = {
  title: "UI Components/InlineFeedback",
  component: InlineFeedback,
  args: neutralInlineFeedbackFixture
} satisfies Meta<typeof InlineFeedback>

export default meta

type Story = StoryObj<typeof meta>

export const Neutral: Story = {}

export const Error: Story = {
  args: errorInlineFeedbackFixture
}

export const Warning: Story = {
  args: warningInlineFeedbackFixture
}

export const Success: Story = {
  args: successInlineFeedbackFixture
}

export const LongMessage: Story = {
  args: longInlineFeedbackFixture
}
