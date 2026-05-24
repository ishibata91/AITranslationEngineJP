import type { MasterDictionaryEntryDetail } from "@application/gateway-contract/master-dictionary"
import type {
  ImportSummary,
  ModalState
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

export interface DictionaryImportPanelProps {
  isImportRunning: boolean
  selectedFileName: string
  chooseXmlFile: () => void
  handleXmlSelected: (event: Event) => void
}

export interface DictionaryImportProgressPanelProps {
  hasStagedFile: boolean
  importProgress: number
  importStatusText: string
  importStatusValue: string
  importSummary: ImportSummary | null
  isImportRunning: boolean
  selectedEntry: MasterDictionaryEntryDetail | null
  selectedFileName: string
  resetImportSelection: () => void
  startImport: () => void
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
