import type { BodyTranslationPhaseScreenViewModel } from "./body-translation-phase-screen-types"

export type BodyTranslationPhaseScreenViewModelListener = (
  viewModel: BodyTranslationPhaseScreenViewModel
) => void

export interface BodyTranslationPhaseScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: BodyTranslationPhaseScreenViewModelListener): () => void
  getViewModel(): BodyTranslationPhaseScreenViewModel
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?(
    searchQuery: string,
    phase?: string
  ): Promise<void>
  setProcessingTargetPage?(page: number, phase?: string): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkOutputReadiness(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
  refreshModelList?(provider: string): Promise<void>
}

export type CreateBodyTranslationPhaseScreenController =
  () => BodyTranslationPhaseScreenControllerContract
