import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
import type { BodyTranslationPhaseGatewayContract } from "@application/gateway-contract/body-translation-phase"
import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import { BodyTranslationPhasePresenter } from "@application/presenter/body-translation-phase"
import { BodyTranslationPhaseStore } from "@application/store/body-translation-phase"
import { BodyTranslationPhaseUseCase } from "@application/usecase/body-translation-phase"

import { BodyTranslationPhaseScreenController } from "./body-translation-phase-screen-controller"

export function createBodyTranslationPhaseScreenControllerFactory(
  gateway: BodyTranslationPhaseGatewayContract | null,
  providerSettingsGateway?: ProviderSettingsGatewayContract | null
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
      useCase,
      gateway: gateway ?? null
    })

    if (providerSettingsGateway) {
      providerSettingsGateway.ListProviderSettings({}).then((response) => {
        const providers = response.providers.map((p) => ({
          value: p.providerId,
          label: p.label
        }))
        controller?.setAvailableProviders(providers)
      }).catch(() => {
        // provider 一覧の取得失敗時はフォールバック（空配列）のまま継続する
      })
    }

    return controller
  }

  return factory
}
