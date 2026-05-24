<script lang="ts">
  import type { CreateBodyTranslationPhaseScreenController } from "@application/contract/body-translation-phase"
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import type { CreatePersonaGenerationPhaseScreenController } from "@application/contract/persona-generation-phase"
  import type { CreateProviderSettingsScreenController } from "@application/contract/provider-settings"
  import type { CreateTermTranslationPhaseScreenController } from "@application/contract/term-translation-phase"
  import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
  import type { CreateTranslationOutputArtifactScreenController } from "@application/contract/translation-output-artifact"
  import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaScreenControllerFactory } from "@controller/master-persona"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createProviderSettingsScreenControllerFactory } from "@controller/provider-settings"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationOutputArtifactScreenControllerFactory } from "@controller/translation-output-artifact"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationJobManagementScreenControllerFactory } from "@controller/translation-job-management"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputScreenControllerFactory } from "@controller/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaGateway } from "@controller/wails/master-persona.gateway"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationJobManagementGateway } from "@controller/wails/translation-job-management.gateway"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputGateway } from "@controller/wails/translation-input.gateway"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createBodyTranslationPhaseScreenController?: CreateBodyTranslationPhaseScreenController | null
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController?: CreateMasterPersonaScreenController | null
    createPersonaGenerationPhaseScreenController?: CreatePersonaGenerationPhaseScreenController | null
    createProviderSettingsScreenController?: CreateProviderSettingsScreenController | null
    createTermTranslationPhaseScreenController?: CreateTermTranslationPhaseScreenController | null
    createTranslationJobManagementScreenController?: CreateTranslationJobManagementScreenController | null
    createTranslationOutputArtifactScreenController?: CreateTranslationOutputArtifactScreenController | null
    createTranslationInputScreenController?: CreateTranslationInputScreenController | null
  }

  let {
    createBodyTranslationPhaseScreenController = null,
    createMasterDictionaryScreenController = null,
    createMasterPersonaScreenController = null,
    createPersonaGenerationPhaseScreenController = null,
    createProviderSettingsScreenController = null,
    createTermTranslationPhaseScreenController = null,
    createTranslationJobManagementScreenController = null,
    createTranslationOutputArtifactScreenController = null,
    createTranslationInputScreenController = null
  }: Props = $props()

  const shellState = createShellState()

  function resolveMasterPersonaScreenControllerFactory(): CreateMasterPersonaScreenController {
    return (
      createMasterPersonaScreenController ??
      createMasterPersonaScreenControllerFactory(createMasterPersonaGateway())
    )
  }

  function resolveProviderSettingsScreenControllerFactory(): CreateProviderSettingsScreenController {
    return (
      createProviderSettingsScreenController ??
      createProviderSettingsScreenControllerFactory(null)
    )
  }

  function resolveTranslationInputScreenControllerFactory(): CreateTranslationInputScreenController {
    return (
      createTranslationInputScreenController ??
      createTranslationInputScreenControllerFactory(
        createTranslationInputGateway()
      )
    )
  }

  function resolveTranslationJobManagementScreenControllerFactory(): CreateTranslationJobManagementScreenController {
    return (
      createTranslationJobManagementScreenController ??
      createTranslationJobManagementScreenControllerFactory(
        createTranslationJobManagementGateway()
      )
    )
  }

  function resolveTranslationOutputArtifactScreenControllerFactory(): CreateTranslationOutputArtifactScreenController {
    return (
      createTranslationOutputArtifactScreenController ??
      createTranslationOutputArtifactScreenControllerFactory(null)
    )
  }
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  defaultTranslationManagementViewId={shellState.defaultTranslationManagementViewId}
  {createBodyTranslationPhaseScreenController}
  {createMasterDictionaryScreenController}
  createMasterPersonaScreenController={resolveMasterPersonaScreenControllerFactory()}
  {createPersonaGenerationPhaseScreenController}
  createProviderSettingsScreenController={resolveProviderSettingsScreenControllerFactory()}
  {createTermTranslationPhaseScreenController}
  createTranslationJobManagementScreenController={resolveTranslationJobManagementScreenControllerFactory()}
  createTranslationOutputArtifactScreenController={resolveTranslationOutputArtifactScreenControllerFactory()}
  createTranslationInputScreenController={resolveTranslationInputScreenControllerFactory()}
  routes={shellState.routes}
/>
