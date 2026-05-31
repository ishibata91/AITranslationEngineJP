import type { Meta, StoryObj } from "@storybook/svelte-vite"

import BodyTranslationPhasePanel from "../BodyTranslationPhasePanel.svelte"
import { bodyTranslationPhasePanelFixture } from "../__fixtures__/body-phase-card-fixture"
import { processingTargetListPanelFixtures } from "@ui/screens/job-run/__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Review/Changed Screens/BodyTranslationPhasePanel",
  component: BodyTranslationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: bodyTranslationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof BodyTranslationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

export const AiSettingsReady: Story = {}

export const AiSettingsMissing: Story = {
  args: {
    viewModel: {
      ...bodyTranslationPhasePanelFixture,
      modelLabel: "-",
      credentialRefLabel: "-"
    }
  }
}

export const RunningLocked: Story = {
  args: {
    viewModel: {
      ...bodyTranslationPhasePanelFixture,
      viewState: "running",
      phaseStateLabel: "実行中",
      statusTitle: "本文翻訳を実行中です"
    }
  }
}

export const WithProcessingTargets: Story = {
  args: {
    processingTargetPageState: {
      ...processingTargetListPanelFixtures.bodyTranslation,
      page: 1,
      pageSize: 50,
      totalCount: 81,
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
