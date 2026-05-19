<script lang="ts">
  import type { Snippet } from "svelte"

  interface Props {
    children?: Snippet
    footer?: Snippet
  }

  let { children, footer }: Props = $props()
</script>

<section class="job-run-phase-region" data-testid="job-run-phase-screen-region">
  {#if children}
    {@render children()}
  {:else}
    <section class="phase-host-placeholder">
      <p class="eyebrow">翻訳段階</p>
      <h3>翻訳段階の表示領域</h3>
      <p>選択中の翻訳段階画面をここに表示します。</p>
    </section>
  {/if}
  {@render footer?.()}
</section>

<style>
  .job-run-phase-region {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
  }

  .phase-host-placeholder {
    display: grid;
    gap: 0.5rem;
    min-width: 0;
    padding: 1.15rem 1.3rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    color: rgba(236, 223, 205, 0.78);
  }

  .eyebrow {
    margin: 0;
    color: rgba(236, 223, 205, 0.78);
  }

  h3 {
    margin: 0;
    color: #fff6ea;
  }
</style>
