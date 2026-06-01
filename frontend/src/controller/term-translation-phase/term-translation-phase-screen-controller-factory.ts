import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import type { TermTranslationPhaseGatewayContract } from "@application/gateway-contract/term-translation-phase"
import { TermTranslationPhasePresenter } from "@application/presenter/term-translation-phase"
import { TermTranslationPhaseStore } from "@application/store/term-translation-phase"
import { TermTranslationPhaseUseCase } from "@application/usecase/term-translation-phase"

import { TermTranslationPhaseScreenController } from "./term-translation-phase-screen-controller"

export function createTermTranslationPhaseScreenControllerFactory(
  gateway: TermTranslationPhaseGatewayContract | null,
  providerSettingsGateway?: ProviderSettingsGatewayContract | null
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
