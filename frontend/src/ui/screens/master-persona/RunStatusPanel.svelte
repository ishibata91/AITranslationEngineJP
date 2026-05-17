<script lang="ts">
  import type { RunStatusPanelProps } from "./master-persona-panel-props"

  let {
    isRunActive,
    progressPercent,
    runStatus,
    interruptGeneration,
    cancelGeneration
  }: RunStatusPanelProps = $props()

  const statusLabel = $derived(isRunActive ? "生成中" : runStatus.runState)
  const lockText = $derived(
    isRunActive
      ? "生成中は編集と削除を行えません。"
      : "一覧と詳細を確認しながら次の操作を選べます。"
  )
</script>

<section
  class="panel run-panel"
  aria-labelledby="runHeading"
  data-testid="master-persona-progress-panel"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">進行状況</p>
      <h2 id="runHeading">{runStatus.message}</h2>
      <p class="support-copy" role="status">{lockText}</p>
    </div>
    <span class:status-danger={isRunActive} class="status-pill"
      >{statusLabel}</span
    >
  </div>

  <div class="progress-track" aria-label="生成進捗">
    <div
      class="progress-fill"
      id="runProgressFill"
      style={`width: ${progressPercent}%;`}
    ></div>
  </div>

  <div class="run-grid">
    <article class="run-card">
      <span>処理済み件数</span>
      <strong>{runStatus.processedCount}</strong>
    </article>
    <article class="run-card">
      <span>作成済み件数</span>
      <strong>{runStatus.successCount}</strong>
    </article>
    <article class="run-card">
      <span>既に作成済み</span>
      <strong>{runStatus.existingSkipCount}</strong>
    </article>
    <article class="run-card current-card">
      <span>現在の対象</span>
      <strong>{runStatus.currentActorLabel || "-"}</strong>
    </article>
  </div>

  <div class="button-row">
    <button
      class="button-secondary"
      disabled={!isRunActive}
      id="interruptGenerationButton"
      onclick={interruptGeneration}
      type="button"
    >
      一時停止
    </button>
    <button
      class="button-danger"
      disabled={!isRunActive}
      id="cancelGenerationButton"
      onclick={cancelGeneration}
      type="button"
    >
      中止
    </button>
  </div>
</section>

<style>
  .panel,
  .run-card {
    border-radius: 20px;
  }

  .panel {
    background: rgba(17, 13, 12, 0.42);
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    display: grid;
    gap: 14px;
    min-width: 0;
    padding: clamp(18px, 3vw, 24px);
  }

  .section-head,
  .button-row {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .eyebrow {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.1em;
    margin: 0 0 6px;
    text-transform: uppercase;
  }

  h2,
  p {
    margin: 0;
  }

  h2,
  strong,
  .support-copy {
    overflow-wrap: anywhere;
  }

  .support-copy {
    color: var(--muted);
    line-height: 1.7;
    margin-top: 8px;
  }

  .status-pill {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    display: inline-flex;
    min-height: 38px;
    padding: 0 14px;
  }

  .status-danger {
    background: rgba(255, 156, 124, 0.14);
    border-color: rgba(255, 156, 124, 0.32);
  }

  .progress-track {
    background: rgba(255, 255, 255, 0.07);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    border-radius: 999px;
    height: 12px;
    overflow: hidden;
  }

  .progress-fill {
    background: linear-gradient(90deg, var(--primary), #9ad7cb);
    height: 100%;
    transition: width 160ms ease;
  }

  .run-grid {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .run-card {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 14px;
  }

  .current-card {
    grid-column: 1 / -1;
  }

  .run-card span {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .button-secondary,
  .button-danger {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    padding: 0 16px;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    color: var(--text);
  }

  .button-danger {
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9c7c 100%);
    border: 0.5px solid transparent;
    color: #35150d;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  @media (max-width: 720px) {
    .run-grid {
      grid-template-columns: 1fr;
    }

    .button-row > * {
      width: 100%;
    }
  }
</style>
