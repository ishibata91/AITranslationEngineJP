import type { Meta, StoryObj } from "@storybook/svelte-vite"

import LatestOutputResultCard from "../LatestOutputResultCard.svelte"
import { latestOutputResultCardFixtures } from "../__fixtures__/translation-output-artifact-fixtures"

const meta = {
  title: "Screens/Translation Output Artifact/LatestOutputResultCard",
  component: LatestOutputResultCard,
  parameters: {
    layout: "fullscreen"
  }
} satisfies Meta<typeof LatestOutputResultCard>

export default meta

type Story = StoryObj<typeof meta>

export const Generated: Story = {
  args: latestOutputResultCardFixtures.generated
}

export const Failed: Story = {
  args: latestOutputResultCardFixtures.failed
}

export const Empty: Story = {
  args: latestOutputResultCardFixtures.empty
}
