import type { Meta, StoryObj } from "@storybook/svelte-vite"

import InputReviewPage from "../InputReviewPage.svelte"
import { createInputReviewPageControllerFixture } from "../__fixtures__/translation-input-panel-fixtures"

const meta = {
  title: "Screens/Translation Input/InputReviewPage",
  component: InputReviewPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createController: createInputReviewPageControllerFixture(),
    onOpenJobRun: () => {}
  }
} satisfies Meta<typeof InputReviewPage>

export default meta

type Story = StoryObj<typeof meta>

export const SelectedInputReady: Story = {}
