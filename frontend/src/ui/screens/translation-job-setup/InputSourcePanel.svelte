<script lang="ts">
  interface InputSourceCandidate {
    id: number
    label: string
    sourceKind: string
    recordCount: number
    registeredAt?: string
  }

  interface Props {
    candidates: InputSourceCandidate[]
    deletingInputSourceId: number | null
    existingJobSummary: string
    isCreating: boolean
    selectedInputLabel: string
    selectedInputRecordCountLabel: string
    selectedInputRegisteredAtLabel: string
    selectedInputSourceId: number | null
    selectedInputSourceKind: string
    formatDate: (timestamp: string) => string
    onDeleteInputSource: (candidateId: number) => void
    onSelectInputSource: (candidateId: number) => void
  }

  let {
    candidates,
    deletingInputSourceId,
    existingJobSummary,
    isCreating,
    selectedInputLabel,
    selectedInputRecordCountLabel,
    selectedInputRegisteredAtLabel,
    selectedInputSourceId,
    selectedInputSourceKind,
    formatDate,
    onDeleteInputSource,
    onSelectInputSource
  }: Props = $props()

  function isSelectedInputCard(candidateId: number): boolean {
    return selectedInputSourceId === candidateId
  }

  function isDeletingInputCard(candidateId: number): boolean {
    return deletingInputSourceId === candidateId
  }
</script>

<section
  class="job-setup-card"
  aria-labelledby="jobSetupInputHeading"
  data-testid="translation-job-setup-input-data-region"
>
  <div class="section-head">
    <div>
      <h3 id="jobSetupInputHeading">入力データ</h3>
    </div>
  </div>
  <div class="input-card-list" aria-label="input data" role="list">
    {#each candidates as candidate (candidate.id)}
      <article
        aria-busy={isDeletingInputCard(candidate.id)}
        class:selected={isSelectedInputCard(candidate.id)}
        class="input-card"
        role="listitem"
      >
        <button
          aria-pressed={isSelectedInputCard(candidate.id)}
          class="input-card-select"
          disabled={isCreating || isDeletingInputCard(candidate.id)}
          onclick={() => onSelectInputSource(candidate.id)}
          type="button"
        >
          <div class="input-card-head">
            <strong class="wrap-value">{candidate.label}</strong>
            {#if isDeletingInputCard(candidate.id)}
              <span class="status-pill">削除中...</span>
            {:else if isSelectedInputCard(candidate.id)}
              <span class="status-pill success">選択中</span>
            {/if}
          </div>
          <dl class="detail-grid compact">
            <div>
              <dt>出自</dt>
              <dd class="wrap-value">{candidate.sourceKind}</dd>
            </div>
            <div>
              <dt>翻訳レコード件数</dt>
              <dd>{candidate.recordCount.toLocaleString("ja-JP")} 件</dd>
            </div>
            <div>
              <dt>登録日時</dt>
              <dd>{formatDate(candidate.registeredAt ?? "")}</dd>
            </div>
          </dl>
        </button>
        <div class="input-card-actions">
          <button
            class="button-secondary"
            disabled={isCreating || deletingInputSourceId !== null}
            onclick={() => onDeleteInputSource(candidate.id)}
            type="button"
          >
            {isDeletingInputCard(candidate.id) ? "削除中..." : "削除"}
          </button>
        </div>
      </article>
    {/each}
  </div>
  <dl class="detail-grid compact">
    <div>
      <dt>入力データ名</dt>
      <dd class="wrap-value">{selectedInputLabel}</dd>
    </div>
    <div>
      <dt>出自</dt>
      <dd class="wrap-value">{selectedInputSourceKind}</dd>
    </div>
    <div>
      <dt>登録日時</dt>
      <dd>{selectedInputRegisteredAtLabel}</dd>
    </div>
    <div>
      <dt>翻訳レコード件数</dt>
      <dd>{selectedInputRecordCountLabel}</dd>
    </div>
    <div>
      <dt>既存 job 状態</dt>
      <dd class="wrap-value">{existingJobSummary}</dd>
    </div>
  </dl>
</section>

<style>
  .job-setup-card,
  .input-card,
  .input-card-actions,
  .input-card-select,
  .input-card-list,
  .detail-grid div {
    display: grid;
    gap: 0.75rem;
  }

  .job-setup-card {
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .section-head,
  .input-card-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .input-card-list {
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  }

  .input-card {
    padding: 0.9rem;
    border-radius: 1rem;
    border: 1px solid rgba(255, 212, 165, 0.12);
    background: rgba(18, 13, 11, 0.52);
  }

  .input-card.selected {
    border-color: rgba(255, 204, 136, 0.72);
    background: rgba(56, 39, 30, 0.78);
    box-shadow: 0 0 0 1px rgba(255, 204, 136, 0.22);
  }

  .input-card-select {
    text-align: left;
    border: 0;
    padding: 0;
    background: transparent;
    color: inherit;
    cursor: pointer;
    font: inherit;
  }

  .input-card-actions {
    align-items: start;
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  dt {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.9rem;
  }

  .status-pill {
    padding: 0.35rem 0.72rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.82rem;
  }

  .status-pill.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .button-secondary {
    padding: 0.8rem 1rem;
    border-radius: 0.9rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    background: rgba(255, 241, 227, 0.08);
    color: #ffe2bf;
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }
</style>
