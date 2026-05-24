import type { Meta, StoryObj } from "@storybook/svelte-vite"
import type { BodyTranslationPhaseScreenViewModel } from "@application/contract/body-translation-phase"
import type { PersonaGenerationPhaseScreenViewModel } from "@application/contract/persona-generation-phase"
import type { TermTranslationPhaseScreenViewModel } from "@application/contract/term-translation-phase"
import type {
  TranslationJobManagementCurrentPhase
} from "@application/gateway-contract/translation-job-management"
import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"

import {
  createBodyTranslationPhasePageControllerFixture,
  createPersonaGenerationPhasePageControllerFixture,
  createTermTranslationPhasePageControllerFixture
} from "../../__fixtures__/screen-page-controller-fixtures"
import JobRunPage from "../JobRunPage.svelte"
import {
  jobRunTargetSummaryFixtures,
  processingTargetListPanelFixtures
} from "../__fixtures__/job-run-shell-fixtures"

const meta = {
  title: "Screens/Job Run/JobRunPage",
  component: JobRunPage,
  parameters: {
    layout: "fullscreen"
  },
  args: {
    createBodyController: createBodyTranslationPhasePageControllerFixture(),
    createController: createTermTranslationPhasePageControllerFixture(),
    createPersonaController: createPersonaGenerationPhasePageControllerFixture(),
    processingTargetItemsByPhase: {
      term: processingTargetListPanelFixtures.termTranslationFirstPage.items,
      persona: processingTargetListPanelFixtures.personaGeneration.items,
      body: processingTargetListPanelFixtures.bodyTranslation.items,
      complete: processingTargetListPanelFixtures.translationComplete.items
    },
    selectedJobTarget: jobRunTargetSummaryFixtures.termPhase.target,
    onOpenJobManagement: () => undefined,
    onOpenOutputManagement: () => undefined,
    onPhaseViewChange: () => undefined
  }
} satisfies Meta<typeof JobRunPage>

export default meta

type Story = StoryObj<typeof meta>

export const TermTranslation: Story = {}

const termWaiting = {
  phaseStateLabel: "待機中",
  statusTitle: "単語翻訳は開始待ちです",
  statusText: "開始条件を満たすまで待機しています。",
  progressLabel: "0 / 180",
  progressPercent: 0,
  progressDetail: "待機中",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "準備完了後に開始できます。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: true,
      blockedReason: "実行中ではありません。",
      tone: "default"
    }
  ]
} satisfies Partial<TermTranslationPhaseScreenViewModel>

const termReady = {
  phaseStateLabel: "準備完了",
  statusTitle: "単語翻訳を開始できます",
  statusText: "単語翻訳用の AI 設定を確認してから開始します。",
  progressLabel: "0 / 180",
  progressPercent: 0,
  progressDetail: "準備完了"
} satisfies Partial<TermTranslationPhaseScreenViewModel>

const termRunning = {
  viewState: "running",
  phaseStateLabel: "実行中",
  statusTitle: "単語翻訳を実行中です",
  statusText: "対象語を AI 翻訳しています。",
  progressLabel: "74 / 180",
  progressPercent: 41,
  progressDetail: "74 件を処理しました。",
  confirmedCountLabel: "74",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "すでに実行中です。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: false,
      blockedReason: "",
      tone: "warning"
    }
  ]
} satisfies Partial<TermTranslationPhaseScreenViewModel>

const termCompleted = {
  viewState: "completed",
  phaseStateLabel: "完了",
  statusTitle: "単語翻訳が完了しました",
  statusText: "次の翻訳段階へ進めます。",
  progressLabel: "180 / 180",
  progressPercent: 100,
  progressDetail: "完了",
  confirmedCountLabel: "180",
  nextPhaseStatusLabel: "NPC ペルソナ生成へ進めます",
  nextPhaseBlockedReason: "",
  actionCards: [
    {
      id: "retry",
      label: "再実行",
      disabled: false,
      blockedReason: "",
      tone: "default"
    }
  ]
} satisfies Partial<TermTranslationPhaseScreenViewModel>

const personaWaiting = {
  phaseStateLabel: "待機中",
  statusTitle: "NPC ペルソナ生成は開始待ちです",
  statusText: "単語翻訳の完了を待っています。",
  progressLabel: "0 / 96",
  progressPercent: 0,
  progressDetail: "待機中",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "単語翻訳の完了後に開始できます。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: true,
      blockedReason: "実行中ではありません。",
      tone: "default"
    }
  ]
} satisfies Partial<PersonaGenerationPhaseScreenViewModel>

const personaReady = {
  phaseStateLabel: "準備完了",
  statusTitle: "NPC ペルソナ生成を開始できます",
  statusText: "NPC ペルソナ生成用の AI 設定を確認してから開始します。",
  progressLabel: "0 / 96",
  progressPercent: 0,
  progressDetail: "準備完了"
} satisfies Partial<PersonaGenerationPhaseScreenViewModel>

const personaRunning = {
  viewState: "running",
  phaseStateLabel: "実行中",
  statusTitle: "NPC ペルソナ生成を実行中です",
  statusText: "NPC ごとのペルソナを生成しています。",
  progressLabel: "52 / 96",
  progressPercent: 54,
  progressDetail: "52 件を処理しました。",
  generatedCountLabel: "52",
  missingCountLabel: "44",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "すでに実行中です。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: false,
      blockedReason: "",
      tone: "warning"
    }
  ],
  screenActionEnablement: {
    canRefresh: true,
    canStart: false,
    canPause: true,
    canResume: false,
    canRetry: false,
    canCancel: true,
    canCheckBodyReadiness: false,
    canStartBodyPhase: false
  }
} satisfies Partial<PersonaGenerationPhaseScreenViewModel>

const personaCompleted = {
  viewState: "completed",
  phaseStateLabel: "完了",
  statusTitle: "NPC ペルソナ生成が完了しました",
  statusText: "本文翻訳の準備確認へ進めます。",
  progressLabel: "96 / 96",
  progressPercent: 100,
  progressDetail: "完了",
  generatedCountLabel: "96",
  missingCountLabel: "0",
  bodyReadinessLabel: "本文翻訳へ進めます",
  bodyReadinessBlockedReason: "",
  actionCards: [
    {
      id: "check-body-readiness",
      label: "本文翻訳を確認",
      disabled: false,
      blockedReason: "",
      tone: "primary"
    }
  ],
  screenActionEnablement: {
    canRefresh: true,
    canStart: false,
    canPause: false,
    canResume: false,
    canRetry: false,
    canCancel: false,
    canCheckBodyReadiness: true,
    canStartBodyPhase: true
  }
} satisfies Partial<PersonaGenerationPhaseScreenViewModel>

const bodyWaiting = {
  viewState: "not_ready",
  phaseStateLabel: "待機中",
  statusTitle: "本文翻訳は開始待ちです",
  statusText: "NPC ペルソナ生成の完了を待っています。",
  progressLabel: "0 / 24",
  progressPercent: 0,
  progressDetail: "待機中",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "NPC ペルソナ生成の完了後に開始できます。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: true,
      blockedReason: "実行中ではありません。",
      tone: "default"
    }
  ]
} satisfies Partial<BodyTranslationPhaseScreenViewModel>

const bodyReady = {
  phaseStateLabel: "準備完了",
  statusTitle: "本文翻訳を開始できます",
  statusText: "本文翻訳用の AI 設定を確認してから開始します。",
  progressLabel: "0 / 24",
  progressPercent: 0,
  progressDetail: "準備完了"
} satisfies Partial<BodyTranslationPhaseScreenViewModel>

const bodyRunning = {
  viewState: "running",
  phaseStateLabel: "実行中",
  statusTitle: "本文翻訳を実行中です",
  statusText: "本文フィールドの訳文を生成しています。",
  progressLabel: "11 / 24",
  progressPercent: 46,
  progressDetail: "11 件を処理しました。",
  processedCountLabel: "11",
  translatedCountLabel: "11",
  outputCountLabel: "11",
  actionCards: [
    {
      id: "start",
      label: "開始",
      disabled: true,
      blockedReason: "すでに実行中です。",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: false,
      blockedReason: "",
      tone: "warning"
    }
  ],
  screenActionEnablement: {
    canRefresh: true,
    canStart: false,
    canPause: true,
    canResume: false,
    canRetry: false,
    canCancel: true,
    canCheckOutputReadiness: false
  }
} satisfies Partial<BodyTranslationPhaseScreenViewModel>

const bodyCompleted = {
  viewState: "completed",
  phaseStateLabel: "完了",
  statusTitle: "本文翻訳が完了しました",
  statusText: "翻訳結果確認へ進めます。",
  progressLabel: "24 / 24",
  progressPercent: 100,
  progressDetail: "完了",
  processedCountLabel: "24",
  translatedCountLabel: "24",
  outputCountLabel: "24",
  resultOutputCountLabel: "24",
  outputReadinessLabel: "翻訳結果確認へ進めます",
  outputReadinessBlockedReason: "",
  outputReadinessStatusLabel: "完了",
  outputReadinessCompletedFieldCountLabel: "24",
  actionCards: [
    {
      id: "check-output-readiness",
      label: "結果確認",
      disabled: false,
      blockedReason: "",
      tone: "primary"
    }
  ],
  screenActionEnablement: {
    canRefresh: true,
    canStart: false,
    canPause: false,
    canResume: false,
    canRetry: false,
    canCancel: false,
    canCheckOutputReadiness: true
  }
} satisfies Partial<BodyTranslationPhaseScreenViewModel>

function createTarget(
  phase: TranslationJobManagementCurrentPhase,
  currentPhaseLabel: string,
  jobId: number,
  stateLabel: string,
  progressLabel: string
): TranslationJobManagementJobRunTarget {
  return {
    ...jobRunTargetSummaryFixtures.termPhase.target,
    jobId,
    currentPhase: phase,
    currentPhaseLabel,
    progressLabel,
    stateLabel,
    stateDescription: `${currentPhaseLabel} は ${stateLabel} です。`
  }
}

function createJobRunStory(
  selectedJobTarget: TranslationJobManagementJobRunTarget,
  termOverride: Partial<TermTranslationPhaseScreenViewModel> = {},
  personaOverride: Partial<PersonaGenerationPhaseScreenViewModel> = {},
  bodyOverride: Partial<BodyTranslationPhaseScreenViewModel> = {}
): Story {
  return {
    args: {
      selectedJobTarget,
      createController: createTermTranslationPhasePageControllerFixture(
        termOverride
      ),
      createPersonaController:
        createPersonaGenerationPhasePageControllerFixture(personaOverride),
      createBodyController:
        createBodyTranslationPhasePageControllerFixture(bodyOverride)
    }
  }
}

export const TermTranslationWaiting: Story = createJobRunStory(
  createTarget("term_translation", "単語翻訳", 105, "待機中", "0 / 180"),
  termWaiting
)
TermTranslationWaiting.name = "単語翻訳 / 待機中"

export const TermTranslationReady: Story = createJobRunStory(
  createTarget("term_translation", "単語翻訳", 105, "準備完了", "0 / 180"),
  termReady
)
TermTranslationReady.name = "単語翻訳 / 準備完了"

export const TermTranslationRunning: Story = createJobRunStory(
  createTarget("term_translation", "単語翻訳", 105, "実行中", "74 / 180"),
  termRunning
)
TermTranslationRunning.name = "単語翻訳 / 実行中"

export const TermTranslationCompleted: Story = createJobRunStory(
  createTarget("term_translation", "単語翻訳", 105, "完了", "180 / 180"),
  termCompleted
)
TermTranslationCompleted.name = "単語翻訳 / 完了"

export const PersonaGeneration: Story = {
  args: {
    selectedJobTarget: jobRunTargetSummaryFixtures.personaPhase.target
  }
}

export const PersonaGenerationWaiting: Story = createJobRunStory(
  createTarget("persona_generation", "NPC ペルソナ生成", 106, "待機中", "0 / 96"),
  {},
  personaWaiting
)
PersonaGenerationWaiting.name = "NPC ペルソナ生成 / 待機中"

export const PersonaGenerationReady: Story = createJobRunStory(
  createTarget(
    "persona_generation",
    "NPC ペルソナ生成",
    106,
    "準備完了",
    "0 / 96"
  ),
  {},
  personaReady
)
PersonaGenerationReady.name = "NPC ペルソナ生成 / 準備完了"

export const PersonaGenerationRunning: Story = createJobRunStory(
  createTarget("persona_generation", "NPC ペルソナ生成", 106, "実行中", "52 / 96"),
  {},
  personaRunning
)
PersonaGenerationRunning.name = "NPC ペルソナ生成 / 実行中"

export const PersonaGenerationCompleted: Story = createJobRunStory(
  createTarget("persona_generation", "NPC ペルソナ生成", 106, "完了", "96 / 96"),
  {},
  personaCompleted
)
PersonaGenerationCompleted.name = "NPC ペルソナ生成 / 完了"

export const BodyTranslation: Story = {
  args: {
    selectedJobTarget: {
      ...jobRunTargetSummaryFixtures.termPhase.target,
      jobId: 108,
      currentPhase: "body_translation",
      currentPhaseLabel: "本文翻訳",
      progressLabel: "78%"
    }
  }
}

export const BodyTranslationWaiting: Story = createJobRunStory(
  createTarget("body_translation", "本文翻訳", 108, "待機中", "0 / 24"),
  {},
  {},
  bodyWaiting
)
BodyTranslationWaiting.name = "本文翻訳 / 待機中"

export const BodyTranslationReady: Story = createJobRunStory(
  createTarget("body_translation", "本文翻訳", 108, "準備完了", "0 / 24"),
  {},
  {},
  bodyReady
)
BodyTranslationReady.name = "本文翻訳 / 準備完了"

export const BodyTranslationRunning: Story = createJobRunStory(
  createTarget("body_translation", "本文翻訳", 108, "実行中", "11 / 24"),
  {},
  {},
  bodyRunning
)
BodyTranslationRunning.name = "本文翻訳 / 実行中"

export const BodyTranslationCompleted: Story = createJobRunStory(
  createTarget("body_translation", "本文翻訳", 108, "完了", "24 / 24"),
  {},
  {},
  bodyCompleted
)
BodyTranslationCompleted.name = "本文翻訳 / 完了"
