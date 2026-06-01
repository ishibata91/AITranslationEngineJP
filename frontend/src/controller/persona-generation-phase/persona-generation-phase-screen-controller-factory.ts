import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import type { PersonaGenerationPhaseGatewayContract } from "@application/gateway-contract/persona-generation-phase"
import { PersonaGenerationPhasePresenter } from "@application/presenter/persona-generation-phase"
import { PersonaGenerationPhaseStore } from "@application/store/persona-generation-phase"
import { PersonaGenerationPhaseUseCase } from "@application/usecase/persona-generation-phase"

import { PersonaGenerationPhaseScreenController } from "./persona-generation-phase-screen-controller"

export function createPersonaGenerationPhaseScreenControllerFactory(
  gateway: PersonaGenerationPhaseGatewayContract | null,
  providerSettingsGateway?: ProviderSettingsGatewayContract | null
): CreatePersonaGenerationPhaseScreenController {
  let controller: PersonaGenerationPhaseScreenController | null = null

  const factory: CreatePersonaGenerationPhaseScreenController = () => {
    if (controller) {
      return controller
    }

    const store = new PersonaGenerationPhaseStore()
    const presenter = new PersonaGenerationPhasePresenter()
    const useCase = new PersonaGenerationPhaseUseCase(gateway, store)

    controller = new PersonaGenerationPhaseScreenController({
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
