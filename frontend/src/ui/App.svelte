<script lang="ts">
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import type { CreateMasterPersonaScreenController } from "@application/contract/master-persona"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaScreenControllerFactory } from "@controller/master-persona"
  // eslint-disable-next-line local/enforce-layer-boundaries
  import { createMasterPersonaGateway } from "@controller/wails/master-persona.gateway"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
    createMasterPersonaScreenController?: CreateMasterPersonaScreenController | null
  }

  let {
    createMasterDictionaryScreenController = null,
    createMasterPersonaScreenController = null
  }: Props = $props()

  const shellState = createShellState()

  function resolveMasterPersonaScreenControllerFactory(): CreateMasterPersonaScreenController {
    return (
      createMasterPersonaScreenController ??
      createMasterPersonaScreenControllerFactory(createMasterPersonaGateway())
    )
  }
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  {createMasterDictionaryScreenController}
  createMasterPersonaScreenController={resolveMasterPersonaScreenControllerFactory()}
  routes={shellState.routes}
/>
