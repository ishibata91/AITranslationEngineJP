import type { TranslationOutputCompletedJobSummary } from "@application/gateway-contract/translation-output-artifact"
import type {
  TranslationOutputArtifactCommandSnapshot,
  TranslationOutputDiffPreviewSnapshot,
  TranslationOutputReviewSnapshot
} from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"

const ignoreAction = (): void => {}
const ignoreString = (value: string): void => {
  void value
}
const ignoreArtifact = (artifactId: number | null): void => {
  void artifactId
}
const ignoreJob = (job: TranslationOutputCompletedJobSummary): void => {
  void job
}

const completedJobs: TranslationOutputCompletedJobSummary[] = [
  {
    jobId: 2401,
    jobStatus: "completed",
    artifactStatus: "current",
    outputReady: true,
    translatedCount: 128,
    outputStatusDistribution: {
      translated: 118,
      unchanged: 10
    }
  },
  {
    jobId: 2402,
    jobStatus: "completed",
    artifactStatus: "stale",
    outputReady: true,
    translatedCount: 42,
    outputStatusDistribution: {
      translated: 39,
      rejected: 3
    }
  }
]

const review: TranslationOutputReviewSnapshot = {
  completedJobs,
  selectedJobId: 2401,
  selectedJobStatus: "completed",
  bodyPhaseStatus: "completed",
  outputReady: true,
  translatedCount: 128,
  rowCount: 128,
  inputSnapshotDigest: "snapshot-digest-synthetic-2401",
  sourceFileDigest: "source-digest-synthetic-2401",
  readiness: true,
  retryable: true,
  artifactId: 901,
  artifactStatus: "current",
  artifactRowCount: 128,
  currentVersion: true,
  rejectionReasons: []
}

const rejectedReview: TranslationOutputReviewSnapshot = {
  ...review,
  readiness: false,
  artifactStatus: "rejected",
  currentVersion: false,
  rejectionKind: "compatibility_rejected",
  rejectionReasons: [
    {
      errorKind: "compatibility_rejected",
      reason: "xTranslator compatibility check rejected 3 rows.",
      retryable: true,
      isRedacted: false
    }
  ]
}

const diffPreview: TranslationOutputDiffPreviewSnapshot = {
  artifactId: 901,
  compatibilityPassed: true,
  compatibilityWarningCount: 1,
  compatibilityRejectCount: 0,
  rows: [
    {
      fieldId: 3001,
      rowDigest: "row-digest-synthetic-001",
      edid: "SAMPLE_WEAPON_NAME",
      rec: "WEAP",
      field: "FULL",
      formId: "000A0001",
      sourceExcerpt: "Iron Sword",
      destExcerpt: "鉄の剣",
      xTranslatorStatus: 1,
      internalOutputStatus: "translated",
      rowReflectionSummary: "dest text will be written",
      canRegenerate: false
    },
    {
      fieldId: 3002,
      rowDigest: "row-digest-synthetic-002",
      edid: "SAMPLE_ARMOR_DESC",
      rec: "ARMO",
      field: "DESC",
      formId: "000A0002",
      sourceExcerpt: "Long source excerpt for wrapping review.",
      destExcerpt:
        "Storybook 用の長い翻訳文です。出力結果の折り返しだけを確認します。",
      xTranslatorStatus: 2,
      internalOutputStatus: "stale",
      rowReflectionSummary: "source digest changed after artifact generation",
      staleReason: "source_digest_mismatch",
      canRegenerate: true
    }
  ]
}

const lastCommand: TranslationOutputArtifactCommandSnapshot = {
  jobId: 2401,
  artifactId: 901,
  artifactStatus: "generated",
  rowCount: 128,
  filePath: "synthetic-output/SamplePlugin.translated.xml",
  targetGame: "skyrim_se",
  operationKind: "generate",
  duplicateRowCreated: false,
  retryable: true,
  isRedacted: false
}

const longPath =
  "synthetic-output/very/long/path/for/storybook/review/SamplePlugin.translated.output.artifact.with.long.name.xml"

export const outputSummaryHeaderFixtures = {
  ready: {
    gatewayStatus: "connected",
    statusTitle: "出力準備完了",
    statusText: "completed job を選択済みです。"
  }
}

export const completedJobListPanelFixtures = {
  populated: {
    completedJobs,
    selectedJobId: 2401,
    refreshDisabled: false,
    onRefresh: ignoreAction,
    onSelectJob: ignoreJob
  },
  empty: {
    completedJobs: [],
    selectedJobId: null,
    refreshDisabled: false,
    onRefresh: ignoreAction,
    onSelectJob: ignoreJob
  }
}

export const selectedJobSummaryCardFixtures = {
  selected: {
    review,
    viewState: "ready"
  },
  rejected: {
    review: rejectedReview,
    viewState: "not_ready"
  },
  awaitingSelection: {
    review: null,
    viewState: "awaiting_selection"
  }
}

export const outputActionPanelFixtures = {
  ready: {
    targetGame: "skyrim_se",
    outputPath: "synthetic-output/SamplePlugin.translated.xml",
    pathState: "valid",
    pathReason: "",
    canGenerate: true,
    canRegenerate: true,
    isSubmitting: false,
    disabledReason: "",
    onTargetGameChange: ignoreString,
    onOutputPathInput: ignoreString,
    onGenerate: ignoreAction,
    onRegenerate: ignoreAction
  },
  disabled: {
    targetGame: "skyrim_se",
    outputPath: "synthetic-output/SamplePlugin.translated.txt",
    pathState: "invalid",
    pathReason: "出力先 path は .xml で終える必要があります。",
    canGenerate: false,
    canRegenerate: false,
    isSubmitting: false,
    disabledReason: "出力先 path を確認してください。",
    onTargetGameChange: ignoreString,
    onOutputPathInput: ignoreString,
    onGenerate: ignoreAction,
    onRegenerate: ignoreAction
  },
  longPath: {
    targetGame: "skyrim_le",
    outputPath: longPath,
    pathState: "valid",
    pathReason: "",
    canGenerate: true,
    canRegenerate: false,
    isSubmitting: false,
    disabledReason: "再出力する artifact がありません。",
    onTargetGameChange: ignoreString,
    onOutputPathInput: ignoreString,
    onGenerate: ignoreAction,
    onRegenerate: ignoreAction
  }
}

export const latestOutputResultCardFixtures = {
  generated: {
    lastCommand
  },
  failed: {
    lastCommand: {
      ...lastCommand,
      artifactStatus: "failed",
      errorKind: "file_write_failed",
      errorReason: "synthetic output path is not writable."
    }
  },
  empty: {
    lastCommand: null
  }
}

export const diffPreviewPanelFixtures = {
  populated: {
    compatibilitySummaryText: "compatibility: warning 1 / rejected 0",
    diffPreview,
    onSelectArtifact: ignoreArtifact
  },
  empty: {
    compatibilitySummaryText: "compatibility: not checked",
    diffPreview: null,
    onSelectArtifact: ignoreArtifact
  }
}
