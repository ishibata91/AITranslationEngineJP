import type { ComponentProps } from "svelte"
import PhaseProgressPanel from "../PhaseProgressPanel.svelte"
import PhaseStatusPanel from "../PhaseStatusPanel.svelte"

type PhaseStatusPanelProps = ComponentProps<typeof PhaseStatusPanel>
type PhaseProgressPanelProps = ComponentProps<typeof PhaseProgressPanel>

export const phaseStatusPanelFixtures = {
  personaAiSettingsReady: {
    eyebrow: "translation-management",
    title: "NPC ペルソナ生成",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、リトライ、キャンセルを判断します。",
    state: "not_started",
    stateLabel: "未開始",
    statusTitle: "NPC ペルソナ生成を開始できます",
    statusText: "NPC ペルソナ生成用の AI 設定を確認してから開始します。",
    errorMessage: "",
    testId: "phase-status-persona-ai-settings-ready",
    metrics: [
      { label: "対象", value: "96" },
      { label: "処理済み", value: "0" },
      { label: "成功", value: "0" },
      { label: "失敗", value: "0" },
      { label: "スキップ", value: "8" }
    ]
  },
  termAiSettingsReady: {
    eyebrow: "translation-management",
    title: "単語翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、リトライを判断します。",
    state: "idle_ready",
    stateLabel: "未開始",
    statusTitle: "単語翻訳を開始できます",
    statusText: "単語翻訳用の AI 設定を確認してから開始します。",
    errorMessage: "",
    testId: "phase-status-term-ai-settings-ready",
    metrics: [
      { label: "対象", value: "180" },
      { label: "処理済み", value: "0" },
      { label: "成功", value: "0" },
      { label: "失敗", value: "0" },
      { label: "スキップ", value: "0" }
    ]
  },
  bodyAiSettingsReady: {
    eyebrow: "translation-management",
    title: "本文翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、回復を判断します。",
    state: "ready",
    stateLabel: "未開始",
    statusTitle: "本文翻訳を開始できます",
    statusText: "本文翻訳用の AI 設定を確認してから開始します。",
    errorMessage: "",
    testId: "phase-status-body-ai-settings-ready",
    metrics: [
      { label: "対象", value: "24" },
      { label: "処理済み", value: "0" },
      { label: "成功", value: "0" },
      { label: "失敗", value: "0" },
      { label: "スキップ", value: "2" }
    ]
  },
  ready: {
    eyebrow: "translation-management",
    title: "単語翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "idle_ready",
    stateLabel: "未開始",
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
    stateLabel: "実行中",
    statusTitle: "実行中",
    statusText: "対象データを処理しています。",
    errorMessage: "",
    testId: "phase-status-running",
    metrics: [
      { label: "対象", value: "120" },
      { label: "処理済み", value: "48" },
      { label: "成功", value: "48" },
      { label: "失敗", value: "1" },
      { label: "スキップ", value: "8" }
    ]
  },
  recoverableFailed: {
    eyebrow: "translation-management",
    title: "本文翻訳",
    gatewayStatus: "接続済み",
    lead: "現在の翻訳段階、進行状況、結果を確認します。",
    state: "recoverable_failed",
    stateLabel: "再試行可能",
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
    stateLabel: "完了",
    statusTitle: "完了しました",
    statusText: "次の翻訳段階へ進めます。",
    errorMessage: "",
    testId: "phase-status-completed"
  }
} satisfies Record<string, PhaseStatusPanelProps>

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
    { label: "対象件数", value: "120" },
    { label: "成功件数", value: "58" }
  ]
}
