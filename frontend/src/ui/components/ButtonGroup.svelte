<script lang="ts">
  import type { Snippet } from "svelte"

  type Align = "start" | "end" | "stretch"

  interface Props {
    ariaLabel?: string
    align?: Align
    children?: Snippet
  }

  let {
    ariaLabel = "操作",
    align = "end",
    children
  }: Props = $props()
</script>

<div class={`button-group button-group-${align}`} aria-label={ariaLabel} role="group">
  {@render children?.()}
</div>

<style>
  .button-group {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.65rem;
    min-width: 0;
  }

  .button-group-start {
    justify-content: flex-start;
  }

  .button-group-end {
    justify-content: flex-end;
  }

  .button-group-stretch {
    align-items: stretch;
  }

  .button-group-stretch :global(button) {
    flex: 1 1 10rem;
  }
</style>
