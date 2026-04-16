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
  dialogueCount: number
  updatedAt: string
}

export interface MasterPersonaDetail extends MasterPersonaListItem {
  personaBody: string
  generationSourceJson: string
  baselineApplied: boolean
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

export interface MasterPersonaDialogueLine {
  index: number
  text: string
}

export interface MasterPersonaDialogueListResponse {
  identityKey: string
  dialogueCount: number
  dialogues: MasterPersonaDialogueLine[]
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
  totalNpcCount: number
  generatableCount: number
  existingSkipCount: number
  zeroDialogueSkipCount: number
  genericNpcCount: number
  status: string
}

export interface MasterPersonaRunStatus {
  runState: string
  targetPlugin: string
  processedCount: number
  successCount: number
  existingSkipCount: number
  zeroDialogueSkipCount: number
  genericNpcCount: number
  currentActorLabel: string
  message: string
  startedAt?: string
  finishedAt?: string
}

export interface MasterPersonaUpdateInput {
  formId: string
  editorId: string
  displayName: string
  race?: string
  sex?: string
  voiceType: string
  className: string
  sourcePlugin: string
  personaBody: string
}

export interface MasterPersonaUpdateRequest {
  identityKey: string
  entry: MasterPersonaUpdateInput
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

export interface MasterPersonaScreenState {
  items: MasterPersonaListItem[]
  pluginGroups: MasterPersonaPluginGroup[]
  selectedIdentityKey: string | null
  selectedEntry: MasterPersonaDetail | null
  dialogueModalOpen: boolean
  dialogues: MasterPersonaDialogueLine[]
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
  preview: MasterPersonaPreviewResult | null
  runStatus: MasterPersonaRunStatus
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
    formId: "",
    editorId: "",
    displayName: "",
    voiceType: "",
    className: "",
    sourcePlugin: "",
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
): MasterPersonaUpdateInput {
  return {
    formId: state.editForm.formId.trim(),
    editorId: state.editForm.editorId.trim(),
    displayName: state.editForm.displayName.trim(),
    race: normalizeOptionalField(state.editForm.race),
    sex: normalizeOptionalField(state.editForm.sex),
    voiceType: state.editForm.voiceType.trim(),
    className: state.editForm.className.trim(),
    sourcePlugin: state.editForm.sourcePlugin.trim(),
    personaBody: state.editForm.personaBody.trim()
  }
}

function normalizeOptionalField(value: string | undefined): string | undefined {
  const trimmed = value?.trim() ?? ""
  return trimmed === "" ? undefined : trimmed
}

export interface MasterPersonaGatewayContract {
  getMasterPersonaPage(
    request: MasterPersonaPageRequest
  ): Promise<MasterPersonaPageResponse>
  getMasterPersonaDetail(
    request: MasterPersonaIdentityRequest
  ): Promise<MasterPersonaDetailResponse>
  getMasterPersonaDialogueList(
    request: MasterPersonaIdentityRequest
  ): Promise<MasterPersonaDialogueListResponse>
  loadMasterPersonaAISettings(): Promise<MasterPersonaAISettings>
  saveMasterPersonaAISettings(
    request: MasterPersonaAISettings
  ): Promise<MasterPersonaAISettings>
  previewMasterPersonaGeneration(
    request: MasterPersonaPreviewRequest
  ): Promise<MasterPersonaPreviewResult>
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
