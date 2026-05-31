import type { Meta, StoryObj } from "@storybook/svelte-vite"

import TermTranslationPhasePanel from "../TermTranslationPhasePanel.svelte"
import { termTranslationPhasePanelFixture } from "../__fixtures__/term-phase-card-fixture"
import { processingTargetListPanelFixtures } from "@ui/screens/job-run/__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Review/Changed Screens/TermTranslationPhasePanel",
  component: TermTranslationPhasePanel,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    viewModel: termTranslationPhasePanelFixture,
    onAction: () => {}
  }
} satisfies Meta<typeof TermTranslationPhasePanel>

export default meta

type Story = StoryObj<typeof meta>

export const AiSettingsReady: Story = {}

export const AiSettingsMissing: Story = {
  args: {
    viewModel: {
      ...termTranslationPhasePanelFixture,
      modelLabel: "-",
      credentialRefLabel: "-"
    }
  }
}

export const ConflictingRunActions: Story = {
  args: {
    viewModel: {
      ...termTranslationPhasePanelFixture,
      actionCards: [
        {
          id: "start",
          label: "開始",
          disabled: false,
          blockedReason: "",
          tone: "primary"
        },
        {
          id: "pause",
          label: "中断",
          disabled: true,
          blockedReason: "実行中ではありません。",
          tone: "warning"
        },
        {
          id: "resume",
          label: "再開",
          disabled: false,
          blockedReason: "",
          tone: "default"
        },
        {
          id: "retry",
          label: "リトライ",
          disabled: false,
          blockedReason: "",
          tone: "default"
        },
        {
          id: "next-phase",
          label: "次の翻訳段階へ進む",
          disabled: false,
          blockedReason: "",
          tone: "primary"
        }
      ]
    }
  }
}

export const RunningLocked: Story = {
  args: {
    viewModel: {
      ...termTranslationPhasePanelFixture,
      viewState: "running",
      phaseStateLabel: "実行中",
      statusTitle: "単語翻訳を実行中です"
    }
  }
}

export const FailureStatus: Story = {
  args: {
    viewModel: {
      ...termTranslationPhasePanelFixture,
      viewState: "failed",
      phaseStateLabel: "失敗",
      statusTitle: "単語翻訳が失敗しました",
      statusText: "失敗理由を確認して再実行を判断します。",
      errorMessage: "AI サービスの応答を取得できませんでした。"
    }
  }
}

export const WithProcessingTargets: Story = {
  args: {
    processingTargetPageState: {
      ...processingTargetListPanelFixtures.termTranslationFirstPage,
      page: 1,
      pageSize: 50,
      totalCount: 125,
      searchQuery: "",
      busy: false
    }
  }
}

export const WithValueOnlyTitleParts: Story = {
  args: {
    processingTargetPageState: {
      ...processingTargetListPanelFixtures.valueonlyTitleParts,
      page: 1,
      pageSize: 50,
      totalCount: 10,
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
