import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
import type { PersonaGenerationPhaseGatewayContract } from "@application/gateway-contract/persona-generation-phase"
import { PersonaGenerationPhasePresenter } from "@application/presenter/persona-generation-phase"
import { PersonaGenerationPhaseStore } from "@application/store/persona-generation-phase"
import { PersonaGenerationPhaseUseCase } from "@application/usecase/persona-generation-phase"

import { PersonaGenerationPhaseScreenController } from "./persona-generation-phase-screen-controller"

export function createPersonaGenerationPhaseScreenControllerFactory(
  gateway: PersonaGenerationPhaseGatewayContract | null
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
      useCase
    })

    return controller
  }

  return factory
}
