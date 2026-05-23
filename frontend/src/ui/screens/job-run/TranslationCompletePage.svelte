<script lang="ts">
  import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"
  import ProcessingTargetListPanel from "@ui/components/ProcessingTargetListPanel.svelte"
  import type { ProcessingTargetListItem } from "@ui/components/processing-target-list-panel-types"

  import TranslationCompleteSummaryPanel from "./TranslationCompleteSummaryPanel.svelte"

  interface Props {
    jobId: number
    rows: BodyTranslationFieldResultItem[]
  }

  let { jobId, rows }: Props = $props()

  function buildCompletedTargetItems(
    resultRows: BodyTranslationFieldResultItem[]
  ): ProcessingTargetListItem[] {
    return resultRows.map((row) => {
      const sourceExcerpt = row.sourceExcerpt || "原文なし"
      const translatedText = row.translatedText || "訳文なし"

      return {
        id: row.fieldId,
        name: `${row.recordTypeLabel} ${row.fieldLabel}`,
        titleParts: [
          { text: `原文: ${sourceExcerpt}` },
          { text: `訳文: ${translatedText}` }
        ],
        detail: "本文翻訳で保持された訳文として出力管理へ進む前に確認する訳文。",
        metadata: [
          { label: "翻訳項目 ID", value: row.fieldId },
          { label: "record type", value: row.recordTypeLabel },
          { label: "field", value: row.fieldLabel },
          { label: "field type", value: row.fieldTypeLabel },
          { label: "FormID", value: row.formIdLabel },
          { label: "EditorID", value: row.editorIdLabel },
          { label: "出力状態", value: row.outputStatus },
          { label: "保護検証", value: row.protectionValidation },
          { label: "retry", value: row.retryCountLabel }
        ]
      }
    })
  }

  const completedTargetItems = $derived(buildCompletedTargetItems(rows))
</script>

<section
  class="complete-shell"
  data-testid="translation-complete-translation-complete-screen"
  id="translationCompleteView"
>
  <TranslationCompleteSummaryPanel {jobId} />
  <ProcessingTargetListPanel items={completedTargetItems} />
</section>

<style>
  .complete-shell {
    display: grid;
    gap: 1rem;
  }
</style>
