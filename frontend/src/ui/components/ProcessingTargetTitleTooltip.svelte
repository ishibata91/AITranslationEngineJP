<script lang="ts">
  import FloatingTooltipTrigger from "./FloatingTooltipTrigger.svelte"

  interface Props {
    text: string
    displayText?: string
    tooltipId: string
    secondary?: boolean
  }

  let { text, displayText, tooltipId, secondary = false }: Props = $props()
</script>

<FloatingTooltipTrigger
  {secondary}
  showWhenTruncatedOnly
  {tooltipId}
  triggerClass="target-title-part-text"
  triggerAriaLabel={text}
  triggerText={displayText ?? text}
  tooltipText={text}
/>

<style>
  :global(.target-title .tooltip-trigger-shell) {
    flex: 1 1 0;
    max-width: 100%;
  }

  :global(.target-title-part-text) {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 100%;
  }
</style>
