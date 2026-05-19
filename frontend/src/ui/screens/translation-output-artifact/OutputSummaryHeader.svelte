<script lang="ts">
  const PLACEHOLDER_LEAD =
    "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

  interface Props {
    gatewayStatus: string
    statusTitle: string
    statusText: string
  }

  let { gatewayStatus, statusTitle, statusText }: Props = $props()
</script>

<section
  class="output-card hero-card"
  data-testid="output-management-output-management-summary"
>
  <div class="hero-head">
    <div>
      <p class="eyebrow">output-management</p>
      <h2>Output Review</h2>
    </div>
    <p class="gateway-status">Gateway: {gatewayStatus}</p>
  </div>
  <p class="lead">
    completed job、result summary、diff preview、artifact 状態を確認し、
    xTranslator XML の出力と再出力を同じ画面で行います。
  </p>
  <p class="status-copy">
    <strong>{statusTitle}</strong>
    <span>{statusText}</span>
  </p>
  <p class="placeholder-copy">{PLACEHOLDER_LEAD}</p>
</section>

<style>
  .output-card {
    border: 1px solid var(--line);
    border-radius: var(--radius-md);
    background: rgba(25, 22, 20, 0.82);
    box-shadow: var(--shadow);
    color: var(--text);
    padding: 1.25rem;
  }

  .hero-card {
    display: grid;
    gap: 0.9rem;
  }

  .hero-head {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow,
  .placeholder-copy {
    color: var(--muted);
    font-size: 0.85rem;
  }

  .lead,
  .status-copy {
    line-height: 1.6;
  }
</style>
