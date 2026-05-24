<script lang="ts">
  import { tick } from "svelte"

  interface Props {
    tooltipId: string
    triggerText: string
    tooltipTitle?: string
    tooltipItems?: string[]
    tooltipText?: string
    secondary?: boolean
    showWhenTruncatedOnly?: boolean
    triggerClass?: string
    triggerAriaLabel?: string
  }

  let {
    tooltipId,
    triggerText,
    tooltipTitle = "",
    tooltipItems = [],
    tooltipText = "",
    secondary = false,
    showWhenTruncatedOnly = false,
    triggerClass = "",
    triggerAriaLabel = ""
  }: Props = $props()

  let tooltipElement = $state<HTMLSpanElement | null>(null)
  let isTooltipVisible = $state(false)
  let tooltipTop = $state(0)
  let tooltipLeft = $state(0)
  let tooltipWidth = $state(0)
  let canShowTooltip = $state(false)

  const viewportPadding = 16
  const tooltipGap = 8
  const maxTooltipWidth = 672

  const tooltipStyle = $derived(
    `top: ${tooltipTop}px; left: ${tooltipLeft}px; max-width: ${tooltipWidth}px;`
  )
  const effectiveTooltipText = $derived(tooltipText || triggerText)
  const shouldRenderTooltipContent = $derived(
    !showWhenTruncatedOnly || isTooltipVisible
  )

  async function showTooltip(anchor: HTMLElement): Promise<void> {
    canShowTooltip = !showWhenTruncatedOnly || isTextTruncated(anchor)
    if (!canShowTooltip) {
      hideTooltip()
      return
    }

    const anchorRect = anchor.getBoundingClientRect()
    const viewportWidth = window.innerWidth
    const availableWidth = Math.max(viewportWidth - viewportPadding * 2, 0)

    tooltipWidth = Math.min(maxTooltipWidth, availableWidth)
    tooltipLeft = clampTooltipLeft(
      anchorRect.left + anchorRect.width / 2 - tooltipWidth / 2,
      viewportWidth
    )
    tooltipTop = Math.max(viewportPadding, anchorRect.top - tooltipGap)
    isTooltipVisible = true

    await tick()

    if (!tooltipElement) {
      return
    }

    const tooltipRect = tooltipElement.getBoundingClientRect()
    tooltipTop = Math.max(
      viewportPadding,
      anchorRect.top - tooltipRect.height - tooltipGap
    )
  }

  function hideTooltip(): void {
    isTooltipVisible = false
  }

  function isTextTruncated(anchor: HTMLElement): boolean {
    return anchor.scrollWidth > anchor.clientWidth
  }

  function clampTooltipLeft(left: number, viewportWidth: number): number {
    const maxLeft = Math.max(
      viewportPadding,
      viewportWidth - viewportPadding - tooltipWidth
    )

    return Math.min(Math.max(left, viewportPadding), maxLeft)
  }
</script>

<span class="tooltip-trigger-shell" class:secondary>
  <button
    class={`tooltip-trigger ${triggerClass}`}
    type="button"
    aria-label={triggerAriaLabel || undefined}
    aria-describedby={canShowTooltip ? tooltipId : undefined}
    onmouseenter={(event) => showTooltip(event.currentTarget)}
    onmouseleave={hideTooltip}
    onfocus={(event) => showTooltip(event.currentTarget)}
    onblur={hideTooltip}
  >
    {triggerText}
  </button>
  <span
    bind:this={tooltipElement}
    class:visible={isTooltipVisible}
    class="floating-tooltip"
    id={tooltipId}
    role="tooltip"
    aria-hidden={!isTooltipVisible}
    style={tooltipStyle}
  >
    {#if shouldRenderTooltipContent}
      {#if tooltipTitle}
        <span class="tooltip-title">{tooltipTitle}</span>
      {/if}
      {#if tooltipItems.length > 0}
        <ul>
          {#each tooltipItems as item (item)}
            <li>{item}</li>
          {/each}
        </ul>
      {:else}
        {effectiveTooltipText}
      {/if}
    {/if}
  </span>
</span>

<style>
  .tooltip-trigger-shell {
    box-sizing: border-box;
    display: inline-block;
    min-width: 0;
  }

  .tooltip-trigger {
    appearance: none;
    background: transparent;
    border: 0;
    border-radius: 4px;
    color: inherit;
    display: block;
    font: inherit;
    min-width: 0;
    padding: 0;
    text-align: left;
  }

  .tooltip-trigger:focus {
    outline: 2px solid rgba(255, 212, 165, 0.44);
    outline-offset: 2px;
  }

  .floating-tooltip {
    background: rgba(18, 16, 15, 0.96);
    border: 1px solid rgba(255, 212, 165, 0.3);
    border-radius: 8px;
    box-shadow: 0 12px 28px rgba(0, 0, 0, 0.34);
    box-sizing: border-box;
    color: #fff6ea;
    display: grid;
    gap: 0.45rem;
    line-height: 1.5;
    max-height: calc(100vh - 2rem);
    opacity: 0;
    overflow-y: auto;
    overflow-wrap: anywhere;
    padding: 0.55rem 0.7rem;
    pointer-events: none;
    position: fixed;
    transform: translateY(-0.12rem);
    transition:
      opacity 120ms ease,
      transform 120ms ease,
      visibility 120ms ease;
    visibility: hidden;
    white-space: normal;
    width: 100%;
    z-index: 30;
  }

  .floating-tooltip.visible {
    opacity: 1;
    transform: translateY(0);
    visibility: visible;
  }

  .tooltip-title {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.76rem;
  }

  ul {
    display: grid;
    gap: 0.35rem;
    margin: 0;
    padding-left: 1.1rem;
  }

  li {
    overflow-wrap: anywhere;
    word-break: break-word;
  }
</style>
