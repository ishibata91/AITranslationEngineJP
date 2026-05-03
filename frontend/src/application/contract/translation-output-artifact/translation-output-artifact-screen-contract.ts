import type { TranslationOutputArtifactScreenViewModel } from "./translation-output-artifact-screen-types"

export type TranslationOutputArtifactScreenViewModelListener = (
  viewModel: TranslationOutputArtifactScreenViewModel
) => void

export interface TranslationOutputArtifactScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(
    listener: TranslationOutputArtifactScreenViewModelListener
  ): () => void
  getViewModel(): TranslationOutputArtifactScreenViewModel
  setJobId(jobId: number | null): Promise<void>
  setArtifactId(artifactId: number | null): Promise<void>
  setTargetGame(targetGame: string): void
  setOutputPath(outputPath: string): void
  refresh(): Promise<void>
  generateArtifact(): Promise<void>
  regenerateArtifact(): Promise<void>
}

export type CreateTranslationOutputArtifactScreenController =
  () => TranslationOutputArtifactScreenControllerContract
