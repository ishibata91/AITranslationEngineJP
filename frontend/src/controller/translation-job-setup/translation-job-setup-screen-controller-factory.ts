import type { CreateTranslationJobSetupScreenController } from "@application/contract/translation-job-setup"
import type { TranslationJobSetupGatewayContract } from "@application/gateway-contract/translation-job-setup"
import { TranslationJobSetupPresenter } from "@application/presenter/translation-job-setup"
import { TranslationJobSetupStore } from "@application/store/translation-job-setup"
import { TranslationJobSetupUseCase } from "@application/usecase/translation-job-setup"

import { TranslationJobSetupScreenController } from "./translation-job-setup-screen-controller"

export function createTranslationJobSetupScreenControllerFactory(
  gateway: TranslationJobSetupGatewayContract | null
): CreateTranslationJobSetupScreenController {
  let controller: TranslationJobSetupScreenController | null = null

  return () => {
    if (controller) {
      return controller
    }

    const store = new TranslationJobSetupStore()
    const presenter = new TranslationJobSetupPresenter()
    const useCase = new TranslationJobSetupUseCase(gateway, store)

    controller = new TranslationJobSetupScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }
}