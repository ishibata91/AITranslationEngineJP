import type { Meta, StoryObj } from "@storybook/svelte-vite"

import PersonaGenerationPhasePanel from "../PersonaGenerationPhasePanel.svelte"
import { personaGenerationPhasePanelFixture } from "../__fixtures__/persona-phase-card-fixture"

const meta = {
  title: "Screen Components/Persona Generation Phase/PersonaGenerationPhasePanel",
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

/**
 * 初回取得中ローディングレイヤー表示（initialFetchDone=false）。
 * selector: persona-generation-phase-processing-target-loading
 */
export const ProcessingTargetLoading: Story = {
  args: {
    initialFetchDone: false
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、件数あり状態。
 */
export const ProcessingTargetWithItems: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [
        {
          id: "npc:1",
          name: "Vilkas",
          detail: "NPC: Vilkas",
          titleParts: [{ text: "Vilkas" }],
          metadata: [{ label: "種族", value: "Nord" }]
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "",
      busy: false
    }
  }
}

/**
 * 初回取得完了後（initialFetchDone=true）、空状態（母数0件）。
 * selector: persona-generation-phase-processing-target-empty
 */
export const ProcessingTargetEmpty: Story = {
  args: {
    initialFetchDone: true,
    processingTargetPageState: {
      items: [],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 0,
      searchQuery: "",
      busy: false
    }
  }
}
