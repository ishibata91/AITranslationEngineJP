import type { CreateTranslationOutputArtifactScreenController } from "@application/contract/translation-output-artifact"
import type { TranslationOutputArtifactGatewayContract } from "@application/gateway-contract/translation-output-artifact"
import { TranslationOutputArtifactPresenter } from "@application/presenter/translation-output-artifact"
import { TranslationOutputArtifactStore } from "@application/store/translation-output-artifact"
import { TranslationOutputArtifactUseCase } from "@application/usecase/translation-output-artifact"

import { TranslationOutputArtifactScreenController } from "./translation-output-artifact-screen-controller"

export function createTranslationOutputArtifactScreenControllerFactory(
  gateway: TranslationOutputArtifactGatewayContract | null
): CreateTranslationOutputArtifactScreenController {
  let controller: TranslationOutputArtifactScreenController | null = null

  return () => {
    if (controller) {
      return controller
    }

    const store = new TranslationOutputArtifactStore()
    const presenter = new TranslationOutputArtifactPresenter()
    const useCase = new TranslationOutputArtifactUseCase(gateway, store)

    controller = new TranslationOutputArtifactScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }
}
