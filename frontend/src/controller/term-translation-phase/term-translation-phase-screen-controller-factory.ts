import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import { TermTranslationPhasePresenter } from "@application/presenter/term-translation-phase"
import { TermTranslationPhaseStore } from "@application/store/term-translation-phase"
import { TermTranslationPhaseUseCase } from "@application/usecase/term-translation-phase"
import type {
  GetTermTranslationNextPhaseReadinessRequestDto,
  GetTermTranslationNextPhaseReadinessResponseDto,
  GetTermTranslationPhaseSummaryRequestDto,
  GetTermTranslationPhaseSummaryResponseDto,
  PauseTermTranslationPhaseRequestDto,
  PauseTermTranslationPhaseResponseDto,
  ResumeTermTranslationPhaseRequestDto,
  ResumeTermTranslationPhaseResponseDto,
  RetryTermTranslationPhaseRequestDto,
  RetryTermTranslationPhaseResponseDto,
  StartTermTranslationPhaseRequestDto,
  StartTermTranslationPhaseResponseDto
} from "@controller/wails/gateway-dto/term-translation-phase"

import { TermTranslationPhaseScreenController } from "./term-translation-phase-screen-controller"

type TermTranslationPhaseGatewayDtoCoverage = {
  getSummaryRequest: GetTermTranslationPhaseSummaryRequestDto
  getSummaryResponse: GetTermTranslationPhaseSummaryResponseDto
  startRequest: StartTermTranslationPhaseRequestDto
  startResponse: StartTermTranslationPhaseResponseDto
  pauseRequest: PauseTermTranslationPhaseRequestDto
  pauseResponse: PauseTermTranslationPhaseResponseDto
  resumeRequest: ResumeTermTranslationPhaseRequestDto
  resumeResponse: ResumeTermTranslationPhaseResponseDto
  retryRequest: RetryTermTranslationPhaseRequestDto
  retryResponse: RetryTermTranslationPhaseResponseDto
  nextPhaseRequest: GetTermTranslationNextPhaseReadinessRequestDto
  nextPhaseResponse: GetTermTranslationNextPhaseReadinessResponseDto
}

export function createTermTranslationPhaseScreenControllerFactory(
  gateway: TermTranslationPhaseGatewayContract | null
): CreateTermTranslationPhaseScreenController {
  let controller: TermTranslationPhaseScreenController | null = null

  const factory: CreateTermTranslationPhaseScreenController & {
    __dtoCoverage?: TermTranslationPhaseGatewayDtoCoverage
  } = () => {
    if (controller) {
      return controller
    }

    const store = new TermTranslationPhaseStore()
    const presenter = new TermTranslationPhasePresenter()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    controller = new TermTranslationPhaseScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }

  return factory
}
