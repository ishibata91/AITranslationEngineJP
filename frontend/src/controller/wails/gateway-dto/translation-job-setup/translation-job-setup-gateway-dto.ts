import type {
  CreateTranslationJobResponse,
  DeleteTranslationJobSetupInputRequest,
  DeleteTranslationJobSetupInputResponse,
  GetTranslationJobSetupSummaryRequest,
  ListTranslationJobSetupProviderModelsRequest,
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupPhaseRuntimeValidationSelection,
  TranslationJobSetupPhaseRuntimeSummary,
  TranslationJobSetupRuntimeSelection,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"

export type GetTranslationJobSetupOptionsResponseDto =
  TranslationJobSetupOptionsResponse

export type ListTranslationJobSetupProviderModelsRequestDto =
  ListTranslationJobSetupProviderModelsRequest
export type ListTranslationJobSetupProviderModelsResponseDto =
  ListTranslationJobSetupProviderModelsResponse

export type TranslationJobSetupPhaseRuntimeValidationSelectionDto = Omit<
  TranslationJobSetupPhaseRuntimeValidationSelection,
  "credentialRef"
>

export interface ValidateTranslationJobSetupRequestDto {
  inputSourceId: number
  runtime: TranslationJobSetupRuntimeSelection
  phaseRuntimeSelections?: TranslationJobSetupPhaseRuntimeValidationSelectionDto[]
}
export type ValidateTranslationJobSetupResponseDto =
  TranslationJobSetupValidationResponse

export interface CreateTranslationJobRequestDto {
  inputSourceId: number
  inputSource: string
  validationStatus: string
  validatedAt: string
  validationPassSlices: string[]
  runtime: TranslationJobSetupRuntimeSelection
  phaseRuntimeSelections?: TranslationJobSetupPhaseRuntimeValidationSelectionDto[]
}

export type TranslationJobSetupPhaseRuntimeSummaryDto = Omit<
  TranslationJobSetupPhaseRuntimeSummary,
  "credentialRef"
>

export type CreateTranslationJobResponseDto = Omit<
  CreateTranslationJobResponse,
  "phaseRuntimeSummaries"
> & {
  phaseRuntimeSummaries?: TranslationJobSetupPhaseRuntimeSummaryDto[]
}
export type DeleteTranslationJobSetupInputRequestDto =
  DeleteTranslationJobSetupInputRequest
export type DeleteTranslationJobSetupInputResponseDto =
  DeleteTranslationJobSetupInputResponse

export type GetTranslationJobSetupSummaryRequestDto =
  GetTranslationJobSetupSummaryRequest
export type GetTranslationJobSetupSummaryResponseDto = Omit<
  TranslationJobSetupSummaryResponse,
  "phaseRuntimeSummaries"
> & {
  phaseRuntimeSummaries?: TranslationJobSetupPhaseRuntimeSummaryDto[]
}
