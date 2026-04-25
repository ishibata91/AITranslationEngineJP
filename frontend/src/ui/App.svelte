<script lang="ts">
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  import type { CreateTranslationInputScreenController } from "@application/contract/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaScreenControllerFactory } from "@controller/master-persona"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputScreenControllerFactory } from "@controller/translation-input"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaGateway } from "@controller/wails/master-persona.gateway"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createTranslationInputGateway } from "@controller/wails/translation-input.gateway"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController?: CreateMasterPersonaScreenController | null
    createTranslationInputScreenController?: CreateTranslationInputScreenController | null
  }

  let {
    createMasterDictionaryScreenController = null,
    createMasterPersonaScreenController = null,
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
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  {createMasterDictionaryScreenController}
  createMasterPersonaScreenController={resolveMasterPersonaScreenControllerFactory()}
  createTranslationInputScreenController={resolveTranslationInputScreenControllerFactory()}
  routes={shellState.routes}
/>
