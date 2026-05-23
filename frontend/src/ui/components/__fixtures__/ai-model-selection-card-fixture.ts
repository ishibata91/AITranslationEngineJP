import type { ComponentProps } from "svelte"
import AIModelSelectionCard from "../AIModelSelectionCard.svelte"

type AIModelSelectionCardProps = ComponentProps<typeof AIModelSelectionCard>

const ignoreSelectionChange = (): void => {}

export const aiModelSelectionCardFixture: AIModelSelectionCardProps = {
  eyebrow: "Storybook sample",
  title: "AI モデル選択",
  titleId: "storybook-ai-model-selection-card-title",
  helperText: "固定 props だけで表示する Storybook 最小 fixture。",
  statusLabel: "設定済み",
  statusTone: "success",
  providerSelectId: "storybook-provider-select",
  providerValue: "sample-provider",
  providerOptions: [{ value: "sample-provider", label: "Sample Provider" }],
  providerDisabled: true,
  onProviderChange: ignoreSelectionChange,
  credentialStatusLabel: "Storybook 用の表示値です。",
  credentialStatusTone: "success",
  showCredentialStatus: true,
  modelSelectId: "storybook-model-select",
  modelValue: "sample-model",
  modelOptions: [{ modelId: "sample-model", label: "Sample Model" }],
  modelDisabled: true,
  onModelChange: ignoreSelectionChange,
  modelStatusText: "backend へ接続しない固定候補です。",
  footerMessage: "fixture は component 横の __fixtures__ に置きます。",
  actionButtonLabel: "確認",
  actionButtonDisabled: true
}

export const aiModelSelectionCardStates = {
  success: aiModelSelectionCardFixture,
  loading: {
    ...aiModelSelectionCardFixture,
    statusLabel: "モデル一覧更新中",
    statusTone: "neutral",
    modelStatusText: "モデル一覧を更新しています。",
    refreshSpinning: true,
    refreshDisabled: true
  },
  failed: {
    ...aiModelSelectionCardFixture,
    statusLabel: "モデル一覧取得失敗",
    statusTone: "warning",
    modelStatusText: "モデル一覧を取得できませんでした。",
    footerWarningText: "時間を置いてモデル一覧を更新してください。"
  },
  credentialMissing: {
    ...aiModelSelectionCardFixture,
    statusLabel: "設定未完了",
    statusTone: "warning",
    credentialStatusLabel: "認証情報が不足しています。",
    credentialStatusTone: "warning",
    showCredentialWarning: true,
    credentialWarningText: "AI サービス設定で認証状態を確認してください。",
    footerWarningText: "認証不足のため開始できません。"
  },
  runningLocked: {
    ...aiModelSelectionCardFixture,
    statusLabel: "固定済み",
    providerDisabled: true,
    modelDisabled: true,
    refreshDisabled: true,
    footerMessage: "実行中は AI 設定を編集できません。"
  }
} satisfies Record<string, AIModelSelectionCardProps>
