import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PersonaGenerationPhasePanel from "../PersonaGenerationPhasePanel.svelte"
import { personaGenerationPhasePanelFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screens/Persona Generation Phase/PersonaGenerationPhasePanel",
  component: PersonaGenerationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: personaGenerationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof PersonaGenerationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

export const AiSettingsReady: Story = {}

export const AiSettingsMissing: Story = {
  args: {
    viewModel: {
      ...personaGenerationPhasePanelFixture,
      modelLabel: "-",
      credentialRefLabel: "-"
    }
  }
}

export const RunningLocked: Story = {
  args: {
    viewModel: {
      ...personaGenerationPhasePanelFixture,
      viewState: "running",
      phaseStateLabel: "実行中",
      statusTitle: "NPC ペルソナ生成を実行中です"
    }
  }
}
