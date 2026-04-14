import type { MasterDictionaryScreenViewModel } from "./master-dictionary-screen-types"

export type MasterDictionaryScreenViewModelListener = (
  viewModel: MasterDictionaryScreenViewModel
) => void

export interface MasterDictionaryScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: MasterDictionaryScreenViewModelListener): () => void
  getViewModel(): MasterDictionaryScreenViewModel
  selectRow(id: string): Promise<void>
  openCreateModal(): void
  openEditModal(): void
  openDeleteModal(): void
  closeEditModal(): void
  closeDeleteModal(): void
  saveCurrentEntry(): Promise<void>
  deleteCurrentEntry(): Promise<void>
  handleSearchInput(event: Event): void
  handleCategoryChange(event: Event): void
  goToPrevPage(): void
  goToNextPage(): void
  stageXmlImport(file: File | null): void
  resetImportSelection(): void
  startImport(): Promise<void>
  setFormSource(event: Event): void
  setFormCategory(event: Event): void
  setFormOrigin(event: Event): void
  setFormTranslation(event: Event): void
}

export type CreateMasterDictionaryScreenController =
  () => MasterDictionaryScreenControllerContract
