import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"

import type {
  JobRunTargetSummaryProps,
  JobUnselectedGuidanceProps,
  PhaseNavigationFooterProps
} from "../job-run-shell-props"

const ignoreAction = (): void => {}

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
