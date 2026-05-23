import type { Meta, StoryObj } from "@storybook/svelte-vite"
import AIModelSelectionCard from "./AIModelSelectionCard.svelte"
import {
  aiModelSelectionCardFixture,
  aiModelSelectionCardStates
} from "./__fixtures__/ai-model-selection-card-fixture"

const meta = {
  title: "UI Components/AIModelSelectionCard",
  component: AIModelSelectionCard,
  args: aiModelSelectionCardFixture
} satisfies Meta<typeof AIModelSelectionCard>

export default meta

type Story = StoryObj<typeof meta>

export const FixedProps: Story = {}

export const ModelListLoading: Story = {
  args: aiModelSelectionCardStates.loading
}

export const ModelListFailed: Story = {
  args: aiModelSelectionCardStates.failed
}

export const CredentialMissing: Story = {
  args: aiModelSelectionCardStates.credentialMissing
}

export const RunningLocked: Story = {
  args: aiModelSelectionCardStates.runningLocked
}
