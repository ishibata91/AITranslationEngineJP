type TranslationJobManagementFilterId =
  import("./translation-job-management-screen-types").TranslationJobManagementFilterId
type TranslationJobManagementScreenViewModel =
  import("./translation-job-management-screen-types").TranslationJobManagementScreenViewModel

export type TranslationJobManagementScreenViewModelListener = (
  viewModel: TranslationJobManagementScreenViewModel
) => void

export interface TranslationJobManagementScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(
    listener: TranslationJobManagementScreenViewModelListener
  ): () => void
  getViewModel(): TranslationJobManagementScreenViewModel
  reload(): Promise<void>
  selectJob(jobId: number): Promise<void>
  setFilter(filterId: TranslationJobManagementFilterId): void
  setSearchQuery(searchQuery: string): void
  requestStop(): Promise<void>
  requestResume(): Promise<void>
  openDeleteConfirmation(): void
  closeDeleteConfirmation(): void
  deleteSelectedJob(): Promise<void>
}

export type CreateTranslationJobManagementScreenController =
  () => TranslationJobManagementScreenControllerContract
