import type {
  MasterDictionaryEntryDetail,
  MasterDictionaryEntrySummary
} from "@application/gateway-contract/master-dictionary"
import type {
  ImportSummary,
  ModalState
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

export interface DictionaryHeaderProps {
  gatewayStatus: string
  errorMessage: string
}

export interface DictionaryImportPanelProps {
  hasStagedFile: boolean
  importProgress: number
  importStatusText: string
  importStatusValue: string
  importSummary: ImportSummary | null
  isImportRunning: boolean
  selectedEntry: MasterDictionaryEntryDetail | null
  selectedFileName: string
  chooseXmlFile: () => void
  handleXmlSelected: (event: Event) => void
  resetImportSelection: () => void
  startImport: () => void
}

export interface DictionaryListPanelProps {
  category: string
  categoryOptions: string[]
  entries: MasterDictionaryEntrySummary[]
  listHeadline: string
  page: number
  pageStatusText: string
  query: string
  selectedId: string | null
  selectionStatusText: string
  totalPages: number
  goToNextPage: () => void
  goToPrevPage: () => void
  handleCategoryChange: (event: Event) => void
  handleSearchInput: (event: Event) => void
  openCreateModal: () => void
  selectRow: (entryId: string) => void
}

export interface DictionaryDetailPanelProps {
  detailSublineText: string
  selectedEntry: MasterDictionaryEntryDetail | null
  openDeleteModal: () => void
  openEditModal: () => void
}

export interface DictionaryEditModalProps {
  categoryOptions: string[]
  formCategory: string
  formOrigin: string
  formSource: string
  formTranslation: string
  modalState: ModalState
  closeEditModal: () => void
  saveCurrentEntry: () => void
  setFormCategory: (event: Event) => void
  setFormOrigin: (event: Event) => void
  setFormSource: (event: Event) => void
  setFormTranslation: (event: Event) => void
}

export interface DictionaryDeleteModalProps {
  modalState: ModalState
  selectedEntry: MasterDictionaryEntryDetail | null
  closeDeleteModal: () => void
  deleteCurrentEntry: () => void
}
