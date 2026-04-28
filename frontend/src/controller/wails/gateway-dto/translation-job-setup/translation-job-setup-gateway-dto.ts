import type {
  CreateTranslationJobRequest,
  CreateTranslationJobResponse,
  GetTranslationJobSetupSummaryRequest,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse,
  ValidateTranslationJobSetupRequest
} from "@application/gateway-contract/translation-job-setup"

export type GetTranslationJobSetupOptionsResponseDto =
  TranslationJobSetupOptionsResponse

export type ValidateTranslationJobSetupRequestDto =
  ValidateTranslationJobSetupRequest
export type ValidateTranslationJobSetupResponseDto =
  TranslationJobSetupValidationResponse

export type CreateTranslationJobRequestDto = CreateTranslationJobRequest
export type CreateTranslationJobResponseDto = CreateTranslationJobResponse

export type GetTranslationJobSetupSummaryRequestDto =
  GetTranslationJobSetupSummaryRequest
export type GetTranslationJobSetupSummaryResponseDto =
  TranslationJobSetupSummaryResponse