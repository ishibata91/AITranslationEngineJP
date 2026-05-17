import type { MasterPersonaEditableFieldMap } from "@application/contract/master-persona/master-persona-screen-contract"
import type { ModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"
import type {
  MasterPersonaAISettings,
  MasterPersonaDetail,
  MasterPersonaListItem,
  MasterPersonaModalState,
  MasterPersonaModelOption,
  MasterPersonaPreviewStateEntry,
  MasterPersonaRunStatus,
  MasterPersonaUpdateInput
} from "@application/gateway-contract/master-persona/master-persona-gateway-contract"

export interface GenerationSetupPanelProps {
  aiSettings: MasterPersonaAISettings
  aiSettingsStatusText: string
  aiSettingsWarningText: string
  aiProviderLabel: string
  canSelectModel: boolean
  canStartGeneration: boolean
  executionMethodOptions: Array<{ value: string; label: string }>
  isRunActive: boolean
  modelOptions: MasterPersonaModelOption[]
  modelSettingsCardViewModel?: ModelSettingsCardViewModel
  preview: MasterPersonaPreviewStateEntry | null
  selectedFileName: string
  selectedFileReference: string | null
  isAISettingsRefreshing: boolean
  handleJsonSelected: (event: Event) => void
  chooseJsonFile: () => void
  resetJsonSelection: () => void
  handleAIProviderChange: (event: Event) => void
  handleAIModelChange: (event: Event) => void
  handleAIExecutionMethodChange: (event: Event) => void
  refreshAISettings: () => Promise<void>
  startGeneration: () => void
  saveAISettings: () => void
}

export interface RunStatusPanelProps {
  isRunActive: boolean
  progressPercent: number
  runStatus: MasterPersonaRunStatus
  interruptGeneration: () => void
  cancelGeneration: () => void
}

export interface PersonaReviewPanelProps {
  canMutate: boolean
  items: MasterPersonaListItem[]
  keyword: string
  page: number
  pageSize: number
  pluginFilter: string
  pluginOptions: Array<{ value: string; label: string }>
  selectedEntry: MasterPersonaDetail | null
  selectedIdentityKey: string | null
  totalCount: number
  totalPages: number
  selectRow: (identityKey: string) => void
  updateKeyword: (event: Event) => void
  updatePluginFilter: (event: Event) => void
  goToPrevPage: () => void
  goToNextPage: () => void
  editCurrent: () => void
  openDelete: () => void
}

export interface PersonaActionModalProps {
  modalState: MasterPersonaModalState
  selectedEntry: MasterPersonaDetail | null
  editForm: MasterPersonaUpdateInput
  errorMessage?: string
  closeEdit: () => void
  closeDelete: () => void
  saveCurrentEntry: () => void
  deleteCurrentEntry: () => void
  setEditFormField: (
    field: keyof MasterPersonaEditableFieldMap,
    event: Event
  ) => void
}
