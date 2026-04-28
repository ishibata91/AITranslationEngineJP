export interface TranslationJobSetupInputCandidate {
  id: number
  label: string
  sourceKind: string
  registeredAt?: string
  recordCount: number
}

export interface TranslationJobSetupExistingJob {
  inputSourceId?: number
  jobId: number
  status: string
  inputSource: string
}

export interface TranslationJobSetupDictionaryOption {
  id: string
  label: string
}

export interface TranslationJobSetupPersonaOption {
  id: string
  label: string
}

export interface TranslationJobSetupRuntimeOption {
  provider: string
  model: string
  mode: string
}

export interface TranslationJobSetupCredentialReference {
  provider: string
  credentialRef: string
  isConfigured: boolean
  isMissingSecret: boolean
}

export interface TranslationJobSetupOptionsResponse {
  inputCandidates: TranslationJobSetupInputCandidate[]
  existingJob?: TranslationJobSetupExistingJob
  sharedDictionaries: TranslationJobSetupDictionaryOption[]
  sharedPersonas: TranslationJobSetupPersonaOption[]
  aiRuntimeOptions: TranslationJobSetupRuntimeOption[]
  credentialRefs: TranslationJobSetupCredentialReference[]
}

export interface TranslationJobSetupRuntimeSelection {
  provider: string
  model: string
  executionMode: string
}

export interface ValidateTranslationJobSetupRequest {
  inputSourceId: number
  runtime: TranslationJobSetupRuntimeSelection
  credentialRef: string
}

export interface TranslationJobSetupValidationResponse {
  status: string
  blockingFailureCategory?: string
  targetSlices: string[]
  validatedAt: string
  canCreate: boolean
  passSlices: string[]
}

export interface CreateTranslationJobRequest {
  inputSourceId: number
  inputSource: string
  validationStatus: string
  validatedAt: string
  validationPassSlices: string[]
  runtime: TranslationJobSetupRuntimeSelection
  credentialRef: string
}

export interface TranslationJobExecutionSummary {
  provider: string
  model: string
  executionMode: string
}

export type TranslationJobSetupCreateErrorKind =
  | "required_setting_missing"
  | "input_not_found"
  | "cache_missing"
  | "foundation_ref_missing"
  | "credential_missing"
  | "provider_mode_unsupported"
  | "provider_unreachable"
  | "duplicate_job_for_input"
  | "validation_stale"
  | "partial_create_failed"
  | "ready_required"

export interface CreateTranslationJobResponse {
  jobId: number
  jobState: string
  inputSource: string
  executionSummary?: TranslationJobExecutionSummary
  validationPassSlices: string[]
  errorKind?: TranslationJobSetupCreateErrorKind
}

export interface GetTranslationJobSetupSummaryRequest {
  jobId: number
}

export interface TranslationJobSetupSummaryResponse {
  jobId: number
  jobState: string
  inputSource: string
  executionSummary: TranslationJobExecutionSummary
  validationPassSlices: string[]
  canStartPhase: boolean
}

export type TranslationJobSetupScreenPhase =
  | "idle"
  | "loading"
  | "ready"
  | "validating"
  | "creating"
  | "summary"

export type TranslationJobSetupValidationState =
  | "not-run"
  | "running"
  | "fresh"
  | "stale"

export interface TranslationJobSetupScreenState {
  phase: TranslationJobSetupScreenPhase
  options: TranslationJobSetupOptionsResponse | null
  selectedInputSourceId: number | null
  selectedRuntimeKey: string | null
  selectedCredentialRef: string
  validationResult: TranslationJobSetupValidationResponse | null
  validationState: TranslationJobSetupValidationState
  dirty: boolean
  errorMessage: string
  createErrorKind: TranslationJobSetupCreateErrorKind | null
  summary: TranslationJobSetupSummaryResponse | null
}

export interface TranslationJobSetupScreenViewModel
  extends TranslationJobSetupScreenState {
  gatewayStatus: string
  selectedInputCandidate: TranslationJobSetupInputCandidate | null
  selectedRuntimeOption: TranslationJobSetupRuntimeOption | null
  availableCredentialRefs: TranslationJobSetupCredentialReference[]
  selectedInputLabel: string
  selectedInputSourceKind: string
  selectedInputRecordCountLabel: string
  selectedInputRegisteredAtLabel: string
  existingJobSummary: string
  dictionaryLabels: string[]
  personaLabels: string[]
  validationStatusLabel: string
  validationStatusText: string
  createStatusText: string
  blockedReasons: string[]
  canValidate: boolean
  canCreate: boolean
  isLoading: boolean
  isValidating: boolean
  isCreating: boolean
  hasExistingJob: boolean
  showCacheMissingGuidance: boolean
  credentialStateText: string
}

export interface TranslationJobSetupGatewayContract {
  getTranslationJobSetupOptions(): Promise<TranslationJobSetupOptionsResponse>
  validateTranslationJobSetup(
    request: ValidateTranslationJobSetupRequest
  ): Promise<TranslationJobSetupValidationResponse>
  createTranslationJob(
    request: CreateTranslationJobRequest
  ): Promise<CreateTranslationJobResponse>
  getTranslationJobSetupSummary(
    request: GetTranslationJobSetupSummaryRequest
  ): Promise<TranslationJobSetupSummaryResponse>
}

export function createTranslationJobSetupRuntimeKey(
  option: TranslationJobSetupRuntimeOption
): string {
  return [option.provider, option.model, option.mode].join("::")
}