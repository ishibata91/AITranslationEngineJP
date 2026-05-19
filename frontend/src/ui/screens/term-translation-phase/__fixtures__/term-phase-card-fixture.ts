import type { ComponentProps } from "svelte"
import TermExecutionSettingsCard from "../TermExecutionSettingsCard.svelte"
import TermResultSummaryCard from "../TermResultSummaryCard.svelte"

type TermExecutionSettingsCardProps = ComponentProps<
  typeof TermExecutionSettingsCard
>
type TermResultSummaryCardProps = ComponentProps<typeof TermResultSummaryCard>

export const termExecutionSettingsCardFixture: TermExecutionSettingsCardProps = {
  providerSkippedLabel: "provider 実行あり",
  details: [
    { label: "provider", value: "Sample Provider" },
    { label: "model", value: "sample-term-model" },
    { label: "execution mode", value: "batch" },
    { label: "credential reference", value: "credential-ref:term-story" },
    { label: "snapshot", value: "term-snapshot-2026-05-18" }
  ]
}

export const termResultSummaryCardFixture: TermResultSummaryCardProps = {
  nextPhaseStatusLabel: "次の翻訳段階へ進めます",
  details: [
    { label: "確定訳語件数", value: "180" },
    { label: "ジョブ内辞書反映件数", value: "42" },
    { label: "置換対象件数", value: "168" },
    { label: "未一致件数", value: "12" },
    {
      label: "次の翻訳段階",
      value: "次の翻訳段階へ進めます",
      note: "ブロック理由はありません。"
    }
  ]
}
