import { createTranslationJobSetupRuntimeKey } from "@application/gateway-contract/translation-job-setup"

type TranslationJobSetupPhaseId =
  import("@application/gateway-contract/translation-job-setup").TranslationJobSetupPhaseId
type TranslationJobSetupScreenViewModel =
  import("@application/gateway-contract/translation-job-setup").TranslationJobSetupScreenViewModel

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
  selectPhaseProvider(phaseId: TranslationJobSetupPhaseId, provider: string): void
  refreshPhaseModels(phaseId: TranslationJobSetupPhaseId): Promise<void>
  selectPhaseModel(phaseId: TranslationJobSetupPhaseId, model: string): void
  togglePhaseBatchMode(
    phaseId: TranslationJobSetupPhaseId,
    enabled: boolean
  ): void
  acknowledgeCredentialConfigured(phaseId: TranslationJobSetupPhaseId): void
  savePhaseCredential?(
    phaseId: TranslationJobSetupPhaseId,
    apiKey: string
  ): Promise<void>
  runValidation(): Promise<void>
  createJob(): Promise<void>
}

export type CreateTranslationJobSetupScreenController =
  () => TranslationJobSetupScreenControllerContract

export { createTranslationJobSetupRuntimeKey }
