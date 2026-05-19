import type { ComponentProps } from "svelte"
import PhaseActionPanel from "../PhaseActionPanel.svelte"
import PhaseFailureInfoCard from "../PhaseFailureInfoCard.svelte"
import PhaseMetricCounterGrid from "../PhaseMetricCounterGrid.svelte"
import PhaseProgressPanel from "../PhaseProgressPanel.svelte"
import PhaseStatusPanel from "../PhaseStatusPanel.svelte"

type PhaseStatusPanelProps = ComponentProps<typeof PhaseStatusPanel>
type PhaseActionPanelProps = ComponentProps<typeof PhaseActionPanel>
type PhaseProgressPanelProps = ComponentProps<typeof PhaseProgressPanel>
type PhaseFailureInfoCardProps = ComponentProps<typeof PhaseFailureInfoCard>
type PhaseMetricCounterGridProps = ComponentProps<typeof PhaseMetricCounterGrid>

const noop = (): void => {}

export const phaseStatusPanelFixtures = {
  ready: {
    eyebrow: "translation-management",
    title: "単語翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "idle_ready",
    stateLabel: "Ready",
    statusTitle: "開始できます",
    statusText: "対象データの準備が完了しています。",
    errorMessage: "",
    testId: "phase-status-ready"
  },
  running: {
    eyebrow: "translation-management",
    title: "NPC ペルソナ生成",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "running",
    stateLabel: "Running",
    statusTitle: "実行中",
    statusText: "対象データを処理しています。",
    errorMessage: "",
    testId: "phase-status-running",
    metrics: [
      { label: "target", value: "120" },
      { label: "generated", value: "48" },
      { label: "failed", value: "1" },
      { label: "skipped", value: "8" }
    ]
  },
  recoverableFailed: {
    eyebrow: "translation-management",
    title: "本文翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "recoverable_failed",
    stateLabel: "RecoverableFailed",
    statusTitle: "再試行できます",
    statusText: "一部の処理が完了していません。",
    errorMessage: "一時的な失敗が発生しました。",
    testId: "phase-status-recoverable-failed"
  },
  completed: {
    eyebrow: "translation-management",
    title: "単語翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "completed",
    stateLabel: "Completed",
    statusTitle: "完了しました",
    statusText: "次の翻訳段階へ進めます。",
    errorMessage: "",
    testId: "phase-status-completed"
  }
} satisfies Record<string, PhaseStatusPanelProps>

export const phaseActionPanelFixture: PhaseActionPanelProps = {
  headingId: "phaseActionStoryHeading",
  testId: "phase-action-story",
  currentPhaseLabel: "単語翻訳",
  actions: [
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
      tone: "default"
    },
    {
      id: "retry",
      label: "リトライ",
      disabled: true,
      blockedReason: "失敗情報がありません。",
      tone: "warning"
    },
    {
      id: "cancel",
      label: "キャンセル",
      disabled: false,
      blockedReason: "",
      tone: "warning"
    }
  ],
  columns: 4,
  onAction: noop
}

export const phaseProgressPanelFixture: PhaseProgressPanelProps = {
  headingId: "phaseProgressStoryHeading",
  testId: "phase-progress-story",
  eyebrow: "翻訳段階の進行状況",
  title: "進行状況",
  progressLabel: "48%",
  progressPercent: 48,
  progressDetail: "120 件中 58 件を処理しました。",
  details: [
    { label: "開始時刻", value: "2026-05-18 10:00" },
    { label: "完了時刻", value: "未完了" },
    { label: "target count", value: "120" },
    { label: "success count", value: "58" }
  ]
}

export const phaseFailureInfoCardFixture: PhaseFailureInfoCardProps = {
  headingId: "phaseFailureStoryHeading",
  testId: "phase-failure-story",
  errorKindLabel: "temporary",
  errorReasonLabel: "一時的な失敗",
  retryableLabel: "再試行可能"
}

export const phaseMetricCounterGridFixture: PhaseMetricCounterGridProps = {
  metrics: [
    { label: "target", value: "120" },
    { label: "generated", value: "96" },
    { label: "failed", value: "2" },
    { label: "skipped", value: "22" }
  ]
}
