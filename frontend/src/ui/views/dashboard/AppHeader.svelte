<script lang="ts">
  import GlobalNavigation from "./GlobalNavigation.svelte"
  import type { AppHeaderProps } from "./dashboard-component-props"

  let {
    currentRoute,
    isMobileNavOpen,
    routes,
    selectRoute,
    toggleMobileNav
  }: AppHeaderProps = $props()
</script>

<header
  class="shell-bar"
  class:is-open={isMobileNavOpen}
  data-testid="dashboard-application-header"
>
  <div class="shell-bar-inner">
    <div class="brand" data-testid="dashboard-brand-display">
      <p class="brand-eyebrow">AITranslationEngineJp</p>
      <strong>翻訳エンジン</strong>
    </div>
    <GlobalNavigation {currentRoute} {routes} {selectRoute} />
    <div class="bar-status">
      <button
        aria-controls="globalNav"
        aria-expanded={isMobileNavOpen ? "true" : "false"}
        class="nav-toggle"
        data-testid="dashboard-mobile-navigation-toggle"
        onclick={toggleMobileNav}
        type="button"
      >
        主要ページ
      </button>
    </div>
  </div>
</header>

<style>
  .shell-bar {
    backdrop-filter: blur(38px);
    background: rgba(22, 19, 17, 0.78);
    border-bottom: 0.5px solid var(--line);
    box-shadow: 0 20px 50px rgba(0, 0, 0, 0.38);
    color: var(--text);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .shell-bar-inner {
    align-items: center;
    display: grid;
    gap: 18px;
    grid-template-columns: auto 1fr auto;
    margin: 0 auto;
    max-width: 1440px;
    padding: 18px 24px;
  }

  .brand {
    display: grid;
    gap: 2px;
  }

  .brand-eyebrow {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .brand strong {
    color: var(--primary);
    font-size: 20px;
    font-style: italic;
    font-weight: 900;
  }

  .bar-status {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: flex-end;
  }

  .nav-toggle {
    background: linear-gradient(135deg, var(--primary) 0%, #f3a114 100%);
    border: 0;
    border-radius: 999px;
    color: #432c00;
    cursor: pointer;
    display: none;
    font: inherit;
    min-height: 42px;
    padding: 0 16px;
  }

  @media (max-width: 1120px) {
    .shell-bar-inner {
      grid-template-columns: 1fr;
    }

    .bar-status {
      justify-content: flex-start;
    }
  }

  @media (max-width: 860px) {
    .nav-toggle {
      align-items: center;
      display: inline-flex;
      justify-content: center;
    }
  }
</style>
