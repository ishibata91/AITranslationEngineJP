import type {
  ImportMasterDictionaryXmlResponse,
  MasterDictionaryEntryDetail,
  MasterDictionaryEntrySummary,
  MasterDictionaryPageState
} from "@application/gateway-contract/master-dictionary"

export type ModalState = "create" | "edit" | "delete" | null
export type ImportStage = "idle" | "ready" | "running" | "done"

export interface ImportSummary {
  fileName: string
  importedCount: number
  totalCount: number
  selectedSource: string
}

export interface MasterDictionaryScreenState {
  entries: MasterDictionaryEntrySummary[]
  selectedEntry: MasterDictionaryEntryDetail | null
  selectedId: string | null
  totalCount: number
  query: string
  category: string
  page: number
  errorMessage: string
  modalState: ModalState
  formSource: string
  formCategory: string
  formOrigin: string
  formTranslation: string
  selectedFileName: string
  selectedFileReference: string | null
  importStage: ImportStage
  importProgress: number
  importSummary: ImportSummary | null
}

export interface MasterDictionaryScreenViewModel extends MasterDictionaryScreenState {
  gatewayStatus: string
  hasStagedFile: boolean
  isImportRunning: boolean
  importStatusValue: string
  importStatusText: string
  categoryOptions: string[]
  totalPages: number
  pageStatusText: string
  listHeadline: string
  selectionStatusText: string
  detailSublineText: string
}

export interface RuntimeImportProgressPayload {
  progress?: number
}

export interface RuntimeImportCompletedPayload {
  page?: MasterDictionaryPageState
  summary?: ImportMasterDictionaryXmlResponse["summary"]
}
