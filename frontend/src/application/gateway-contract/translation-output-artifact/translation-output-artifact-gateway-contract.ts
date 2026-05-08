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

export interface GetTranslationOutputReviewRequest {
  selectedJobId?: number
}

export interface GetTranslationOutputDiffPreviewRequest {
  jobId: number
  artifactId: number
}

export interface GenerateXTranslatorOutputArtifactRequest {
  jobId: number
  targetGame: string
  outputPath: string
}

export interface RegenerateXTranslatorOutputArtifactRequest {
  jobId: number
  artifactId: number
  targetGame: string
  outputPath: string
}

export interface TranslationOutputCompletedJobSummary {
  jobId: number
  jobStatus: string
  artifactStatus: string
  outputReady: boolean
  translatedCount: number
  outputStatusDistribution?: Record<string, number>
}

export interface TranslationOutputInputProvenanceSummary {
  inputSnapshotDigest: string
  sourceFileDigest: string
}

export interface TranslationOutputResultSummary {
  translatedCount: number
  rowCount: number
  inputProvenance: TranslationOutputInputProvenanceSummary
}

export interface TranslationOutputSelectedJobSummary {
  jobId: number
  jobStatus: string
  bodyPhaseStatus: string
  outputReady: boolean
  resultSummary: TranslationOutputResultSummary
}

export interface TranslationOutputReadinessSummary {
  ready: boolean
  retryable: boolean
  rejectionKind?: TranslationOutputArtifactErrorKind
}

export interface TranslationOutputArtifactStatusSummary {
  artifactId: number
  status: string
  rowCount: number
  currentVersion: boolean
}

export interface TranslationOutputArtifactErrorSummary {
  errorKind: TranslationOutputArtifactErrorKind
  reason: string
  retryable: boolean
  isRedacted: boolean
}

export interface TranslationOutputReviewResponse {
  completedJobs: TranslationOutputCompletedJobSummary[]
  hasSelectedJob?: boolean
  selectedJob: TranslationOutputSelectedJobSummary
  outputReadiness: TranslationOutputReadinessSummary
  artifactStatus: TranslationOutputArtifactStatusSummary
  rejectionReasons?: TranslationOutputArtifactErrorSummary[]
}

export interface TranslationOutputDiffPreviewRow {
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

export interface TranslationOutputCompatibilitySummary {
  passed: boolean
  warningCount: number
  rejectCount: number
}

export interface TranslationOutputDiffPreviewResponse {
  jobId: number
  artifactId: number
  rows: TranslationOutputDiffPreviewRow[]
  compatibilitySummary: TranslationOutputCompatibilitySummary
}

export interface TranslationOutputOperationSummary {
  operationKind: string
  replacedArtifactId: number
  affectedFieldIds?: number[]
  duplicateRowCreated: boolean
}

export interface TranslationOutputArtifactCommandResponse {
  jobId: number
  artifactId: number
  artifactStatus: string
  rowCount: number
  filePath?: string
  targetGame: string
  errorSummary?: TranslationOutputArtifactErrorSummary
  operationSummary?: TranslationOutputOperationSummary
}

export interface TranslationOutputArtifactGatewayContract {
  getTranslationOutputReview(
    request: GetTranslationOutputReviewRequest
  ): Promise<TranslationOutputReviewResponse>
  getTranslationOutputDiffPreview(
    request: GetTranslationOutputDiffPreviewRequest
  ): Promise<TranslationOutputDiffPreviewResponse>
  generateXTranslatorOutputArtifact(
    request: GenerateXTranslatorOutputArtifactRequest
  ): Promise<TranslationOutputArtifactCommandResponse>
  regenerateXTranslatorOutputArtifact(
    request: RegenerateXTranslatorOutputArtifactRequest
  ): Promise<TranslationOutputArtifactCommandResponse>
}
