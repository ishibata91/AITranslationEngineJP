<script lang="ts">
  import type { GlobalNavigationProps } from "./dashboard-component-props"

  let { currentRoute, routes, selectRoute }: GlobalNavigationProps = $props()
</script>

<nav
  aria-label="グローバルナビゲーション"
  class="global-nav"
  data-testid="dashboard-global-navigation"
  id="globalNav"
>
  {#each routes as route (route.id)}
    <a
      aria-current={route.id === currentRoute.id ? "page" : undefined}
      class="nav-link"
      class:is-active={route.id === currentRoute.id}
      data-testid="dashboard-global-navigation-item"
      href={`#${route.id}`}
      onclick={() => selectRoute(route.id)}
    >
      {route.label}
    </a>
  {/each}
</nav>

<style>
  .global-nav {
    color: var(--text);
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: center;
  }

  .nav-link {
    align-items: center;
    border: 0.5px solid transparent;
    border-radius: 999px;
    color: var(--muted);
    cursor: pointer;
    display: inline-flex;
    font: inherit;
    min-height: 42px;
    padding: 10px 14px;
    text-decoration: none;
    transition:
      background var(--transition),
      color var(--transition),
      border-color var(--transition),
      transform var(--transition);
  }

  .nav-link:hover,
  .nav-link:focus-visible {
    background: rgba(255, 255, 255, 0.04);
    border-color: rgba(255, 186, 56, 0.14);
    color: var(--text);
    outline: none;
  }

  .nav-link.is-active {
    background: linear-gradient(135deg, var(--primary) 0%, #e8a31d 100%);
    box-shadow: 0 0 20px rgba(255, 186, 56, 0.18);
    color: var(--bg-strong);
  }

  @media (max-width: 860px) {
    .global-nav {
      display: none;
      justify-content: flex-start;
    }

    :global(.shell-bar.is-open) .global-nav {
      display: flex;
    }
  }
</style>
