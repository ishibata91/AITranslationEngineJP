import type { ComponentProps } from "svelte"
import type { TranslationManagementViewContract } from "@ui/stores/shell-state"

import JobCard from "../JobCard.svelte"
import JobListPanel from "../JobListPanel.svelte"
import TranslationJobManagementDeleteModal from "../TranslationJobManagementDeleteModal.svelte"
import TranslationManagementStepper from "../TranslationManagementStepper.svelte"

type JobCardProps = ComponentProps<typeof JobCard>
type JobListPanelProps = ComponentProps<typeof JobListPanel>
type DeleteModalProps = ComponentProps<
  typeof TranslationJobManagementDeleteModal
>
type StepperProps = ComponentProps<typeof TranslationManagementStepper>

const noop = (): void => {}
const noopAsync = async (): Promise<void> => {}

const enabledStopOperation = {
  kind: "stop",
  label: "停止",
  helperText: "実行中のジョブを停止します。",
  enabled: true,
  reasonText: "",
  busy: false
} as const

const enabledResumeOperation = {
  kind: "resume",
  label: "再開",
  helperText: "停止中のジョブを再開します。",
  enabled: true,
  reasonText: "",
  busy: false
} as const

const enabledDeleteOperation = {
  kind: "delete",
  label: "削除",
  helperText: "ジョブを削除します。",
  enabled: true,
  reasonText: "",
  busy: false
} as const

const disabledStopOperation = {
  ...enabledStopOperation,
  enabled: false,
  reasonText: "実行中のジョブだけ停止できます。"
}

const disabledResumeOperation = {
  ...enabledResumeOperation,
  enabled: false,
  reasonText: "停止中または再開可能な失敗状態だけ再開できます。"
}

const disabledDeleteOperation = {
  ...enabledDeleteOperation,
  enabled: false,
  reasonText: "実行中のジョブは削除できません。"
}

export const selectedRunningJob = {
  jobId: 4101,
  title: "SampleQuestPlugin.esp / 会話本文翻訳",
  jobState: "Running",
  stateLabel: "実行中",
  stateDescription: "本文翻訳を実行しています。",
  stateTone: "info",
  inputSourceLabel: "SampleQuestPlugin.esp",
  sourcePath: "/synthetic/skyrim/mods/SampleQuestPlugin.esp",
  currentPhase: "body_translation",
  currentPhaseLabel: "本文翻訳",
  progressLabel: "42%",
  lastUpdatedLabel: "2026-05-18 12:30 更新",
  canOpenPhase: true,
  openBlockedReason: null,
  openBlockedReasonText: "",
  jobRunTarget: {
    jobId: 4101,
    stateLabel: "実行中",
    stateDescription: "本文翻訳を実行しています。",
    currentPhase: "body_translation",
    currentPhaseLabel: "本文翻訳",
    progressLabel: "42%",
    inputSourceLabel: "SampleQuestPlugin.esp",
    sourcePath: "/synthetic/skyrim/mods/SampleQuestPlugin.esp"
  },
  isSelected: true,
  stopOperation: enabledStopOperation,
  resumeOperation: disabledResumeOperation,
  deleteOperation: disabledDeleteOperation
} satisfies JobCardProps["job"]

export const pausedJob = {
  jobId: 4102,
  title: "LongNamedPluginForStorybookReviewWithOverflowCheck.esp / 単語翻訳",
  jobState: "Paused",
  stateLabel: "停止中",
  stateDescription: "再開待ちです。",
  stateTone: "warning",
  inputSourceLabel: "LongNamedPluginForStorybookReviewWithOverflowCheck.esp",
  sourcePath:
    "/synthetic/skyrim/mods/very/long/path/LongNamedPluginForStorybookReviewWithOverflowCheck.esp",
  currentPhase: "term_translation",
  currentPhaseLabel: "単語翻訳",
  progressLabel: "18%",
  lastUpdatedLabel: "2026-05-18 11:40 更新",
  canOpenPhase: false,
  openBlockedReason: {
    category: "stop_requested",
    title: "再開が必要",
    detail: "現在の翻訳段階へ進む前に再開してください。"
  },
  openBlockedReasonText: "現在の翻訳段階へ進む前に再開してください。",
  jobRunTarget: null,
  isSelected: false,
  stopOperation: disabledStopOperation,
  resumeOperation: enabledResumeOperation,
  deleteOperation: enabledDeleteOperation
} satisfies JobCardProps["job"]

export const recoverableFailedJob = {
  ...pausedJob,
  jobId: 4103,
  title: "RecoverableFailurePlugin.esp / NPC ペルソナ生成",
  jobState: "RecoverableFailed",
  stateLabel: "再開可能な失敗",
  stateDescription: "前回の実行で失敗しました。",
  stateTone: "danger",
  currentPhase: "persona_generation",
  currentPhaseLabel: "NPC ペルソナ生成",
  progressLabel: "64%",
  lastUpdatedLabel: "2026-05-18 10:10 更新",
  openBlockedReasonText: "失敗内容を確認してから再開してください。"
} satisfies JobCardProps["job"]

export const jobCardFixtures = {
  selectedRunning: {
    job: selectedRunningJob,
    onOpenJob: noopAsync,
    onOperation: noopAsync
  },
  pausedWithDisabledReason: {
    job: pausedJob,
    onOpenJob: noopAsync,
    onOperation: noopAsync
  },
  recoverableFailed: {
    job: recoverableFailedJob,
    onOpenJob: noopAsync,
    onOperation: noopAsync
  }
} satisfies Record<string, JobCardProps>

const baseFilterChips = [
  { id: "all", label: "すべて", count: 3, selected: true },
  { id: "Running", label: "実行中", count: 1, selected: false },
  { id: "Paused", label: "停止中", count: 1, selected: false },
  {
    id: "RecoverableFailed",
    label: "再開可能な失敗",
    count: 1,
    selected: false
  }
] satisfies JobListPanelProps["filterChips"]

const baseJobListProps = {
  phase: "ready",
  jobs: [selectedRunningJob, pausedJob, recoverableFailedJob],
  searchQuery: "",
  filterChips: baseFilterChips,
  listEmptyTitle: "未完了ジョブはありません",
  listEmptyDescription: "新規翻訳を開始すると、この一覧に表示されます。",
  listErrorTitle: "一覧を取得できませんでした",
  listErrorDescription: "時間を置いて再読み込みしてください。",
  onSearchQueryChange: noop,
  onFilterChange: noop,
  onOpenJob: noopAsync,
  onOperation: noopAsync
} satisfies JobListPanelProps

export const jobListPanelFixtures = {
  ready: baseJobListProps,
  loading: {
    ...baseJobListProps,
    phase: "loading",
    jobs: []
  },
  empty: {
    ...baseJobListProps,
    jobs: []
  },
  error: {
    ...baseJobListProps,
    phase: "error",
    jobs: []
  }
} satisfies Record<string, JobListPanelProps>

export const deleteModalFixtures = {
  closed: {
    confirmation: null,
    onClose: noop,
    onConfirm: noop
  },
  open: {
    confirmation: {
      title: "ジョブ #4102 を削除しますか",
      lines: [
        "対象: SampleQuestPlugin.esp",
        "翻訳途中の状態とキャッシュを削除します。",
        "入力データとマスター辞書は残ります。"
      ],
      confirmLabel: "削除する",
      cancelLabel: "キャンセル",
      busy: false
    },
    onClose: noop,
    onConfirm: noop
  },
  deleting: {
    confirmation: {
      title: "ジョブ #4102 を削除しています",
      lines: [
        "対象: SampleQuestPlugin.esp",
        "削除が完了するまで閉じられません。"
      ],
      confirmLabel: "削除する",
      cancelLabel: "キャンセル",
      busy: true
    },
    onClose: noop,
    onConfirm: noop
  },
  deleteFailed: {
    confirmation: {
      title: "ジョブ #4102 を削除しますか",
      lines: [
        "対象: SampleQuestPlugin.esp",
        "前回の削除は完了できませんでした。",
        "状態を確認してから再実行してください。"
      ],
      confirmLabel: "削除を再実行",
      cancelLabel: "キャンセル",
      busy: false
    },
    onClose: noop,
    onConfirm: noop
  }
} satisfies Record<string, DeleteModalProps>

export const translationManagementViews = [
  {
    id: "job-management",
    label: "未完了のジョブ",
    description: "新しい翻訳を始めるか、途中のジョブを再開します。",
    stepNumber: 1,
    directNavigation: true
  },
  {
    id: "input-review",
    label: "入力データの確認",
    description: "翻訳に使う入力データを選び、登録結果を確認します。",
    stepNumber: 2,
    directNavigation: false
  },
  {
    id: "term-translation",
    label: "単語翻訳",
    description: "選択したジョブで、単語翻訳を実行します。",
    stepNumber: 3,
    directNavigation: false
  },
  {
    id: "persona-generation",
    label: "NPC ペルソナ生成",
    description: "単語翻訳の完了後に、NPC の話し方や役割を整理します。",
    stepNumber: 4,
    directNavigation: false
  },
  {
    id: "body-translation",
    label: "本文翻訳",
    description: "NPC ペルソナを参照できる状態で、本文の翻訳を実行します。",
    stepNumber: 5,
    directNavigation: false
  }
] satisfies TranslationManagementViewContract[]

export const stepperFixtures = {
  jobManagementCurrent: {
    currentViewId: "job-management",
    views: translationManagementViews,
    onSelect: noop
  },
  termTranslationCurrent: {
    currentViewId: "term-translation",
    views: translationManagementViews,
    onSelect: noop
  }
} satisfies Record<string, StepperProps>
