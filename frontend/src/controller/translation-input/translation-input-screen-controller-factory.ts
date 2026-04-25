import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
import type { TranslationInputGatewayContract } from "@application/gateway-contract/translation-input"
import { TranslationInputPresenter } from "@application/presenter/translation-input"
import { TranslationInputStore } from "@application/store/translation-input"
import { TranslationInputUseCase } from "@application/usecase/translation-input"

import { TranslationInputScreenController } from "./translation-input-screen-controller"

export function createTranslationInputScreenControllerFactory(
  gateway: TranslationInputGatewayContract | null
): CreateTranslationInputScreenController {
  let controller: TranslationInputScreenController | null = null

  return () => {
    if (controller) {
      return controller
    }

    const store = new TranslationInputStore()
    const presenter = new TranslationInputPresenter()
    const useCase = new TranslationInputUseCase(gateway, store)

    controller = new TranslationInputScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }
}