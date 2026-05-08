import type {
  TranslationJobManagementJobDetail,
  TranslationJobManagementJobSummary,
  TranslationJobManagementBlockedReason,
  TranslationJobManagementCurrentPhase,
  TranslationJobManagementOperationKind,
  TranslationJobManagementReasonCategory,
  TranslationJobManagementStateTone
} from "@application/gateway-contract/translation-job-management"

export type TranslationJobManagementScreenPhase =
  | "idle"
  | "loading"
  | "ready"
  | "empty"
  | "error"

export type TranslationJobManagementDetailPhase =
  | "idle"
  | "loading"
  | "ready"
  | "stale"

export type TranslationJobManagementFilterId =
  | "all"
  | "Ready"
  | "Running"
  | "Paused"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"

export interface TranslationJobManagementFeedbackState {
  tone: "info" | "success" | "warning" | "danger"
  title: string
  message: string
  category?:
    | TranslationJobManagementReasonCategory
    | "stop_requested"
    | "delete_success"
    | "resume_success"
}

export interface TranslationJobManagementScreenState {
  phase: TranslationJobManagementScreenPhase
  jobs: TranslationJobManagementJobSummary[]
  selectedJobId: number | null
  selectedJobDetail: TranslationJobManagementJobDetail | null
  detailPhase: TranslationJobManagementDetailPhase
  filterId: TranslationJobManagementFilterId
  searchQuery: string
  isReloading: boolean
  activeOperation: TranslationJobManagementOperationKind | "reload" | null
  isDeleteConfirmationOpen: boolean
  feedback: TranslationJobManagementFeedbackState | null
}

export interface TranslationJobManagementFilterChipViewModel {
  id: TranslationJobManagementFilterId
  label: string
  count: number
  selected: boolean
}

export interface TranslationJobManagementOperationViewModel {
  kind: TranslationJobManagementOperationKind
  label: string
  helperText: string
  enabled: boolean
  reasonText: string
  busy: boolean
}

export type TranslationJobManagementJobState =
  import("@application/gateway-contract/translation-job-management/translation-job-management-gateway-contract").TranslationJobManagementJobState

export interface TranslationJobManagementJobCardViewModel {
  jobId: number
  title: string
  jobState: TranslationJobManagementJobState
  stateLabel: string
  stateDescription: string
  stateTone: TranslationJobManagementStateTone
  inputSourceLabel: string
  sourcePath: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  lastUpdatedLabel: string
  canOpenPhase: boolean
  openBlockedReason: TranslationJobManagementBlockedReason | null
  openBlockedReasonText: string
  jobRunTarget: TranslationJobManagementJobRunTarget | null
  isSelected: boolean
  stopOperation: TranslationJobManagementOperationViewModel
  resumeOperation: TranslationJobManagementOperationViewModel
  deleteOperation: TranslationJobManagementOperationViewModel
}

export interface TranslationJobManagementReasonItemViewModel {
  category: TranslationJobManagementReasonCategory
  categoryLabel: string
  title: string
  detail: string
}

export interface TranslationJobManagementSelectedJobViewModel {
  jobId: number
  jobIdLabel: string
  jobState: TranslationJobManagementJobState
  stateLabel: string
  stateDescription: string
  stateTone: TranslationJobManagementStateTone
  inputSourceLabel: string
  inputSourceKindLabel: string
  sourcePath: string
  pluginName: string
  extractedJsonLabel: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  lastUpdatedLabel: string
  cacheStateLabel: string
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialStateLabel: string
  stopOperation: TranslationJobManagementOperationViewModel
  resumeOperation: TranslationJobManagementOperationViewModel
  deleteOperation: TranslationJobManagementOperationViewModel
  resumeBlockedReasons: TranslationJobManagementReasonItemViewModel[]
  warnings: TranslationJobManagementReasonItemViewModel[]
  deleteImpactLines: string[]
}

export interface TranslationJobManagementDeleteConfirmationViewModel {
  title: string
  lines: string[]
  confirmLabel: string
  cancelLabel: string
  busy: boolean
}

export interface TranslationJobManagementJobRunTarget {
  jobId: number
  stateLabel: string
  stateDescription: string
  currentPhase: TranslationJobManagementCurrentPhase
  currentPhaseLabel: string
  progressLabel: string
  inputSourceLabel: string
  sourcePath: string
}

export interface TranslationJobManagementScreenViewModel {
  gatewayStatus: string
  pageTitle: string
  pageLead: string
  headerCountLabel: string
  listEmptyTitle: string
  listEmptyDescription: string
  listErrorTitle: string
  listErrorDescription: string
  detailPlaceholderTitle: string
  detailPlaceholderDescription: string
  phase: TranslationJobManagementScreenPhase
  detailPhase: TranslationJobManagementDetailPhase
  isReloading: boolean
  searchQuery: string
  filterChips: TranslationJobManagementFilterChipViewModel[]
  jobs: TranslationJobManagementJobCardViewModel[]
  feedback: TranslationJobManagementFeedbackState | null
  selectedJob: TranslationJobManagementSelectedJobViewModel | null
  deleteConfirmation: TranslationJobManagementDeleteConfirmationViewModel | null
  jobRunTarget: TranslationJobManagementJobRunTarget | null
}
