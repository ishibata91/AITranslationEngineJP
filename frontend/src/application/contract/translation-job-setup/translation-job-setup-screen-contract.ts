import type { TranslationJobSetupScreenViewModel } from "@application/gateway-contract/translation-job-setup"

export type TranslationJobSetupScreenViewModelListener = (
  viewModel: TranslationJobSetupScreenViewModel
) => void

export interface TranslationJobSetupScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: TranslationJobSetupScreenViewModelListener): () => void
  getViewModel(): TranslationJobSetupScreenViewModel
  selectInputSource(inputSourceId: number): void
  selectRuntime(runtimeKey: string): void
  selectCredentialRef(credentialRef: string): void
  runValidation(): Promise<void>
  createJob(): Promise<void>
}

export type CreateTranslationJobSetupScreenController =
  () => TranslationJobSetupScreenControllerContract