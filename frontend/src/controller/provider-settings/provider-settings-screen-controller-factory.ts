import type { CreateProviderSettingsScreenController } from "@application/contract/provider-settings"
import type { ProviderSettingsGatewayContract } from "@application/gateway-contract/provider-settings"
import { ProviderSettingsPresenter } from "@application/presenter/provider-settings"
import { ProviderSettingsStore } from "@application/store/provider-settings"
import { ProviderSettingsUseCase } from "@application/usecase/provider-settings"

import { ProviderSettingsScreenController } from "./provider-settings-screen-controller"

export function createProviderSettingsScreenControllerFactory(
  gateway: ProviderSettingsGatewayContract | null
): CreateProviderSettingsScreenController {
  let controller: ProviderSettingsScreenController | null = null

  return () => {
    if (controller) {
      return controller
    }

    const store = new ProviderSettingsStore()
    const presenter = new ProviderSettingsPresenter()
    const useCase = new ProviderSettingsUseCase(gateway, store)

    controller = new ProviderSettingsScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase
    })

    return controller
  }
}
