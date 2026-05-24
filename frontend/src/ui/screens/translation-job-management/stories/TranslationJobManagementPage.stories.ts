import type { Meta, StoryObj } from "@storybook/svelte-vite"

import { createTranslationJobManagementPageControllerFixture } from "../../__fixtures__/screen-page-controller-fixtures"
import TranslationJobManagementPage from "../TranslationJobManagementPage.svelte"

const meta = {
  title: "Screens/Translation Job Management/TranslationJobManagementPage",
  component: TranslationJobManagementPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createTranslationJobManagementPageControllerFixture(),
    onJobRunTargetChange: () => undefined,
    onOpenInputReview: () => undefined,
    onOpenJobRun: () => undefined
  }
} satisfies Meta<typeof TranslationJobManagementPage>

export default meta

type Story = StoryObj<typeof meta>

export const Disconnected: Story = {}
