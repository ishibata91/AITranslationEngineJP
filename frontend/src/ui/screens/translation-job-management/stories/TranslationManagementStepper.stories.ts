import type { Meta, StoryObj } from "@storybook/svelte-vite"

import TranslationManagementStepper from "../TranslationManagementStepper.svelte"
import { stepperFixtures } from "../__fixtures__/translation-job-management-fixtures"

const meta = {
  title: "UI Components/TranslationManagementStepper",
  component: TranslationManagementStepper,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof TranslationManagementStepper>

export default meta

type Story = StoryObj<typeof meta>

export const JobManagementCurrent: Story = {
  args: stepperFixtures.jobManagementCurrent
}

export const TermTranslationCurrent: Story = {
  args: stepperFixtures.termTranslationCurrent
}
