import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
import type { BodyTranslationPhaseGatewayContract } from "@application/gateway-contract/body-translation-phase"
import { BodyTranslationPhasePresenter } from "@application/presenter/body-translation-phase"
import { BodyTranslationPhaseStore } from "@application/store/body-translation-phase"
import { BodyTranslationPhaseUseCase } from "@application/usecase/body-translation-phase"

import { BodyTranslationPhaseScreenController } from "./body-translation-phase-screen-controller"

export function createBodyTranslationPhaseScreenControllerFactory(
  gateway: BodyTranslationPhaseGatewayContract | null
): CreateBodyTranslationPhaseScreenController {
  let controller: BodyTranslationPhaseScreenController | null = null

  const factory: CreateBodyTranslationPhaseScreenController = () => {
    if (controller) {
      return controller
    }

    const store = new BodyTranslationPhaseStore()
    const presenter = new BodyTranslationPhasePresenter()
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    controller = new BodyTranslationPhaseScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }

  return factory
}
