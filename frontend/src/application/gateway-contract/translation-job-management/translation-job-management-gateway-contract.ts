export type TranslationJobManagementJobState =
  | "Ready"
  | "Running"
  | "Paused"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"
  | "Completed"

export type TranslationJobManagementStateTone =
  | "neutral"
  | "info"
  | "warning"
  | "danger"

export type TranslationJobManagementCredentialState =
  | "configured"
  | "missing"
  | "inaccessible"

export type TranslationJobManagementCacheState = "available" | "missing"

export type TranslationJobManagementCurrentPhase =
  | "term_translation"
  | "persona_generation"
  | "body_translation"

export type TranslationJobManagementReasonCategory =
  | "cache_missing"
  | "terminal_state"
  | "state_projection_inconsistent"
  | "phase_progress_aggregation_failed"
  | "stale_selection"
  | "list_load_failure"
  | "running_delete_blocked"
  | "stop_requested"
  | "stop_failed"
  | "delete_failed"
  | "resume_failed"

export type TranslationJobManagementOperationKind = "stop" | "resume" | "delete"

export interface TranslationJobManagementInputSourceSummary {
  inputSourceId: number
  inputSourceLabel: string
  inputSourceKindLabel: string
  sourcePath: string
  pluginName: string
  extractedJsonLabel: string
}

export interface TranslationJobManagementProgressSummary {
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  percent: number
  progressLabel: string
  lastUpdatedLabel: string
}

export interface TranslationJobManagementProtectedSettingSummary {
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialState: TranslationJobManagementCredentialState
  credentialStateLabel: string
}

export interface TranslationJobManagementOperationAvailability {
  kind: TranslationJobManagementOperationKind
  enabled: boolean
  label: string
  helperText: string
  reasonCategory?: TranslationJobManagementReasonCategory
  reasonText?: string
}

export interface TranslationJobManagementBlockedReason {
  category: TranslationJobManagementReasonCategory
  title: string
  detail: string
}

export interface TranslationJobManagementJobSummary {
  jobId: number
  jobState: TranslationJobManagementJobState
  jobStateLabel: string
  stateTone: TranslationJobManagementStateTone
  canOpenPhase?: boolean
  openBlockedReason?: TranslationJobManagementBlockedReason
  inputSource: TranslationJobManagementInputSourceSummary
  progress: TranslationJobManagementProgressSummary
  stopAvailability: TranslationJobManagementOperationAvailability
  resumeAvailability: TranslationJobManagementOperationAvailability
  deleteAvailability: TranslationJobManagementOperationAvailability
}

export interface TranslationJobManagementJobDetail extends TranslationJobManagementJobSummary {
  cacheState: TranslationJobManagementCacheState
  cacheStateLabel: string
  runtimeSummary: TranslationJobManagementProtectedSettingSummary
  resumeBlockedReasons: TranslationJobManagementBlockedReason[]
  warnings: TranslationJobManagementBlockedReason[]
  deleteImpactLines: string[]
}

export interface TranslationJobManagementListResponse {
  jobs: TranslationJobManagementJobSummary[]
}

export interface TranslationJobManagementGetDetailRequest {
  jobId: number
}

export interface TranslationJobManagementDeleteRequest {
  jobId: number
}

export interface TranslationJobManagementActionRequest {
  jobId: number
}

export interface TranslationJobManagementActionResponse {
  message: string
  tone: "info" | "success" | "warning" | "danger"
  detail?: TranslationJobManagementJobDetail
  deletedJobId?: number
  reasonCategory?: TranslationJobManagementReasonCategory
}

export interface TranslationJobManagementGatewayContract {
  ListIncompleteJobs(): Promise<TranslationJobManagementListResponse>
  GetJobDetail(
    request: TranslationJobManagementGetDetailRequest
  ): Promise<TranslationJobManagementJobDetail>
  RequestStop(
    request: TranslationJobManagementActionRequest
  ): Promise<TranslationJobManagementActionResponse>
  ResumeJob(
    request: TranslationJobManagementActionRequest
  ): Promise<TranslationJobManagementActionResponse>
  DeleteJob(
    request: TranslationJobManagementDeleteRequest
  ): Promise<TranslationJobManagementActionResponse>
}
