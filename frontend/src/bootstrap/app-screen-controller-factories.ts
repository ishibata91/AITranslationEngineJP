import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
import type { CreateProviderSettingsScreenController } from "@application/contract/provider-settings"
import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
import type { CreateTranslationJobSetupScreenController } from "@application/contract/translation-job-setup"
import type { CreateTranslationOutputArtifactScreenController } from "@application/contract/translation-output-artifact"
import {
  createNoopFrontendDiagnosticLogger,
  createPinoFrontendDiagnosticLogger
} from "@application/diagnostic"
import { createBodyTranslationPhaseScreenControllerFactory } from "@controller/body-translation-phase"
import { createMasterDictionaryScreenControllerFactory } from "@controller/master-dictionary"
import { createMasterPersonaScreenControllerFactory } from "@controller/master-persona"
import { createPersonaGenerationPhaseScreenControllerFactory } from "@controller/persona-generation-phase"
import { createProviderSettingsScreenControllerFactory } from "@controller/provider-settings"
import { createTermTranslationPhaseScreenControllerFactory } from "@controller/term-translation-phase"
import { createTranslationJobManagementScreenControllerFactory } from "@controller/translation-job-management"
import { createTranslationInputScreenControllerFactory } from "@controller/translation-input"
import { createTranslationJobSetupScreenControllerFactory } from "@controller/translation-job-setup"
import { createTranslationOutputArtifactScreenControllerFactory } from "@controller/translation-output-artifact"
import { createBodyTranslationPhaseGateway } from "@controller/wails/body-translation-phase.gateway"
import { createMasterDictionaryGateway } from "@controller/wails/master-dictionary.gateway"
import { createMasterPersonaGateway } from "@controller/wails/master-persona.gateway"
import { createPersonaGenerationPhaseGateway } from "@controller/wails/persona-generation-phase.gateway"
import { createProviderSettingsGateway } from "@controller/wails/provider-settings.gateway"
import { createTermTranslationPhaseGateway } from "@controller/wails/term-translation-phase.gateway"
import { createTranslationJobManagementGateway } from "@controller/wails/translation-job-management.gateway"
import { createTranslationInputGateway } from "@controller/wails/translation-input.gateway"
import { createTranslationJobSetupGateway } from "@controller/wails/translation-job-setup.gateway"
import { createTranslationOutputArtifactGateway } from "@controller/wails/translation-output-artifact.gateway"

type FrontendDiagnosticLogger = ReturnType<
  typeof createNoopFrontendDiagnosticLogger
>

interface AppScreenControllerFactories {
  diagnosticLogger: FrontendDiagnosticLogger
  createBodyTranslationPhaseScreenController: CreateBodyTranslationPhaseScreenController
  createMasterDictionaryScreenController: CreateMasterDictionaryScreenController
  createMasterPersonaScreenController: CreateMasterPersonaScreenController
  createPersonaGenerationPhaseScreenController: CreatePersonaGenerationPhaseScreenController
  createProviderSettingsScreenController: CreateProviderSettingsScreenController
  createTermTranslationPhaseScreenController: CreateTermTranslationPhaseScreenController
  createTranslationJobManagementScreenController: CreateTranslationJobManagementScreenController
  createTranslationInputScreenController: CreateTranslationInputScreenController
  createTranslationJobSetupScreenController: CreateTranslationJobSetupScreenController
  createTranslationOutputArtifactScreenController: CreateTranslationOutputArtifactScreenController
}

export function createProductionAppFactories(): AppScreenControllerFactories {
  const masterDictionaryGateway = createMasterDictionaryGateway()
  const termTranslationPhaseGateway = createTermTranslationPhaseGateway()
  const bodyTranslationPhaseGateway = createBodyTranslationPhaseGateway()
  const personaGenerationPhaseGateway = createPersonaGenerationPhaseGateway()
  const providerSettingsGateway = createProviderSettingsGateway()
  const translationInputGateway = createTranslationInputGateway()
  const translationJobSetupGateway = createTranslationJobSetupGateway()
  const translationJobManagementGateway =
    createTranslationJobManagementGateway()
  const translationOutputArtifactGateway =
    createTranslationOutputArtifactGateway()
  const masterPersonaGateway = createMasterPersonaGateway()
  const diagnosticLogger = createPinoFrontendDiagnosticLogger({
    source: "frontend"
  })

  return {
    diagnosticLogger,
    createBodyTranslationPhaseScreenController:
      createBodyTranslationPhaseScreenControllerFactory(
        bodyTranslationPhaseGateway
      ),
    createMasterDictionaryScreenController:
      createMasterDictionaryScreenControllerFactory(
        masterDictionaryGateway,
        diagnosticLogger
      ),
    createMasterPersonaScreenController:
      createMasterPersonaScreenControllerFactory(masterPersonaGateway),
    createPersonaGenerationPhaseScreenController:
      createPersonaGenerationPhaseScreenControllerFactory(
        personaGenerationPhaseGateway
      ),
    createProviderSettingsScreenController:
      createProviderSettingsScreenControllerFactory(providerSettingsGateway),
    createTermTranslationPhaseScreenController:
      createTermTranslationPhaseScreenControllerFactory(
        termTranslationPhaseGateway
      ),
    createTranslationJobManagementScreenController:
      createTranslationJobManagementScreenControllerFactory(
        translationJobManagementGateway
      ),
    createTranslationInputScreenController:
      createTranslationInputScreenControllerFactory(translationInputGateway),
    createTranslationJobSetupScreenController:
      createTranslationJobSetupScreenControllerFactory(
        translationJobSetupGateway
      ),
    createTranslationOutputArtifactScreenController:
      createTranslationOutputArtifactScreenControllerFactory(
        translationOutputArtifactGateway
      )
  }
}
