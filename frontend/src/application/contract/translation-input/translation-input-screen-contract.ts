import type {
  CreateTranslationJobFromInputResponse,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"

export type TranslationInputScreenViewModelListener = (
  viewModel: TranslationInputScreenViewModel
) => void

export interface TranslationInputScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: TranslationInputScreenViewModelListener): () => void
  getViewModel(): TranslationInputScreenViewModel
  selectItem(localId: string): void
  stageJsonImport(file: File | null): Promise<void>
  resetImportSelection(): void
  startImport(): Promise<void>
  rebuildSelected(): Promise<void>
  createTranslationJobFromSelected?: () => Promise<CreateTranslationJobFromInputResponse | null>
}

export type CreateTranslationInputScreenController =
  () => TranslationInputScreenControllerContract
