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
