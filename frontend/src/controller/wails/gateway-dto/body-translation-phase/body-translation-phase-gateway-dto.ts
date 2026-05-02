import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseSummaryResponse,
  CancelBodyTranslationPhaseRequest,
  GetBodyTranslationOutputReadinessRequest,
  GetBodyTranslationPhaseSummaryRequest,
  PauseBodyTranslationPhaseRequest,
  ResumeBodyTranslationPhaseRequest,
  RetryBodyTranslationPhaseRequest,
  StartBodyTranslationPhaseRequest
} from "@application/gateway-contract/body-translation-phase"

export type GetBodyTranslationPhaseSummaryRequestDto =
  GetBodyTranslationPhaseSummaryRequest
export type GetBodyTranslationPhaseSummaryResponseDto =
  BodyTranslationPhaseSummaryResponse

export type StartBodyTranslationPhaseRequestDto =
  StartBodyTranslationPhaseRequest
export type StartBodyTranslationPhaseResponseDto =
  BodyTranslationPhaseCommandResponse

export type PauseBodyTranslationPhaseRequestDto =
  PauseBodyTranslationPhaseRequest
export type PauseBodyTranslationPhaseResponseDto =
  BodyTranslationPhaseCommandResponse

export type ResumeBodyTranslationPhaseRequestDto =
  ResumeBodyTranslationPhaseRequest
export type ResumeBodyTranslationPhaseResponseDto =
  BodyTranslationPhaseCommandResponse

export type RetryBodyTranslationPhaseRequestDto =
  RetryBodyTranslationPhaseRequest
export type RetryBodyTranslationPhaseResponseDto =
  BodyTranslationPhaseCommandResponse

export type CancelBodyTranslationPhaseRequestDto =
  CancelBodyTranslationPhaseRequest
export type CancelBodyTranslationPhaseResponseDto =
  BodyTranslationPhaseCommandResponse

export type GetBodyTranslationOutputReadinessRequestDto =
  GetBodyTranslationOutputReadinessRequest
export type GetBodyTranslationOutputReadinessResponseDto =
  BodyTranslationOutputReadinessResponse
