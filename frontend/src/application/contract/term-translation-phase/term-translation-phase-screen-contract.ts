import type { TermTranslationPhaseScreenViewModel } from "./term-translation-phase-screen-types"

export type TermTranslationPhaseScreenViewModelListener = (
  viewModel: TermTranslationPhaseScreenViewModel
) => void

export interface TermTranslationPhaseScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: TermTranslationPhaseScreenViewModelListener): () => void
  getViewModel(): TermTranslationPhaseScreenViewModel
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?(searchQuery: string): Promise<void>
  setProcessingTargetPage?(page: number): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
}

export type CreateTermTranslationPhaseScreenController =
  () => TermTranslationPhaseScreenControllerContract
