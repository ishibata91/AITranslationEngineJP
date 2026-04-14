<script lang="ts">
  import type { CreateMasterDictionaryScreenController } from "@application/contract/master-dictionary"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    createMasterDictionaryScreenController?: CreateMasterDictionaryScreenController | null
  }

  let { createMasterDictionaryScreenController = null }: Props = $props()

  const shellState = createShellState()
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  {createMasterDictionaryScreenController}
  routes={shellState.routes}
/>
