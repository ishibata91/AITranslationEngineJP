import type { PersonaGenerationPhaseScreenViewModel } from "./persona-generation-phase-screen-types"

export type PersonaGenerationPhaseScreenViewModelListener = (
  viewModel: PersonaGenerationPhaseScreenViewModel
) => void

export interface PersonaGenerationPhaseScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: PersonaGenerationPhaseScreenViewModelListener): () => void
  getViewModel(): PersonaGenerationPhaseScreenViewModel
  setJobId(jobId: number | null): Promise<void>
  setProcessingTargetSearchQuery?(searchQuery: string): Promise<void>
  setProcessingTargetPage?(page: number): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkBodyReadiness(): Promise<void>
  startBodyPhase(): Promise<void>
  saveAISettings?: (request: {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }) => Promise<void>
  refreshModelList?(provider: string): Promise<void>
}

export type CreatePersonaGenerationPhaseScreenController =
  () => PersonaGenerationPhaseScreenControllerContract
