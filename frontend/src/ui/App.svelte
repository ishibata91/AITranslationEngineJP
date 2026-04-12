<script lang="ts">
  import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
  import { createShellState } from "@ui/stores/shell-state"
  import AppShell from "@ui/views/AppShell.svelte"

  interface Props {
    masterDictionaryGateway?: MasterDictionaryGatewayContract | null
  }

  let { masterDictionaryGateway = null }: Props = $props()

  const shellState = createShellState()
</script>

<AppShell
  defaultRouteId={shellState.defaultRouteId}
  masterDictionaryGateway={masterDictionaryGateway}
  routes={shellState.routes}
/>
