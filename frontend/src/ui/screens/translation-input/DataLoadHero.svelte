<script lang="ts">
  interface Props {
    gatewayStatus: string
    operationStatusLabel: string
    operationStatusText: string
    errorMessage: string
  }

  let {
    gatewayStatus,
    operationStatusLabel,
    operationStatusText,
    errorMessage
  }: Props = $props()
</script>

<section class="panel hero-panel">
  <div class="hero-head">
    <div class="title-stack">
      <p class="eyebrow">翻訳管理</p>
      <h2 class="page-title">データロード</h2>
      <p class="lead">
        JSON を読み込み、登録結果と再構築の要否を同じ画面で確認する。
      </p>
    </div>
    <p class="gateway-status">接続状態: {gatewayStatus}</p>
  </div>

  <div class="status-band" aria-label="現在の作業状態">
    <strong>{operationStatusLabel}</strong>
    <p>{operationStatusText}</p>
  </div>

  <p class="error-text" hidden={!errorMessage}>
    {errorMessage}
  </p>
</section>

<style>
  .panel {
    border: 1px solid var(--line);
    border-radius: 20px;
    padding: 1.25rem;
    background:
      radial-gradient(
        circle at top right,
        rgba(255, 186, 56, 0.16),
        transparent 32%
      ),
      rgba(28, 23, 20, 0.74);
    box-shadow: var(--shadow);
    backdrop-filter: blur(18px);
  }

  .hero-panel {
    display: grid;
    gap: 1rem;
  }

  .hero-head {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .title-stack {
    display: grid;
    gap: 0.45rem;
  }

  .eyebrow {
    margin: 0;
    color: var(--muted);
    font-size: 0.76rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  h2,
  p {
    margin: 0;
  }

  .page-title {
    font-size: clamp(1.7rem, 4vw, 2.3rem);
    line-height: 1.15;
  }

  .lead {
    max-width: 44rem;
    color: var(--muted);
  }

  .gateway-status {
    margin: 0;
    padding: 0.32rem 0.7rem;
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--primary);
    font-size: 0.82rem;
    white-space: nowrap;
  }

  .status-band {
    display: grid;
    gap: 0.2rem;
    padding: 1rem;
    border-radius: 18px;
    border: 1px solid rgba(255, 186, 56, 0.16);
    background: rgba(17, 15, 14, 0.72);
  }

  .status-band p {
    color: var(--muted);
  }

  .error-text {
    margin: 0;
    padding: 0.8rem 1rem;
    border-radius: 14px;
    background: rgba(255, 104, 63, 0.16);
    color: #ffd8c5;
  }

  @media (max-width: 720px) {
    .hero-head {
      flex-direction: column;
    }

    .gateway-status {
      white-space: normal;
    }
  }
</style>
