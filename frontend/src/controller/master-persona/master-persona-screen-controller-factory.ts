import type {
  CreateMasterPersonaScreenController,
  MasterPersonaScreenControllerContract
} from "@application/contract/master-persona"
import type { MasterPersonaGatewayContract } from "@application/gateway-contract/master-persona"
import { MasterPersonaPresenter } from "@application/presenter/master-persona"
import { MasterPersonaStore } from "@application/store/master-persona"
import { MasterPersonaUseCase } from "@application/usecase/master-persona"
// eslint-disable-next-line local/enforce-layer-boundaries
import { MasterPersonaRuntimePollingAdapter } from "@controller/runtime/master-persona"

import { MasterPersonaScreenController } from "./master-persona-screen-controller"

export function createMasterPersonaScreenControllerFactory(
  gateway: MasterPersonaGatewayContract | null
): CreateMasterPersonaScreenController {
  return (): MasterPersonaScreenControllerContract => {
    const store = new MasterPersonaStore()
    const presenter = new MasterPersonaPresenter()
    const useCase = new MasterPersonaUseCase(gateway, store)
    const runtimePollingAdapter = new MasterPersonaRuntimePollingAdapter()

    return new MasterPersonaScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase,
      runtimePollingAdapter
    })
  }
}
