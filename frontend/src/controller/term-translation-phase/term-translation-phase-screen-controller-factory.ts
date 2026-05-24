import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import { TermTranslationPhasePresenter } from "@application/presenter/term-translation-phase"
import { TermTranslationPhaseStore } from "@application/store/term-translation-phase"
import { TermTranslationPhaseUseCase } from "@application/usecase/term-translation-phase"

import { TermTranslationPhaseScreenController } from "./term-translation-phase-screen-controller"

export function createTermTranslationPhaseScreenControllerFactory(
  gateway: TermTranslationPhaseGatewayContract | null
): CreateTermTranslationPhaseScreenController {
  let controller: TermTranslationPhaseScreenController | null = null

  const factory: CreateTermTranslationPhaseScreenController = () => {
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
