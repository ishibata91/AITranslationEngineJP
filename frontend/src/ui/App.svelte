<script lang="ts">
  import AppShell from "@ui/views/AppShell.svelte"
  import { createShellState } from "@ui/stores/shell-state"

  const shellState = createShellState()
</script>

<AppShell defaultRouteId={shellState.defaultRouteId} routes={shellState.routes} />
