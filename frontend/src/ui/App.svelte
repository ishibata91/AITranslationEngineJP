<script lang="ts">
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import type { CreateTranslationJobSetupScreenController } from "@application/contract/translation-job-setup"
  import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaScreenControllerFactory } from "@controller/master-persona"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationJobSetupScreenControllerFactory } from "@controller/translation-job-setup"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputScreenControllerFactory } from "@controller/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaGateway } from "@controller/wails/master-persona.gateway"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationJobSetupGateway } from "@controller/wails/translation-job-setup.gateway"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputGateway } from "@controller/wails/translation-input.gateway"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController?: CreateMasterPersonaScreenController | null
    createTranslationJobSetupScreenController?: CreateTranslationJobSetupScreenController | null
    createTranslationInputScreenController?: CreateTranslationInputScreenController | null
  }

  let {
    createMasterDictionaryScreenController = null,
    createMasterPersonaScreenController = null,
    createTranslationJobSetupScreenController = null,
    createTranslationInputScreenController = null
  }: Props = $props()

  const shellState = createShellState()

  function resolveMasterPersonaScreenControllerFactory(): CreateMasterPersonaScreenController {
    return (
      createMasterPersonaScreenController ??
      createMasterPersonaScreenControllerFactory(createMasterPersonaGateway())
    )
  }

  function resolveTranslationInputScreenControllerFactory(): CreateTranslationInputScreenController {
    return (
      createTranslationInputScreenController ??
      createTranslationInputScreenControllerFactory(createTranslationInputGateway())
    )
  }

  function resolveTranslationJobSetupScreenControllerFactory(): CreateTranslationJobSetupScreenController {
    return (
      createTranslationJobSetupScreenController ??
      createTranslationJobSetupScreenControllerFactory(
        createTranslationJobSetupGateway()
      )
    )
  }
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  defaultTranslationManagementViewId={shellState.defaultTranslationManagementViewId}
  {createMasterDictionaryScreenController}
  createMasterPersonaScreenController={resolveMasterPersonaScreenControllerFactory()}
  createTranslationJobSetupScreenController={resolveTranslationJobSetupScreenControllerFactory()}
  createTranslationInputScreenController={resolveTranslationInputScreenControllerFactory()}
  routes={shellState.routes}
  translationManagementViews={shellState.translationManagementViews}
/>
