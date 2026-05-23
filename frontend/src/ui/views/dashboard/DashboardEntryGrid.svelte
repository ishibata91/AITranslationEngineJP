<script lang="ts">
  import DashboardEntryCard from "./DashboardEntryCard.svelte"
  import type { DashboardEntryGridProps } from "./dashboard-component-props"

  let { routes, selectRoute }: DashboardEntryGridProps = $props()
</script>

<section class="hero-grid" id="dashboardView">
  <section
    class="panel entry-panel"
    data-testid="dashboard-primary-page-entry-region"
  >
    <div class="section-head">
      <div>
        <p class="page-label">主要ページ</p>
        <h2>作業を選ぶ</h2>
      </div>
    </div>
    <div class="entry-grid">
      {#each routes as route (route.id)}
        <DashboardEntryCard {route} {selectRoute} />
      {/each}
    </div>
  </section>
</section>

<style>
  .hero-grid {
    display: grid;
    gap: 20px;
    grid-template-columns: 1fr;
  }

  .panel {
    backdrop-filter: blur(38px);
    background: var(--surface);
    border: 0.5px solid var(--line);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow);
    color: var(--text);
  }

  .entry-panel {
    padding: 24px;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    justify-content: space-between;
    margin-bottom: 18px;
  }

  .section-head h2 {
    font-size: 24px;
  }

  .page-label {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .entry-grid {
    display: grid;
    gap: 14px;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  @media (max-width: 860px) {
    .entry-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
