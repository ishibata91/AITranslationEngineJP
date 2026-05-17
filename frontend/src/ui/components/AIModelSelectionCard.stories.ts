import type { Meta, StoryObj } from "@storybook/svelte-vite"
import AIModelSelectionCard from "./AIModelSelectionCard.svelte"
import { aiModelSelectionCardFixture } from "./__fixtures__/ai-model-selection-card-fixture"

const meta = {
  title: "UI Components/AIModelSelectionCard",
  component: AIModelSelectionCard,
  args: aiModelSelectionCardFixture
} satisfies Meta<typeof AIModelSelectionCard>

export default meta

type Story = StoryObj<typeof meta>

export const FixedProps: Story = {}
