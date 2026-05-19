<script lang="ts">
  import type { DictionaryHeaderProps } from "./dictionary-panel-props"

  let { gatewayStatus, errorMessage }: DictionaryHeaderProps = $props()
</script>

<section class="shell-card" data-testid="master-dictionary-screen-header">
  <div class="hero-top">
    <div>
      <p class="eyebrow">基盤データ</p>
      <h2>マスター辞書</h2>
    </div>
    <p class="gateway-status" id="gatewayStatus">Gateway: {gatewayStatus}</p>
  </div>
  <p class="lead">
    一覧、詳細、作成、更新、削除、XML 取り込みを同じ画面で操作できます。
  </p>
  <p class="error-text" hidden={!errorMessage} id="masterDictionaryError">
    {errorMessage}
  </p>
</section>

<style>
  .shell-card {
    padding: 18px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(16, 13, 11, 0.58);
    color: var(--text);
  }

  .hero-top {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  .lead,
  .gateway-status {
    color: var(--muted);
  }

  .error-text {
    color: #ffc0ab;
  }

  .error-text[hidden] {
    display: none;
  }

  .eyebrow {
    font-size: 12px;
    letter-spacing: 0.08em;
  }
</style>
