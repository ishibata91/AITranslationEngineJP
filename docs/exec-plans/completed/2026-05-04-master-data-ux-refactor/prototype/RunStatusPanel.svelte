<script lang="ts">
  type PreviewState = "empty" | "ready" | "running" | "complete" | "error"

  interface Props {
    previewState: PreviewState
    pauseGenerationAction: () => void
    cancelGenerationAction: () => void
  }

  let { previewState, pauseGenerationAction, cancelGenerationAction }: Props = $props()

  const isRunning = $derived(previewState === "running")
  const percent = $derived(
    previewState === "running" ? 54 : previewState === "complete" ? 100 : 0
  )
  const message = $derived(
    previewState === "running"
      ? "ペルソナを生成中"
      : previewState === "complete"
        ? "生成が完了しました"
        : previewState === "error"
          ? "生成を開始できません"
          : "生成はまだ開始されていません"
  )
</script>

<section class="panel run-panel" aria-labelledby="runStatusHeading">
  <div class="section-head">
    <div>
      <p class="eyebrow">進行状況</p>
      <h2 id="runStatusHeading">{message}</h2>
    </div>
    <span class="status-pill">{isRunning ? "実行中" : "待機中"}</span>
  </div>

  <div class="progress-track" aria-label="生成進捗">
    <div class="progress-fill" style={`width: ${percent}%`}></div>
  </div>

  <div class="run-grid">
    <article>
      <span>処理済み</span>
      <strong>{isRunning ? 23 : previewState === "complete" ? 42 : 0}</strong>
    </article>
    <article>
      <span>作成済み</span>
      <strong>{isRunning ? 9 : previewState === "complete" ? 18 : 0}</strong>
    </article>
    <article>
      <span>スキップ</span>
      <strong>{isRunning ? 14 : previewState === "complete" ? 24 : 0}</strong>
    </article>
  </div>

  <div class="current-box">
    <span>現在の対象</span>
    <strong>{isRunning ? "ファレンガー" : "-"}</strong>
  </div>

  <div class="button-row">
    <p>生成中も一覧と詳細を確認できます。</p>
    <div>
      <button
        class="button-secondary"
        disabled={!isRunning}
        onclick={pauseGenerationAction}
        type="button"
      >
        一時停止
      </button>
      <button
        class="button-secondary danger"
        disabled={!isRunning}
        onclick={cancelGenerationAction}
        type="button"
      >
        中止
      </button>
    </div>
  </div>
</section>

<style>
  .panel {
    background: var(--surface);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: var(--shadow);
    min-width: 0;
    padding: clamp(16px, 2vw, 22px);
  }

  .section-head,
  .button-row,
  .button-row div {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .progress-track {
    background: rgba(255, 255, 255, 0.07);
    border: 1px solid var(--line);
    border-radius: 999px;
    height: 12px;
    margin-top: 14px;
    overflow: hidden;
  }

  .progress-fill {
    background: linear-gradient(90deg, var(--primary), var(--accent));
    height: 100%;
  }

  .run-grid {
    display: grid;
    gap: 8px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    margin-top: 14px;
  }

  .run-grid article,
  .current-box {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 6px;
    padding: 12px;
  }

  .current-box {
    margin-top: 8px;
  }

  .eyebrow,
  .run-grid span,
  .current-box span {
    color: var(--primary);
    font-size: 0.78rem;
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--accent);
    flex: none;
    padding: 7px 10px;
    white-space: nowrap;
  }

  .button-row {
    color: var(--muted);
    line-height: 1.6;
    margin-top: 14px;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--line);
    border-radius: 7px;
    color: var(--text);
    cursor: pointer;
    min-height: 40px;
    padding: 9px 13px;
  }

  .button-secondary.danger {
    border-color: rgba(255, 140, 120, 0.5);
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  h2,
  p {
    margin: 0;
  }

  @media (max-width: 720px) {
    .run-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
