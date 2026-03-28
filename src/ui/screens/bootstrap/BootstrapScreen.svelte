<script lang="ts">
  import { onMount } from "svelte";
  import { loadBootstrapStatus } from "@application/bootstrap/load-bootstrap-status";
  import type { BootstrapStatus } from "@shared/contracts/bootstrap-status";

  let status: BootstrapStatus | null = null;
  let errorMessage = "";

  onMount(async () => {
    try {
      status = await loadBootstrapStatus();
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : "Unknown bootstrap failure";
    }
  });
</script>

<section class="panel">
  <p class="eyebrow">Implementation Bootstrap</p>
  <h1>AITranslationEngineJp</h1>
  <p class="lede">
    Frontend は `Gateway` 経由で backend bootstrap status を取得し、最小の責務境界を可視化します。
  </p>

  {#if errorMessage}
    <p class="status error">{errorMessage}</p>
  {:else if status}
    <dl class="status-grid">
      <div>
        <dt>Backend Version</dt>
        <dd>{status.backendVersion}</dd>
      </div>
      <div>
        <dt>Boundary Ready</dt>
        <dd>{status.boundaryReady ? "true" : "false"}</dd>
      </div>
      <div>
        <dt>Frontend Entry</dt>
        <dd>{status.frontendEntry}</dd>
      </div>
    </dl>
  {:else}
    <p class="status loading">Loading bootstrap status...</p>
  {/if}
</section>

<style>
  .panel {
    width: min(100%, 42rem);
    padding: 2rem;
    border: 1px solid #c8d2e0;
    border-radius: 1.25rem;
    background: rgba(255, 255, 255, 0.88);
    box-shadow: 0 24px 60px rgba(20, 32, 51, 0.12);
  }

  .eyebrow {
    margin: 0 0 0.75rem;
    font-size: 0.8rem;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: #5b6a7d;
  }

  h1 {
    margin: 0;
    font-size: clamp(2rem, 4vw, 3rem);
  }

  .lede {
    margin: 1rem 0 1.5rem;
    line-height: 1.6;
    color: #334155;
  }

  .status-grid {
    display: grid;
    gap: 1rem;
    margin: 0;
  }

  .status-grid div {
    padding: 1rem;
    border-radius: 0.9rem;
    background: #f3f6fb;
  }

  dt {
    font-size: 0.85rem;
    color: #526277;
  }

  dd {
    margin: 0.3rem 0 0;
    font-size: 1.1rem;
    font-weight: 600;
  }

  .status {
    margin: 0;
    padding: 1rem;
    border-radius: 0.9rem;
  }

  .loading {
    background: #edf4ff;
  }

  .error {
    background: #fff1f2;
    color: #9f1239;
  }
</style>
