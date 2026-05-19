<script lang="ts">
  import type { CurrentPageHeroProps } from "./dashboard-component-props"

  let { currentRoute, dataTestId = undefined }: CurrentPageHeroProps = $props()
</script>

<section class="panel hero-panel" data-testid={dataTestId}>
  <div class="hero-top">
    <div>
      <p class="page-label">現在のページ</p>
      <h1>{currentRoute.label}</h1>
    </div>
  </div>
  <p class="hero-lead">{currentRoute.lead}</p>
</section>

<style>
  .panel {
    backdrop-filter: blur(38px);
    background: var(--surface);
    border: 0.5px solid var(--line);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow);
    color: var(--text);
  }

  .hero-panel {
    display: grid;
    gap: 16px;
    overflow: hidden;
    padding: 28px;
    position: relative;
  }

  .hero-panel::before {
    background: radial-gradient(
      circle,
      rgba(255, 186, 56, 0.14) 0%,
      transparent 72%
    );
    content: "";
    height: 220px;
    inset: 0 auto auto 0;
    pointer-events: none;
    position: absolute;
    width: 220px;
  }

  .hero-top {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    justify-content: space-between;
  }

  .page-label {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .hero-panel h1 {
    font-size: clamp(40px, 6vw, 72px);
    font-style: italic;
    letter-spacing: -0.05em;
    line-height: 0.96;
  }

  .hero-lead {
    color: var(--muted);
    font-size: 14px;
    line-height: 1.7;
  }
</style>
