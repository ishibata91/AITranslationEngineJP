<script lang="ts">
	import AppShell from "./views/AppShell.svelte";
	import { createShellState } from "./stores/shell-state";

	const shellState = createShellState();
</script>

<AppShell title={shellState.title} />
