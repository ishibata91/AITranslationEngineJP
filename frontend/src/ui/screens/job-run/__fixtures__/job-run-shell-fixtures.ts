import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"

import type {
  JobRunTargetSummaryProps,
  JobUnselectedGuidanceProps,
  PhaseNavigationFooterProps,
  ProcessingTargetListItem
} from "../job-run-shell-props"
import type { ComponentProps } from "svelte"
import type ProcessingTargetListPanel from "@ui/components/ProcessingTargetListPanel.svelte"

const ignoreAction = (): void => {}
type ProcessingTargetListPanelProps = ComponentProps<typeof ProcessingTargetListPanel>

const baseTarget: TranslationJobManagementJobRunTarget = {
  jobId: 105,
  stateLabel: "実行中",
  stateDescription: "現在の翻訳段階を実行しています。",
  currentPhase: "term_translation",
  currentPhaseLabel: "単語翻訳",
  progressLabel: "42%",
  inputSourceLabel: "sample-input.json",
  sourcePath: "synthetic://translation-input/sample-input.json"
}

export const jobRunTargetSummaryFixtures = {
  termPhase: {
    target: baseTarget
  },
  personaPhase: {
    target: {
      ...baseTarget,
      jobId: 106,
      currentPhase: "persona_generation",
      currentPhaseLabel: "NPC ペルソナ生成",
      progressLabel: "64%"
    }
  },
  longSourcePath: {
    target: {
      ...baseTarget,
      jobId: 107,
      inputSourceLabel: "sample-input-with-long-name-for-review.json",
      sourcePath:
        "synthetic://translation-input/review/sample-input-with-long-name-for-job-run-target-summary-overflow-check.json"
    }
  }
} satisfies Record<string, JobRunTargetSummaryProps>

export const jobUnselectedGuidanceFixtures = {
  default: {
    onOpenJobManagement: ignoreAction
  }
} satisfies Record<string, JobUnselectedGuidanceProps>

function createProcessingTargetItems(
  prefix: string,
  detail: string,
  count: number
): ProcessingTargetListItem[] {
  return Array.from({ length: count }, (_, index) => {
    const itemNumber = index + 1
    return {
      id: `${prefix}-${itemNumber}`,
      name: `${prefix} ${itemNumber.toString().padStart(3, "0")}`,
      detail,
      metadata: [
        { label: "対象 ID", value: `${prefix}-${itemNumber}` },
        { label: "原文 field", value: `${prefix} 原文 ${itemNumber}` },
        { label: "処理内容", value: detail },
        { label: "確認状態", value: "未確認" }
      ]
    }
  })
}

export const processingTargetListPanelFixtures = {
  termTranslationFirstPage: {
    items: Array.from({ length: 125 }, (_, index): ProcessingTargetListItem => {
      const itemNumber = index + 1
      const paddedNumber = itemNumber.toString().padStart(3, "0")
      return {
        id: `term-translation-${paddedNumber}`,
        name: `用語候補 ${paddedNumber}`,
        titleParts: [
          { text: `原語 ${paddedNumber}: Dragonborn` },
          { text: `訳語 ${paddedNumber}: ドラゴンボーン` },
          { text: `レコード種別 ${paddedNumber}: NPC_:FULL` }
        ],
        detail: "AIサービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存する候補。",
        metadata: [
          { label: "辞書 ID", value: `term-translation-${paddedNumber}` },
          { label: "カテゴリ", value: "共通辞書対象外" },
          { label: "登録元", value: `term.source.${paddedNumber}` },
          { label: "最終更新", value: "未登録" },
          {
            label: "メモ",
            value: "AIサービスへ送る訳語候補"
          }
        ]
      }
    })
  },
  personaGeneration: {
    items: Array.from({ length: 64 }, (_, index): ProcessingTargetListItem => {
      const itemNumber = index + 1
      const paddedNumber = itemNumber.toString().padStart(3, "0")
      return {
        id: `persona-generation-${paddedNumber}`,
        name: `NPC ${paddedNumber}`,
        detail:
          "NPC の原文発話、NPC 属性、会話文脈、共通ペルソナ参照からペルソナ参照情報を作る入力。",
        metadata: [
          { label: "FormID", value: `FE${paddedNumber}001` },
          { label: "EditorID", value: `NPC_DIALOGUE_${paddedNumber}` },
          { label: "対象プラグイン", value: "SamplePersonaPlugin.esp" },
          { label: "元プラグイン", value: "Skyrim.esm" },
          { label: "声", value: "MaleNord" },
          { label: "話し方", value: "共通ペルソナ参照から補完" },
          { label: "ペルソナ本文", value: "生成前" },
          { label: "最終更新", value: "未生成" }
        ]
      }
    })
  },
  bodyTranslation: {
    items: Array.from({ length: 81 }, (_, index): ProcessingTargetListItem => {
      const itemNumber = index + 1
      const paddedNumber = itemNumber.toString().padStart(3, "0")
      return {
        id: `body-translation-${paddedNumber}`,
        name: `本文翻訳項目 ${paddedNumber}`,
        titleParts: [
          {
            text: `原文 ${paddedNumber}: The ancient stone door remains sealed.`
          },
          {
            text: `訳文 ${paddedNumber}: 古代の石扉は封印されたままです。`
          }
        ],
        detail: "辞書とペルソナを参照して訳文を作る翻訳項目。",
        metadata: [
          { label: "翻訳項目 ID", value: `body-translation-${paddedNumber}` },
          { label: "原文 field", value: `body.source.${paddedNumber}` },
          { label: "参照情報", value: "翻訳ジョブ内辞書、ペルソナ参照情報" },
          { label: "確認状態", value: "未確認" }
        ]
      }
    })
  },
  translationComplete: {
    items: Array.from({ length: 52 }, (_, index): ProcessingTargetListItem => {
      const itemNumber = index + 1
      const paddedNumber = itemNumber.toString().padStart(3, "0")
      return {
        id: `translated-text-${paddedNumber}`,
        name: `原文 ${paddedNumber} / 訳文 ${paddedNumber}`,
        titleParts: [
          {
            text: `原文 ${paddedNumber}: The ancient stone door remains sealed.`
          },
          {
            text: `訳文 ${paddedNumber}: 古代の石扉は封印されたままです。`
          }
        ],
        detail: "本文翻訳で保持された訳文として出力管理へ進む前に確認する訳文。",
        metadata: [
          { label: "訳文 ID", value: `translated-text-${paddedNumber}` },
          { label: "原文 field", value: `source.field.${paddedNumber}` },
          { label: "確認状態", value: itemNumber === 1 ? "確認済み" : "未確認" }
        ]
      }
    })
  },
  finalPage: {
    items: createProcessingTargetItems(
      "最終ページ確認候補",
      "現在ページの表示範囲だけを画面要素として確認する候補。",
      51
    ),
    initialPage: 2
  },
  valueonlyTitleParts: {
    items: Array.from({ length: 10 }, (_, index): ProcessingTargetListItem => {
      const itemNumber = index + 1
      const paddedNumber = itemNumber.toString().padStart(3, "0")
      return {
        id: `valueonly-${paddedNumber}`,
        name: `用語候補 ${paddedNumber}`,
        titleParts: [
          { text: `Dragonborn${paddedNumber}` },
          { text: `ドラゴンボーン${paddedNumber}` }
        ],
        detail: "ラベルなし値のみ形式 titleParts のレンダリング確認。",
        metadata: [
          { label: "対象 ID", value: `valueonly-${paddedNumber}` }
        ]
      }
    }),
    titleColumnLabels: ["原語", "訳語", "レコード種別"]
  },
  longText: {
    items: [
      {
        id: "long-processing-target",
        name: "長い処理対象名を持つ本文翻訳項目",
        titleParts: [
          {
            text:
              "原文: When the sealed archive beneath the ancient stone keep opens, every hidden inscription must be translated without losing the quest context."
          },
          {
            text:
              "訳文: 古代の石造りの砦の地下にある封印書庫が開いたとき、隠された碑文はクエスト文脈を失わずに翻訳される必要があります。"
          }
        ],
        detail:
          "長い row title part は省略表示し、hover または focus の tooltip で全文を確認する。",
        metadata: [
          { label: "対象 ID", value: "long-processing-target" },
          {
            label: "原文 field",
            value: "body.dialogue.original.text.with.long.field.path.for.review"
          },
          {
            label: "参照情報",
            value: "翻訳ジョブ内辞書、ペルソナ参照情報、本文翻訳対象原文"
          },
          { label: "確認状態", value: "未確認" }
        ]
      }
    ]
  }
} satisfies Record<string, ProcessingTargetListPanelProps>

export const phaseNavigationFooterFixtures = {
  termReady: {
    title: "単語翻訳の次の作業",
    titleId: "termPhaseNavigationFooter",
    description: "単語翻訳が完了し、辞書を参照できる場合だけ次へ進めます。",
    reasons: [],
    primaryDisabled: false,
    onBack: ignoreAction,
    onPrimary: ignoreAction,
    dataTestId: "job-run-next-action-footer"
  },
  termBlocked: {
    title: "単語翻訳の次の作業",
    titleId: "termPhaseNavigationFooter",
    description: "単語翻訳が完了し、辞書を参照できる場合だけ次へ進めます。",
    reasons: [
      "次へ進めません。単語翻訳の完了状況と辞書の参照状態を確認してください。"
    ],
    primaryDisabled: true,
    onBack: ignoreAction,
    onPrimary: ignoreAction,
    dataTestId: "job-run-next-action-footer"
  },
  complete: {
    title: "翻訳完了後の次の作業",
    titleId: "completeNavigationFooter",
    description: "翻訳結果を確認した後は、出力管理で出力対象を選びます。",
    reasons: [],
    primaryLabel: "出力管理へ進む",
    primaryDisabled: false,
    onBack: ignoreAction,
    onPrimary: ignoreAction,
    dataTestId: "translation-complete-post-completion-next-action"
  }
} satisfies Record<string, PhaseNavigationFooterProps>
