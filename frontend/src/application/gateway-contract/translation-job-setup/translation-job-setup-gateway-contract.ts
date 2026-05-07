import type { ModelSettingsCardState } from "@application/gateway-contract/model-settings-card"

export interface TranslationJobSetupInputCandidate {
  id: number
  label: string
  sourceKind: string
  registeredAt?: string
  recordCount: number
  existingJob?: TranslationJobSetupExistingJob
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

export type TranslationJobSetupPhaseId =
  | "word_translation"
  | "npc_persona_generation"
  | "text_translation"

export type TranslationJobSetupCredentialRequirement =
  | "required"
  | "not_required"

export type TranslationJobSetupCredentialStatus =
  | "configured"
  | "missing"
  | "not_required"

export type TranslationJobSetupBatchMode =
  | "disabled"
  | "enabled"
  | "unsupported"

export type TranslationJobSetupProviderModelListStatus =
  | "not_updated"
  | "loading"
  | "success"
  | "failed"
  | "credential_missing"
  | "credential_not_required"

export interface TranslationJobSetupProviderCapability {
  provider: string
  credentialRequirement: TranslationJobSetupCredentialRequirement
  supportedExecutionModes: string[]
  supportsBatchMode: boolean
}

export interface TranslationJobSetupCredentialReference {
  provider: string
  isConfigured: boolean
  isMissingSecret: boolean
}

export interface TranslationJobSetupPhaseRuntimeDraft {
  phaseId: TranslationJobSetupPhaseId
  provider: string
  model: string
  credentialStatus: TranslationJobSetupCredentialStatus
  executionMode: string
  batchMode: TranslationJobSetupBatchMode
}

export interface TranslationJobSetupOptionsResponse {
  inputCandidates: TranslationJobSetupInputCandidate[]
  existingJob?: TranslationJobSetupExistingJob
  sharedDictionaries: TranslationJobSetupDictionaryOption[]
  sharedPersonas: TranslationJobSetupPersonaOption[]
  aiRuntimeOptions: TranslationJobSetupRuntimeOption[]
  providerCapabilities?: TranslationJobSetupProviderCapability[]
  phaseRuntimeDrafts?: TranslationJobSetupPhaseRuntimeDraft[]
}

export interface TranslationJobSetupRuntimeSelection {
  provider: string
  model: string
  executionMode: string
}

export interface TranslationJobSetupPhaseRuntimeSelection {
  phaseId: TranslationJobSetupPhaseId
  provider: string
  model: string
  credentialStatus: TranslationJobSetupCredentialStatus
  executionMode: string
  batchMode: TranslationJobSetupBatchMode
}

export interface TranslationJobSetupPhaseRuntimeValidationSelection extends TranslationJobSetupPhaseRuntimeSelection {
  modelListFreshnessToken: string
}

export interface ListTranslationJobSetupProviderModelsRequest {
  phaseId: TranslationJobSetupPhaseId
  provider: string
  credentialStatus: TranslationJobSetupCredentialStatus
  requestToken: string
}

export interface TranslationJobSetupProviderModelOption {
  modelId: string
  label: string
}

export interface ListTranslationJobSetupProviderModelsResponse {
  phaseId: TranslationJobSetupPhaseId
  provider: string
  credentialStatus: TranslationJobSetupCredentialStatus
  requestToken: string
  sourceToken: string
  status: TranslationJobSetupProviderModelListStatus
  models: TranslationJobSetupProviderModelOption[]
  failureKind?: TranslationJobSetupCreateErrorKind
}

export interface ValidateTranslationJobSetupRequest {
  inputSourceId: number
  runtime: TranslationJobSetupRuntimeSelection
  phaseRuntimeSelections?: TranslationJobSetupPhaseRuntimeValidationSelection[]
}

export interface TranslationJobSetupPhaseValidationResult {
  phaseId: TranslationJobSetupPhaseId
  status: string
  blockingFailureCategory?: string
  canCreate: boolean
  modelListState: TranslationJobSetupProviderModelListStatus
  isModelSelectionStale: boolean
}

export interface TranslationJobSetupValidationResponse {
  status: string
  blockingFailureCategory?: string
  targetSlices: string[]
  validatedAt: string
  canCreate: boolean
  passSlices: string[]
  phaseResults?: TranslationJobSetupPhaseValidationResult[]
  staleModelListPhaseIds?: TranslationJobSetupPhaseId[]
}

export interface CreateTranslationJobRequest {
  inputSourceId: number
  inputSource: string
  validationStatus: string
  validatedAt: string
  validationPassSlices: string[]
  runtime: TranslationJobSetupRuntimeSelection
  phaseRuntimeSelections?: TranslationJobSetupPhaseRuntimeValidationSelection[]
}

export interface TranslationJobExecutionSummary {
  provider: string
  model: string
  executionMode: string
}

export interface TranslationJobSetupPhaseRuntimeSummary {
  phaseId: TranslationJobSetupPhaseId
  provider: string
  model: string
  credentialStatus: TranslationJobSetupCredentialStatus
  executionMode: string
  batchMode: TranslationJobSetupBatchMode
}

export type TranslationJobSetupCreateErrorKind =
  | "phase_runtime_missing"
  | "required_setting_missing"
  | "input_not_found"
  | "cache_missing"
  | "foundation_ref_missing"
  | "credential_missing"
  | "model_list_credential_missing"
  | "model_list_failed"
  | "model_selection_stale"
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
  phaseRuntimeSummaries?: TranslationJobSetupPhaseRuntimeSummary[]
}

export interface DeleteTranslationJobSetupInputRequest {
  inputSourceId: number
}

export type TranslationJobSetupDeleteInputErrorKind =
  | "input_not_found"
  | "input_delete_blocked"

export interface DeleteTranslationJobSetupInputResponse {
  deletedInputSourceId?: number
  errorKind?: TranslationJobSetupDeleteInputErrorKind
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
  phaseRuntimeSummaries?: TranslationJobSetupPhaseRuntimeSummary[]
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
  deletingInputSourceId: number | null
  selectedRuntimeKey: string | null
  selectedCredentialRef: string
  phaseRuntimeSelections?: TranslationJobSetupPhaseRuntimeSelection[]
  providerModelLists?: ListTranslationJobSetupProviderModelsResponse[]
  modelSettingsCards?: ModelSettingsCardState[]
  validationResult: TranslationJobSetupValidationResponse | null
  validationState: TranslationJobSetupValidationState
  dirty: boolean
  errorMessage: string
  createErrorKind: TranslationJobSetupCreateErrorKind | null
  summary: TranslationJobSetupSummaryResponse | null
}

export interface TranslationJobSetupScreenViewModel extends TranslationJobSetupScreenState {
  gatewayStatus: string
  selectedInputCandidate: TranslationJobSetupInputCandidate | null
  selectedRuntimeOption: TranslationJobSetupRuntimeOption | null
  availableCredentialRefs: TranslationJobSetupCredentialReference[]
  phaseValidationResults?: TranslationJobSetupPhaseValidationResult[]
  phaseRuntimeSummaries?: TranslationJobSetupPhaseRuntimeSummary[]
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
  listTranslationJobSetupProviderModels(
    request: ListTranslationJobSetupProviderModelsRequest
  ): Promise<ListTranslationJobSetupProviderModelsResponse>
  validateTranslationJobSetup(
    request: ValidateTranslationJobSetupRequest
  ): Promise<TranslationJobSetupValidationResponse>
  createTranslationJob(
    request: CreateTranslationJobRequest
  ): Promise<CreateTranslationJobResponse>
  deleteTranslationJobSetupInput(
    request: DeleteTranslationJobSetupInputRequest
  ): Promise<DeleteTranslationJobSetupInputResponse>
  getTranslationJobSetupSummary(
    request: GetTranslationJobSetupSummaryRequest
  ): Promise<TranslationJobSetupSummaryResponse>
}

export function createTranslationJobSetupRuntimeKey(
  option: TranslationJobSetupRuntimeOption
): string {
  return [option.provider, option.model, option.mode].join("::")
}
