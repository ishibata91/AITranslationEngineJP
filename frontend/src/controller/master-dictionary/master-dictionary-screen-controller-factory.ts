import type {
  CreateMasterDictionaryScreenController,
  MasterDictionaryScreenControllerContract
} from "@application/contract/master-dictionary"
import type {
  RuntimeImportCompletedPayload,
  RuntimeImportProgressPayload
} from "@application/contract/master-dictionary/master-dictionary-screen-types"
import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
import { MasterDictionaryPresenter } from "@application/presenter/master-dictionary"
import { MasterDictionaryStore } from "@application/store/master-dictionary"
import { MasterDictionaryUseCase } from "@application/usecase/master-dictionary"
import { MasterDictionaryRuntimeEventAdapter } from "@controller/runtime/master-dictionary"

import { MasterDictionaryScreenController } from "./master-dictionary-screen-controller"

export function createMasterDictionaryScreenControllerFactory(
  gateway: MasterDictionaryGatewayContract | null
): CreateMasterDictionaryScreenController {
  return (): MasterDictionaryScreenControllerContract => {
    const store = new MasterDictionaryStore()
    const presenter = new MasterDictionaryPresenter()
    const useCase = new MasterDictionaryUseCase(gateway, store)
    const runtimeEventAdapter = new MasterDictionaryRuntimeEventAdapter({
      onImportProgress: (payload: RuntimeImportProgressPayload) => {
        useCase.handleImportProgress(payload)
      },
      onImportCompleted: (payload: RuntimeImportCompletedPayload) => {
        void useCase.handleImportCompleted(payload)
      }
    })

    return new MasterDictionaryScreenController({
      isGatewayConnected: gateway !== null,
      store,
      presenter,
      useCase,
      runtimeEventAdapter
    })
  }
}
