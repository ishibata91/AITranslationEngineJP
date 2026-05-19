import type { Meta, StoryObj } from "@storybook/svelte-vite"
import TermExecutionSettingsCard from "../TermExecutionSettingsCard.svelte"
import { termExecutionSettingsCardFixture } from "../__fixtures__/term-phase-card-fixture"

const meta = {
  title: "Screens/Term Translation Phase/TermExecutionSettingsCard",
  component: TermExecutionSettingsCard,
  args: termExecutionSettingsCardFixture
} satisfies Meta<typeof TermExecutionSettingsCard>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {}
