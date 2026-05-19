import type { ComponentProps } from "svelte"
import BodyExecutionSummaryCard from "../BodyExecutionSummaryCard.svelte"
import BodyInputSummaryCard from "../BodyInputSummaryCard.svelte"
import BodyResultSummaryCard from "../BodyResultSummaryCard.svelte"
import FieldResultListPanel from "../FieldResultListPanel.svelte"
import OutputReadinessCard from "../OutputReadinessCard.svelte"

type BodyInputSummaryCardProps = ComponentProps<typeof BodyInputSummaryCard>
type BodyExecutionSummaryCardProps = ComponentProps<
  typeof BodyExecutionSummaryCard
>
type BodyResultSummaryCardProps = ComponentProps<typeof BodyResultSummaryCard>
type FieldResultListPanelProps = ComponentProps<typeof FieldResultListPanel>
type OutputReadinessCardProps = ComponentProps<typeof OutputReadinessCard>

export const bodyInputSummaryCardFixture: BodyInputSummaryCardProps = {
  readinessReason: "出力準備は本文翻訳結果の整合確認後に可能です。",
  details: [
    { label: "dictionary digest", value: "dictionary-digest:sample-body" },
    { label: "persona digest", value: "persona-digest:sample-body" },
    { label: "metadata digest", value: "metadata-digest:sample-body" },
    { label: "prompt digest", value: "prompt-digest:sample-body" },
    { label: "input snapshot", value: "input-snapshot:body-story" },
    { label: "skipped reasons", value: "完全一致辞書により 2 件を除外" }
  ]
}

export const bodyExecutionSummaryCardFixture: BodyExecutionSummaryCardProps = {
  providerStateLabel: "provider 実行完了",
  details: [
    { label: "provider", value: "Sample Provider" },
    { label: "model", value: "sample-body-model" },
    { label: "execution mode", value: "batch" },
    { label: "credential ref", value: "credential-ref:body-story" },
    { label: "request unit count", value: "8" },
    { label: "provider target count", value: "24" },
    { label: "exact dictionary excluded", value: "2" },
    { label: "partial dictionary constrained", value: "5" },
    { label: "output count", value: "22" },
    { label: "late response rejected", value: "0" }
  ]
}

export const bodyResultSummaryCardFixture: BodyResultSummaryCardProps = {
  outputReadinessLabel: "出力準備完了",
  details: [
    { label: "translated count", value: "22" },
    { label: "failed count", value: "0" },
    { label: "skipped count", value: "2" },
    { label: "output ready count", value: "24" },
    { label: "result output count", value: "22" },
    { label: "status consistency", value: "整合しています" }
  ]
}

export const fieldResultListPanelFixture: FieldResultListPanelProps = {
  availabilityLabel: "2 件を表示",
  items: [
    {
      fieldId: "field-name-001",
      fieldLabel: "Name",
      recordTypeLabel: "NPC_",
      fieldTypeLabel: "FULL",
      formIdLabel: "0x00001234",
      editorIdLabel: "SampleActorOne",
      sourceExcerpt: "Sample source line for review.",
      translatedText: "確認用のサンプル訳文です。",
      outputStatus: "ready",
      protectionValidation: "passed",
      retryCountLabel: "0",
      rawItem: null
    },
    {
      fieldId: "field-dialogue-002",
      fieldLabel: "Dialogue",
      recordTypeLabel: "INFO",
      fieldTypeLabel: "NAM1",
      formIdLabel: "0x00005678",
      editorIdLabel: "SampleDialogueTopic",
      sourceExcerpt:
        "A longer synthetic source line that checks wrapping inside the card.",
      translatedText: "カード内で折り返しを確認するための長めの合成訳文です。",
      outputStatus: "ready",
      protectionValidation: "passed",
      retryCountLabel: "1",
      rawItem: null
    }
  ]
}

export const emptyFieldResultListPanelFixture: FieldResultListPanelProps = {
  availabilityLabel: "0 件",
  items: []
}

export const outputReadinessCardFixture: OutputReadinessCardProps = {
  outputReadinessLabel: "出力準備完了",
  details: [
    { label: "readiness", value: "出力準備完了" },
    { label: "completed field count", value: "24" },
    { label: "status consistency", value: "整合しています" },
    { label: "blocked reason", value: "ブロック理由はありません。" }
  ]
}

export const blockedOutputReadinessCardFixture: OutputReadinessCardProps = {
  outputReadinessLabel: "出力準備未達",
  details: [
    { label: "readiness", value: "出力準備未達" },
    { label: "completed field count", value: "18 / 24" },
    { label: "status consistency", value: "未完了の field result があります" },
    { label: "blocked reason", value: "失敗項目を再試行してください。" }
  ]
}
