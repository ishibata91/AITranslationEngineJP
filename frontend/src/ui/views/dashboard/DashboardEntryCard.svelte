<script lang="ts">
  import type { DashboardEntryCardProps } from "./dashboard-component-props"

  let { route, selectRoute }: DashboardEntryCardProps = $props()
</script>

<a
  class="entry-card"
  data-testid="dashboard-primary-page-card"
  href={`#${route.id}`}
  onclick={() => selectRoute(route.id)}
>
  <div class="entry-head">
    <div>
      <p class="card-tag">{route.id}</p>
      <h3>{route.label}</h3>
    </div>
    <span class="entry-state">{route.state}</span>
  </div>
  <p>{route.description}</p>
  <span class="entry-action">開く</span>
</a>

<style>
  .entry-card {
    background: rgba(17, 13, 12, 0.34);
    border: 0.5px solid rgba(255, 186, 56, 0.12);
    border-radius: var(--radius-md);
    color: var(--text);
    display: block;
    overflow: hidden;
    padding: 18px;
    position: relative;
    text-decoration: none;
    transition:
      transform var(--transition),
      border-color var(--transition),
      background var(--transition),
      box-shadow var(--transition);
  }

  .entry-card::before {
    color: rgba(255, 186, 56, 0.26);
    content: "✦";
    font-size: 12px;
    position: absolute;
    right: 16px;
    top: 14px;
  }

  .entry-card:hover,
  .entry-card:focus-visible {
    background: rgba(255, 186, 56, 0.08);
    border-color: rgba(255, 186, 56, 0.24);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.26);
    outline: none;
    transform: translateY(-2px);
  }

  .entry-head {
    align-items: flex-start;
    display: flex;
    gap: 12px;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .card-tag {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  .entry-card h3 {
    font-size: 20px;
  }

  .entry-card p {
    color: var(--muted);
    font-size: 14px;
    line-height: 1.7;
  }

  .entry-state {
    align-items: center;
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 186, 56, 0.14);
    border-radius: 999px;
    color: var(--muted);
    display: inline-flex;
    font-size: 12px;
    letter-spacing: 0.08em;
    min-height: 28px;
    padding: 0 10px;
    text-transform: uppercase;
  }

  .entry-action {
    align-items: center;
    color: var(--primary);
    display: inline-flex;
    font-size: 13px;
    gap: 8px;
    letter-spacing: 0.08em;
    margin-top: 14px;
    text-transform: uppercase;
  }
</style>
