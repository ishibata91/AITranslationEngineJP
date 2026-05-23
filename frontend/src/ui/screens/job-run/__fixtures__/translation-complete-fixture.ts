import type { ComponentProps } from "svelte"
import TranslationCompleteSummaryPanel from "../TranslationCompleteSummaryPanel.svelte"

type TranslationCompleteSummaryPanelProps = ComponentProps<
  typeof TranslationCompleteSummaryPanel
>

export const translationCompleteSummaryPanelFixture: TranslationCompleteSummaryPanelProps =
  {
    jobId: 1001
  }
