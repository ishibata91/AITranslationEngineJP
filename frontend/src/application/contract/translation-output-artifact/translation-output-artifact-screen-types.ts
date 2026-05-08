import type {
  TranslationOutputArtifactErrorSummary,
  TranslationOutputArtifactStatusSummary,
  TranslationOutputCompletedJobSummary
} from "@application/gateway-contract/translation-output-artifact"

export type TranslationOutputArtifactErrorKind =
  | "not_completed"
  | "canceled"
  | "status_mismatch"
  | "missing_required_row_field"
  | "unknown_output_status"
  | "xml_serialization_failed"
  | "file_write_failed"
  | "artifact_save_failed"
  | "compatibility_rejected"
  | "secret_redacted"

export type TranslationOutputArtifactViewState =
  | "loading"
  | "empty"
  | "awaiting_selection"
  | "not_ready"
  | "ready"
  | "generating"
  | "success"
  | "failed"
  | "stale"

export type TranslationOutputArtifactActionKind = "generate" | "regenerate"

export interface TranslationOutputReviewSnapshot {
  completedJobs: TranslationOutputCompletedJobSummary[]
  selectedJobId: number
  selectedJobStatus: string
  bodyPhaseStatus: string
  outputReady: boolean
  translatedCount: number
  rowCount: number
  inputSnapshotDigest: string
  sourceFileDigest: string
  readiness: boolean
  retryable: boolean
  artifactId: number
  artifactStatus: string
  artifactRowCount: number
  currentVersion: boolean
  rejectionKind?: TranslationOutputArtifactErrorKind
  rejectionReasons: TranslationOutputArtifactErrorSummary[]
}

export interface TranslationOutputDiffPreviewSnapshotRow {
  fieldId: number
  rowDigest: string
  edid: string
  rec: string
  field: string
  formId: string
  sourceExcerpt: string
  destExcerpt: string
  xTranslatorStatus: number
  internalOutputStatus: string
  rowReflectionSummary: string
  staleReason?: string
  canRegenerate: boolean
}

export interface TranslationOutputDiffPreviewSnapshot {
  artifactId: number
  rows: TranslationOutputDiffPreviewSnapshotRow[]
  compatibilityPassed: boolean
  compatibilityWarningCount: number
  compatibilityRejectCount: number
}

export interface TranslationOutputArtifactCommandSnapshot {
  jobId: number
  artifactId: number
  artifactStatus: string
  rowCount: number
  filePath?: string
  targetGame: string
  operationKind?: string
  replacedArtifactId?: number
  affectedFieldIds?: number[]
  duplicateRowCreated?: boolean
  errorKind?: TranslationOutputArtifactErrorKind
  errorReason?: string
  retryable?: boolean
  isRedacted?: boolean
}

export interface TranslationOutputArtifactActionDisablement {
  actionKind: TranslationOutputArtifactActionKind
  disabled: boolean
  reason?: string
  errorKind?: TranslationOutputArtifactErrorKind
}

export interface TranslationOutputArtifactScreenState {
  phase: "idle" | "loading" | "ready" | "submitting"
  viewState: TranslationOutputArtifactViewState
  completedJobs: TranslationOutputCompletedJobSummary[]
  selectedJobId: number | null
  selectedArtifactId: number | null
  review: TranslationOutputReviewSnapshot | null
  diffPreview: TranslationOutputDiffPreviewSnapshot | null
  lastCommand: TranslationOutputArtifactCommandSnapshot | null
  actionDisablements: TranslationOutputArtifactActionDisablement[]
  refreshPending: boolean
  targetGame: string
  outputPath: string
  pathState: "empty" | "valid" | "invalid" | "readonly"
  pathReason: string
  errorMessage: string
  pendingAction: TranslationOutputArtifactActionKind | "refresh" | null
  hasLoaded: boolean
}

export interface TranslationOutputArtifactScreenViewModel extends TranslationOutputArtifactScreenState {
  canGenerate: boolean
  canRegenerate: boolean
  primaryAction: TranslationOutputArtifactActionKind
  disabledReason?: string
  selectedJobStatus?: string
  selectedArtifactStatus?: string
  gatewayStatus: string
  isLoading: boolean
  isSubmitting: boolean
  statusTitle: string
  statusText: string
  artifactStatusSummary: TranslationOutputArtifactStatusSummary | null
  compatibilitySummaryText: string
}
