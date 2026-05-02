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
  refresh(): Promise<void>
  startPhase(): Promise<void>
  pausePhase(): Promise<void>
  resumePhase(): Promise<void>
  retryPhase(): Promise<void>
  cancelPhase(): Promise<void>
  checkBodyReadiness(): Promise<void>
  startBodyPhase(): Promise<void>
}

export type CreatePersonaGenerationPhaseScreenController =
  () => PersonaGenerationPhaseScreenControllerContract
