<script lang="ts">
  import type { Snippet } from "svelte"

  interface Props {
    ariaLabel: string
    children: Snippet
  }

  let { ariaLabel, children }: Props = $props()
</script>

<div class="overflow-x-auto">
  <table class="table table-zebra" aria-label={ariaLabel}>
    {@render children()}
  </table>
</div>
