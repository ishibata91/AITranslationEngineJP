import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
import type { TranslationJobManagementGatewayContract } from "@application/gateway-contract/translation-job-management"
import { TranslationJobManagementPresenter } from "@application/presenter/translation-job-management"
import { TranslationJobManagementStore } from "@application/store/translation-job-management"
import { TranslationJobManagementUseCase } from "@application/usecase/translation-job-management"

import { TranslationJobManagementScreenController } from "./translation-job-management-screen-controller"

export function createTranslationJobManagementScreenControllerFactory(
  gateway: TranslationJobManagementGatewayContract | null
): CreateTranslationJobManagementScreenController {
  let controller: TranslationJobManagementScreenController | null = null

  return () => {
    if (controller) {
      return controller
    }

    const store = new TranslationJobManagementStore()
    const presenter = new TranslationJobManagementPresenter()
    const useCase = new TranslationJobManagementUseCase(gateway, store)

    controller = new TranslationJobManagementScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }
}
