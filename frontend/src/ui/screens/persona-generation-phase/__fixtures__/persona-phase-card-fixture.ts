import type { ComponentProps } from "svelte"
import BodyReadinessInputCard from "../BodyReadinessInputCard.svelte"
import PersonaExecutionSettingsCard from "../PersonaExecutionSettingsCard.svelte"
import PersonaResultSummaryCard from "../PersonaResultSummaryCard.svelte"
import PersonaTargetSummaryCard from "../PersonaTargetSummaryCard.svelte"

type PersonaTargetSummaryCardProps = ComponentProps<
  typeof PersonaTargetSummaryCard
>
type PersonaExecutionSettingsCardProps = ComponentProps<
  typeof PersonaExecutionSettingsCard
>
type PersonaResultSummaryCardProps = ComponentProps<
  typeof PersonaResultSummaryCard
>
type BodyReadinessInputCardProps = ComponentProps<typeof BodyReadinessInputCard>

export const personaTargetSummaryCardFixture: PersonaTargetSummaryCardProps = {
  details: [
    { label: "NPC count", value: "96" },
    { label: "common persona hit", value: "24" },
    { label: "common persona miss", value: "72" },
    { label: "対象外理由", value: "翻訳対象外 8 件" },
    { label: "target snapshot", value: "persona-target-snapshot-2026-05-18" }
  ]
}

export const personaExecutionSettingsCardFixture: PersonaExecutionSettingsCardProps =
  {
    details: [
      { label: "provider", value: "Sample Provider" },
      { label: "model", value: "sample-persona-model" },
      { label: "execution mode", value: "batch" },
      { label: "credential ref", value: "credential-ref:persona-story" },
      { label: "input count", value: "96" },
      { label: "output count", value: "84" },
      { label: "prompt digest", value: "digest:persona-story" },
      { label: "error kind", value: "none" }
    ]
  }

export const personaResultSummaryCardFixture: PersonaResultSummaryCardProps = {
  bodyReadinessLabel: "本文翻訳を開始できます",
  details: [
    { label: "persona snapshot", value: "persona-result-snapshot-2026-05-18" },
    { label: "snapshot 参照状態", value: "参照可能" },
    { label: "persona count", value: "84" },
    { label: "missing count", value: "12" },
    {
      label: "body readiness",
      value: "本文翻訳を開始できます",
      note: "ブロック理由はありません。"
    }
  ]
}

export const bodyReadinessInputCardFixture: BodyReadinessInputCardProps = {
  details: [
    { label: "入力 summary", value: "辞書、ペルソナ、本文入力が揃っています。" },
    { label: "evidence refs", value: "term-snapshot, persona-snapshot" }
  ]
}
