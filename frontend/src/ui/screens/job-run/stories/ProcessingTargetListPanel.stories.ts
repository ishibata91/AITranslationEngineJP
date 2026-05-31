import type { Meta, StoryObj } from "@storybook/svelte-vite"

import ProcessingTargetListPanel from "@ui/components/ProcessingTargetListPanel.svelte"
import { processingTargetListPanelFixtures } from "../__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Review/Changed Components/ProcessingTargetListPanel",
  component: ProcessingTargetListPanel,
  args: processingTargetListPanelFixtures.termTranslationFirstPage,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof ProcessingTargetListPanel>

export default meta

type Story = StoryObj<typeof meta>

export const TermTranslation: Story = {}

export const PersonaGeneration: Story = {
  args: processingTargetListPanelFixtures.personaGeneration
}

export const BodyTranslation: Story = {
  args: processingTargetListPanelFixtures.bodyTranslation
}

export const TranslationComplete: Story = {
  args: processingTargetListPanelFixtures.translationComplete
}

export const FirstPage: Story = {
  args: processingTargetListPanelFixtures.termTranslationFirstPage
}

export const FinalPage: Story = {
  args: processingTargetListPanelFixtures.finalPage
}

export const LongText: Story = {
  args: processingTargetListPanelFixtures.longText
}

export const TooltipLongTitle: Story = {
  args: processingTargetListPanelFixtures.longText
}

export const ValueOnlyTitleParts: Story = {
  args: processingTargetListPanelFixtures.valueonlyTitleParts
}
