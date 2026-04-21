export interface MasterPersonaFrontendRefresh {
  keyword: string
  pluginFilter: string
  page: number
  pageSize: number
}

export interface MasterPersonaPageRequest {
  refresh: MasterPersonaFrontendRefresh
  preferredIdentityKey?: string
}

export interface MasterPersonaPluginGroup {
  targetPlugin: string
  count: number
}

export interface MasterPersonaListItem {
  identityKey: string
  targetPlugin: string
  formId: string
  recordType: string
  editorId: string
  displayName: string
  race?: string
  sex?: string
  voiceType: string
  className: string
  sourcePlugin: string
  personaSummary: string
  updatedAt: string
}

export interface MasterPersonaDetail extends MasterPersonaListItem {
  personaBody: string
  speechStyle?: string
  runLockReason: string
}

export interface MasterPersonaPageState {
  items: MasterPersonaListItem[]
  pluginGroups: MasterPersonaPluginGroup[]
  totalCount: number
  page: number
  pageSize: number
  selectedIdentityKey?: string
}

export interface MasterPersonaPageResponse {
  page: MasterPersonaPageState
}

export interface MasterPersonaIdentityRequest {
  identityKey: string
}

export interface MasterPersonaDetailResponse {
  entry: MasterPersonaDetail
}

export interface MasterPersonaAISettings {
  provider: string
  model: string
  apiKey: string
}

export interface MasterPersonaPreviewRequest {
  filePath: string
  aiSettings: MasterPersonaAISettings
}

export interface MasterPersonaPreviewResult {
  fileName: string
  targetPlugin: string
  candidateCount: number
  newlyAddableCount: number
  existingCount: number
  status: string
}

/** @public */
export interface MasterPersonaPreviewStateEntry {
  fileName?: string
  targetPlugin?: string
  status?: string
  candidateCount?: number
  newlyAddableCount?: number
  existingCount?: number
}

export interface MasterPersonaRunStatus {
  runState: string
  targetPlugin: string
  processedCount: number
  successCount: number
  existingSkipCount: number
  currentActorLabel: string
  message: string
  startedAt?: string
  finishedAt?: string
}

export interface MasterPersonaUpdateInput {
  personaSummary?: string
  speechStyle?: string
  personaBody: string
  displayName?: string
  formId?: string
  editorId?: string
  race?: string
  sex?: string
  voiceType?: string
  className?: string
  sourcePlugin?: string
}

export interface MasterPersonaUpdateRequest {
  identityKey: string
  entry: { personaBody: string; personaSummary?: string; speechStyle?: string }
  refresh: MasterPersonaFrontendRefresh
}

export interface MasterPersonaDeleteRequest {
  identityKey: string
  refresh: MasterPersonaFrontendRefresh
}

export interface MasterPersonaMutationResponse {
  page: MasterPersonaPageState
  changedEntry?: MasterPersonaDetail
  deletedEntryId?: string
}

export type MasterPersonaModalState = "edit" | "delete" | null

// State-level type extends the public seam with optional legacy backend fields
// so test fixtures that still pass zeroDialogueSkipCount / genericNpcCount compile.
type MasterPersonaRunStatusState = MasterPersonaRunStatus & {
  zeroDialogueSkipCount?: number
  genericNpcCount?: number
}

export interface MasterPersonaScreenState {
  items: MasterPersonaListItem[]
  pluginGroups: MasterPersonaPluginGroup[]
  selectedIdentityKey: string | null
  selectedEntry: MasterPersonaDetail | null
  keyword: string
  pluginFilter: string
  page: number
  pageSize: number
  totalCount: number
  errorMessage: string
  aiSettings: MasterPersonaAISettings
  aiSettingsMessage: string
  selectedFileName: string
  selectedFileReference: string | null
  preview: MasterPersonaPreviewStateEntry | null
  runStatus: MasterPersonaRunStatusState
  modalState: MasterPersonaModalState
  editForm: MasterPersonaUpdateInput
}

export interface MasterPersonaScreenViewModel extends MasterPersonaScreenState {
  gatewayStatus: string
  pluginOptions: Array<{ value: string; label: string }>
  totalPages: number
  pageStatusText: string
  selectionStatusText: string
  listHeadline: string
  detailLockText: string
  detailStatusText: string
  canStartPreview: boolean
  canStartGeneration: boolean
  canMutate: boolean
  isRunActive: boolean
  hasPreview: boolean
  aiProviderLabel: string
  promptTemplateDescription: string
  progressPercent: number
}

export const MASTER_PERSONA_PAGE_SIZE = 30
const MASTER_PERSONA_DEFAULT_PROVIDER = ""
const MASTER_PERSONA_DEFAULT_MODEL = ""
export const MASTER_PERSONA_IDLE_RUN_STATE = "入力待ち"
export const MASTER_PERSONA_PROMPT_TEMPLATE_DESCRIPTION =
  "プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。"

export function createDefaultMasterPersonaAISettings(): MasterPersonaAISettings {
  return {
    provider: MASTER_PERSONA_DEFAULT_PROVIDER,
    model: MASTER_PERSONA_DEFAULT_MODEL,
    apiKey: ""
  }
}

export function createEmptyMasterPersonaUpdateInput(): MasterPersonaUpdateInput {
  return {
    personaSummary: "",
    speechStyle: "",
    personaBody: ""
  }
}

export function buildMasterPersonaRefresh(
  keyword: string,
  pluginFilter: string,
  page: number
): MasterPersonaFrontendRefresh {
  return {
    keyword: keyword.trim(),
    pluginFilter: pluginFilter === "すべてのプラグイン" ? "" : pluginFilter,
    page,
    pageSize: MASTER_PERSONA_PAGE_SIZE
  }
}

export function buildMasterPersonaUpdateInput(
  state: MasterPersonaScreenState
): { personaSummary: string; speechStyle: string; personaBody: string } {
  return {
    personaSummary: (state.editForm.personaSummary ?? "").trim(),
    speechStyle: (state.editForm.speechStyle ?? "").trim(),
    personaBody: state.editForm.personaBody.trim()
  }
}

export interface MasterPersonaGatewayContract {
  getMasterPersonaPage(
    request: MasterPersonaPageRequest
  ): Promise<MasterPersonaPageResponse>
  getMasterPersonaDetail(
    request: MasterPersonaIdentityRequest
  ): Promise<MasterPersonaDetailResponse>
  loadMasterPersonaAISettings(): Promise<MasterPersonaAISettings>
  saveMasterPersonaAISettings(
    request: MasterPersonaAISettings
  ): Promise<MasterPersonaAISettings>
  previewMasterPersonaGeneration(
    request: MasterPersonaPreviewRequest
  ): Promise<MasterPersonaPreviewStateEntry>
  executeMasterPersonaGeneration(
    request: MasterPersonaPreviewRequest
  ): Promise<MasterPersonaRunStatus>
  getMasterPersonaRunStatus(): Promise<MasterPersonaRunStatus>
  interruptMasterPersonaGeneration(): Promise<MasterPersonaRunStatus>
  cancelMasterPersonaGeneration(): Promise<MasterPersonaRunStatus>
  updateMasterPersona(
    request: MasterPersonaUpdateRequest
  ): Promise<MasterPersonaMutationResponse>
  deleteMasterPersona(
    request: MasterPersonaDeleteRequest
  ): Promise<MasterPersonaMutationResponse>
}
