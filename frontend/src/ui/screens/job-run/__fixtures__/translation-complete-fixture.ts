import type { ComponentProps } from "svelte"
import TranslationCompleteSummaryPanel from "../TranslationCompleteSummaryPanel.svelte"
import TranslationResultListPanel from "../TranslationResultListPanel.svelte"

type TranslationCompleteSummaryPanelProps = ComponentProps<
  typeof TranslationCompleteSummaryPanel
>
type TranslationResultListPanelProps = ComponentProps<
  typeof TranslationResultListPanel
>

export const translationCompleteSummaryPanelFixture: TranslationCompleteSummaryPanelProps =
  {
    jobId: 1001
  }

export const translationResultListPanelFixture: TranslationResultListPanelProps =
  {
    rows: [
      {
        fieldId: "field-name-001",
        fieldLabel: "Name",
        recordTypeLabel: "NPC_",
        fieldTypeLabel: "FULL",
        formIdLabel: "0x00001234",
        editorIdLabel: "SampleActorOne",
        sourceExcerpt: "Sample source line for completion review.",
        translatedText: "完了画面確認用のサンプル訳文です。",
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
          "A longer synthetic source line for pagination and wrapping review.",
        translatedText:
          "ページングと折り返しを確認するための長めの合成訳文です。",
        outputStatus: "ready",
        protectionValidation: "passed",
        retryCountLabel: "1",
        rawItem: null
      }
    ],
    pageIndex: 0,
    pageCount: 2,
    pageLabel: "1 / 2",
    onPreviousPage: () => undefined,
    onNextPage: () => undefined
  }

export const emptyTranslationResultListPanelFixture: TranslationResultListPanelProps =
  {
    ...translationResultListPanelFixture,
    rows: [],
    pageCount: 1,
    pageLabel: "1 / 1"
  }
