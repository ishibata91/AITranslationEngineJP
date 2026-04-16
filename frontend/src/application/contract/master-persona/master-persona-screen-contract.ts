import type { MasterPersonaScreenViewModel } from "./master-persona-screen-types"

export type MasterPersonaScreenViewModelListener = (
  viewModel: MasterPersonaScreenViewModel
) => void

export interface MasterPersonaScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: MasterPersonaScreenViewModelListener): () => void
  getViewModel(): MasterPersonaScreenViewModel
  selectRow(identityKey: string): Promise<void>
  handleSearchInput(event: Event): void
  handlePluginFilterChange(event: Event): void
  goToPrevPage(): void
  goToNextPage(): void
  stageJsonSelection(file: File | null): void
  resetJsonSelection(): void
  previewGeneration(): Promise<void>
  executeGeneration(): Promise<void>
  interruptGeneration(): Promise<void>
  cancelGeneration(): Promise<void>
  saveAISettings(): Promise<void>
  setAIProvider(event: Event): void
  setAIModel(event: Event): void
  setAPIKey(event: Event): void
  openDialogueModal(): Promise<void>
  closeDialogueModal(): void
  openEditModal(): void
  closeEditModal(): void
  openDeleteModal(): void
  closeDeleteModal(): void
  saveCurrentEntry(): Promise<void>
  deleteCurrentEntry(): Promise<void>
  setEditFormField(field: keyof MasterPersonaEditableFieldMap, event: Event): void
}

export interface MasterPersonaEditableFieldMap {
  formId: HTMLInputElement
  editorId: HTMLInputElement
  displayName: HTMLInputElement
  race: HTMLInputElement
  sex: HTMLInputElement
  voiceType: HTMLInputElement
  className: HTMLInputElement
  sourcePlugin: HTMLInputElement
  personaBody: HTMLTextAreaElement
}

export type CreateMasterPersonaScreenController =
  () => MasterPersonaScreenControllerContract
