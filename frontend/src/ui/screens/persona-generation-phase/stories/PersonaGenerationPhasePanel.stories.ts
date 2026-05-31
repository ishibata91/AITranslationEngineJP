import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PersonaGenerationPhasePanel from "../PersonaGenerationPhasePanel.svelte"
import { personaGenerationPhasePanelFixture } from "../__fixtures__/persona-phase-card-fixture"
import { processingTargetListPanelFixtures } from "@ui/screens/job-run/__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Review/Changed Screens/PersonaGenerationPhasePanel",
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

export const WithProcessingTargets: Story = {
  args: {
    processingTargetPageState: {
      ...processingTargetListPanelFixtures.personaGeneration,
      page: 1,
      pageSize: 50,
      totalCount: 64,
      searchQuery: "",
      busy: false
    }
  }
}

export const EmptyProcessingTargets: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [],
      page: 1,
      pageSize: 50,
      totalCount: 0,
      searchQuery: "",
      busy: false
    }
  }
}
